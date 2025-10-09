package tournamentmanagerrepo

import (
	"math/rand"
	"time"
	"tournament-manager/domain"
)

func (r *tournamentManagerRepo) GenerateGroups(tournament_id int, groupCount int, approvedParticipants []*domain.Participant) error {
	// Start a transaction
	tx, err := r.db.Begin()
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

	// Step 1: Clear existing groups and group_teams for the tournament
	_, err = tx.Exec("DELETE FROM group_teams WHERE group_id IN (SELECT id FROM groups WHERE tournament_id = $1)", tournament_id)
	if err != nil {
		return err
	}
	_, err = tx.Exec("DELETE FROM groups WHERE tournament_id = $1", tournament_id)
	if err != nil {
		return err
	}

	// Step 2: Create groups
	groupIDs := make([]int, groupCount)
	for i := 0; i < groupCount; i++ {
		groupName := string(rune('A' + i)) // Group names: A, B, C, ...
		var groupID int
		err = tx.QueryRow("INSERT INTO groups (name, tournament_id) VALUES ($1, $2) RETURNING id", groupName, tournament_id).Scan(&groupID)
		if err != nil {
			return err
		}
		groupIDs[i] = groupID
	}
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(approvedParticipants), func(i, j int) {
		approvedParticipants[i], approvedParticipants[j] = approvedParticipants[j], approvedParticipants[i]
	})
	// Step 3: Distribute participants into groups
	for i, participant := range approvedParticipants {
		groupID := groupIDs[i%groupCount]
		_, err = tx.Exec("INSERT INTO group_teams (group_id, participant_id) VALUES ($1, $2)", groupID, participant.ID)
		if err != nil {
			return err
		}
	}

	return nil
}
