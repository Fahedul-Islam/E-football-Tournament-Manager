package tournamentmanagerrepo

import (
	"context"
	"fmt"
	"math/rand"
	"time"
	"tournament-manager/internal/domain"
)

func (r *tournamentManagerRepo) CreateMatchSchedules(ctx context.Context, tournamentID int, groupCount int, approvedParticipants []*domain.Participant) error {
	// 1️⃣ Generate groups first
	if err := r.GenerateGroups(ctx, tournamentID, groupCount, approvedParticipants); err != nil {
		return fmt.Errorf("failed to generate groups: %w", err)
	}

	// 2️⃣ Fetch all group IDs for this tournament
	groupRows, err := r.db.QueryContext(ctx, `SELECT id FROM groups WHERE tournament_id = $1`, tournamentID)
	if err != nil {
		return fmt.Errorf("failed to fetch groups: %w", err)
	}
	defer groupRows.Close()

	var groupIDs []int
	for groupRows.Next() {
		var id int
		if err := groupRows.Scan(&id); err != nil {
			return fmt.Errorf("failed to scan group ID: %w", err)
		}
		groupIDs = append(groupIDs, id)
	}
	if len(groupIDs) == 0 {
		return fmt.Errorf("no groups found for tournament %d", tournamentID)
	}

	// 3️⃣ Get participants per group
	participantsPerGroup := make(map[int][]int)
	for _, groupID := range groupIDs {
		rows, err := r.db.QueryContext(ctx, `SELECT participant_id FROM group_teams WHERE group_id = $1`, groupID)
		if err != nil {
			return fmt.Errorf("failed to fetch participants for group %d: %w", groupID, err)
		}

		var pids []int
		for rows.Next() {
			var pid int
			if err := rows.Scan(&pid); err != nil {
				rows.Close()
				return fmt.Errorf("failed to scan participant ID for group %d: %w", groupID, err)
			}
			pids = append(pids, pid)
		}
		rows.Close()
		participantsPerGroup[groupID] = pids
	}

	// 4️⃣ Begin transaction
	tx, err := r.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() { _ = tx.Rollback() }()

	// 5️⃣ Prepare insert statements
	matchStmt, err := tx.PrepareContext(ctx, `
		INSERT INTO matches (tournament_id, group_id, round, participant_a_id, participant_b_id)
		VALUES ($1, $2, $3, $4, $5)
	`)
	if err != nil {
		return fmt.Errorf("failed to prepare match insert: %w", err)
	}
	defer matchStmt.Close()

	playerStatStmt, err := tx.PrepareContext(ctx, `
		INSERT INTO player_stats (tournament_id, group_id, participant_id)
		VALUES ($1, $2, $3)
		ON CONFLICT (group_id, participant_id) DO NOTHING
	`)
	if err != nil {
		return fmt.Errorf("failed to prepare player_stats insert: %w", err)
	}
	defer playerStatStmt.Close()

	// 6️⃣ Generate round-robin matches and insert player stats
	for groupID, participants := range participantsPerGroup {
		// Insert player stats for each participant
		for _, pid := range participants {
			if _, err := playerStatStmt.ExecContext(ctx, tournamentID, groupID, pid); err != nil {
				return fmt.Errorf("failed to insert player stats for participant %d in group %d: %w", pid, groupID, err)
			}
		}

		// Create round-robin matches
		for i := 0; i < len(participants); i++ {
			for j := i + 1; j < len(participants); j++ {
				if _, err := matchStmt.ExecContext(ctx, tournamentID, groupID, "Group Stage", participants[i], participants[j]); err != nil {
					return fmt.Errorf("failed to insert match for group %d: %w", groupID, err)
				}
			}
		}
	}

	// 7️⃣ Commit transaction
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (r *tournamentManagerRepo) LeagueStyleSchedule(ctx context.Context, tournament_id int, approvedParticipants []*domain.Participant) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(approvedParticipants), func(i, j int) {
		approvedParticipants[i], approvedParticipants[j] = approvedParticipants[j], approvedParticipants[i]
	})

	for i, participantA := range approvedParticipants {
		playerStatStmt, err := tx.PrepareContext(ctx, `
			INSERT INTO player_stats (tournament_id, participant_id) VALUES ($1, $2)
		`)
		if err != nil {
			return err
		}
		defer playerStatStmt.Close()
		if _, err := playerStatStmt.ExecContext(ctx, tournament_id, participantA.ID); err != nil {
			return fmt.Errorf("fail to insert player stat: %w", err)
		}
		for j, participantB := range approvedParticipants {
			if i != j {
				matchStmt, err := tx.PrepareContext(ctx, `
				INSERT INTO matches (tournament_id, group_id, round, participant_a_id, participant_b_id)
				VALUES ($1, $2, $3, $4, $5)`)
				if err != nil {
					return fmt.Errorf("failed to prepare match insert: %w", err)
				}
				defer matchStmt.Close()
				if _, err := matchStmt.ExecContext(ctx, tournament_id, nil, "league", participantA.ID, participantB.ID); err != nil {
					return fmt.Errorf("failed to insert matches: %w", err)
				}
			}
		}
	}
	return nil
}
