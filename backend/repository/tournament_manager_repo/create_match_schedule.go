package tournamentmanagerrepo

import (
	"fmt"
	"tournament-manager/domain"
)

func (r *tournamentManagerRepo) CreateMatchSchedules(tournament_id int, groupCount int, approvedParticipants []*domain.Participant) error {

	err := r.GenerateGroups(tournament_id, groupCount, approvedParticipants)
	if err != nil {
		return fmt.Errorf("failed to generate groups: %w", err)
	}
	// 1️⃣ Fetch all group IDs for this tournament
	rows, err := r.db.Query("SELECT id FROM groups WHERE tournament_id=$1", tournament_id)
	if err != nil {
		return fmt.Errorf("failed to fetch groups: %w", err)
	}
	defer rows.Close()

	var groupIDs []int
	for rows.Next() {
		var id int
		if err := rows.Scan(&id); err != nil {
			return fmt.Errorf("failed to scan group ID: %w", err)
		}
		groupIDs = append(groupIDs, id)
	}

	if len(groupIDs) == 0 {
		return fmt.Errorf("no groups found for tournament %d", tournament_id)
	}
	// 2️⃣ Get participants per group
	participantsPerGroup := make(map[int][]int)
	for _, groupID := range groupIDs {
		rows, err := r.db.Query("SELECT participant_id FROM group_teams WHERE group_id = $1", groupID)
		if err != nil {
			return fmt.Errorf("failed to fetch participants for group %d: %w", groupID, err)
		}
		defer rows.Close()

		for rows.Next() {
			var pid int
			if err := rows.Scan(&pid); err != nil {
				return fmt.Errorf("failed to scan participant ID: %w", err)
			}
			participantsPerGroup[groupID] = append(participantsPerGroup[groupID], pid)
		}
	}

	// 3️⃣ Begin transaction
	tx, err := r.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() { _ = tx.Rollback() }()

	// 4️⃣ Prepare insert statement
	stmt, err := tx.Prepare(`
		INSERT INTO matches (tournament_id, group_id, round, participant_a_id, participant_b_id)
		VALUES ($1, $2, $3, $4, $5)
	`)
	if err != nil {
		return fmt.Errorf("failed to prepare insert statement: %w", err)
	}
	defer stmt.Close()

	// 5️⃣ Generate round-robin matches
	for groupID, participants := range participantsPerGroup {
		for i := 0; i < len(participants); i++ {
			for j := i + 1; j < len(participants); j++ {
				if _, err := stmt.Exec(tournament_id, groupID, "Group Stage", participants[i], participants[j]); err != nil {
					return fmt.Errorf("failed to insert match for group %d: %w", groupID, err)
				}
			}
		}
	}

	// 6️⃣ Commit
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}
