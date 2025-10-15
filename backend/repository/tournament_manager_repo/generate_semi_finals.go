package tournamentmanagerrepo

import (
	"database/sql"
	"fmt"
	"math/rand"
	"time"
)

func (r *tournamentManagerRepo) GenerateSemiFinals(tournament_id int) (bool, error) {
	rows, err := r.db.Query(`
		SELECT 
			CASE 
				WHEN participant_a_score > participant_b_score THEN participant_a_id
				WHEN participant_b_score > participant_a_score THEN participant_b_id
				ELSE NULL
			END AS winner_id
		FROM matches
		WHERE tournament_id = $1 AND round = 'Quarterfinals' AND status = 'completed'
	`, tournament_id)
	if err != nil {
		return false, fmt.Errorf("failed to fetch quarterfinal results: %w", err)
	}
	defer rows.Close()

	var winners []int
	for rows.Next() {
		var winner sql.NullInt64
		if err := rows.Scan(&winner); err != nil {
			return false, fmt.Errorf("failed to scan winner: %w", err)
		}
		if winner.Valid {
			winners = append(winners, int(winner.Int64))
		}
	}

	if len(winners) != 4 {
		return false, fmt.Errorf("expected 4 winners from Quarterfinals, got %d", len(winners))
	}

	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(winners), func(i, j int) {
		winners[i], winners[j] = winners[j], winners[i]
	})

	tx, err := r.db.Begin()
	if err != nil {
		return false, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() { _ = tx.Rollback() }()

	stmt, err := tx.Prepare(`
		INSERT INTO matches (tournament_id, group_id, round, participant_a_id, participant_b_id)
		VALUES ($1, $2, $3, $4, $5)
	`)
	if err != nil {
		return false, fmt.Errorf("failed to prepare insert statement: %w", err)
	}
	defer stmt.Close()

	for i := 0; i < len(winners); i += 2 {
		_, err := stmt.Exec(tournament_id, -1, "Semifinals", winners[i], winners[i+1])
		if err != nil {
			return false, fmt.Errorf("failed to insert semifinal match: %w", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return false, fmt.Errorf("failed to commit semifinal matches: %w", err)
	}

	return true, nil
}
