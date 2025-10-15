package tournamentmanagerrepo

import (
	"fmt"
	"tournament-manager/domain"
)

func (r *tournamentManagerRepo) CreateMatchSchedules(tournamentID int, groupCount int, approvedParticipants []*domain.Participant) error {
	// 1️⃣ Generate groups first
	if err := r.GenerateGroups(tournamentID, groupCount, approvedParticipants); err != nil {
		return fmt.Errorf("failed to generate groups: %w", err)
	}

	// 2️⃣ Fetch all group IDs for this tournament
	groupRows, err := r.db.Query(`SELECT id FROM groups WHERE tournament_id = $1`, tournamentID)
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
		rows, err := r.db.Query(`SELECT participant_id FROM group_teams WHERE group_id = $1`, groupID)
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
	matchStmt, err := tx.Prepare(`
		INSERT INTO matches (tournament_id, group_id, round, participant_a_id, participant_b_id)
		VALUES ($1, $2, $3, $4, $5)
	`)
	if err != nil {
		return fmt.Errorf("failed to prepare match insert: %w", err)
	}
	defer matchStmt.Close()

	playerStatStmt, err := tx.Prepare(`
		INSERT INTO player_stats (group_id, participant_id)
		VALUES ($1, $2)
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
			if _, err := playerStatStmt.Exec(groupID, pid); err != nil {
				return fmt.Errorf("failed to insert player stats for participant %d in group %d: %w", pid, groupID, err)
			}
		}

		// Create round-robin matches
		for i := 0; i < len(participants); i++ {
			for j := i + 1; j < len(participants); j++ {
				if _, err := matchStmt.Exec(tournamentID, groupID, "Group Stage", participants[i], participants[j]); err != nil {
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
