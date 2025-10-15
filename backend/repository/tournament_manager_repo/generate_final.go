package tournamentmanagerrepo

import (
	"database/sql"
	"fmt"
)

func (r *tournamentManagerRepo) GenerateFinal(tournament_id int) (bool, error) {
	rows, err := r.db.Query(`
		SELECT 
			CASE 
				WHEN participant_a_score > participant_b_score THEN participant_a_id
				WHEN participant_b_score > participant_a_score THEN participant_b_id
				ELSE NULL
			END AS winner_id
		FROM matches
		WHERE tournament_id = $1 AND round = 'Semifinals' AND status = 'completed'
	`, tournament_id)
	if err != nil {
		return false, fmt.Errorf("failed to fetch semifinal results: %w", err)
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

	if len(winners) != 2 {
		return false, fmt.Errorf("expected 2 winners from Semifinals, got %d", len(winners))
	}

	_, err = r.db.Exec(`
		INSERT INTO matches (tournament_id, group_id, round, participant_a_id, participant_b_id)
		VALUES ($1, $2, $3, $4, $5)
	`, tournament_id, -1, "Final", winners[0], winners[1])
	if err != nil {
		return false, fmt.Errorf("failed to insert final match: %w", err)
	}

	return true, nil
}
