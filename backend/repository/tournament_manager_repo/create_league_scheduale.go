package tournamentmanagerrepo

import (
	"fmt"
	"math/rand"
	"time"
	"tournament-manager/domain"
)

func (r *tournamentManagerRepo) LeagueStyleSchedule(tournament_id int, approvedParticipants []*domain.Participant) error {
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
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(approvedParticipants), func(i, j int) {
		approvedParticipants[i], approvedParticipants[j] = approvedParticipants[j], approvedParticipants[i]
	})

	for i, participantA := range approvedParticipants {
		playerStatStmt, err := tx.Prepare(`
			INSERT INTO player_stats participant_id VALUES ($1)
		`)
		if err != nil {
			return err
		}
		defer playerStatStmt.Close()
		if _, err := playerStatStmt.Exec(participantA.ID); err != nil {
			return fmt.Errorf("fail to insert player stat")
		}
		for j, participantB := range approvedParticipants {
			if i != j {
				matchStmt, err := tx.Prepare(`
				INSERT INTO matches (tournament_id, group_id, round, participant_a_id, participant_b_id)
				VALUES ($1, $2, $3, $4, $5)`)
				if err != nil {
					return fmt.Errorf("failed to prepare match insert: %w", err)
				}
				defer matchStmt.Close()
				if _,err:= matchStmt.Exec(tournament_id, nil, "leauge",participantA.ID, participantB.ID); err!=nil{
					return fmt.Errorf("failed to insert matches")
				}
			}
		}
	}
	return nil
}
