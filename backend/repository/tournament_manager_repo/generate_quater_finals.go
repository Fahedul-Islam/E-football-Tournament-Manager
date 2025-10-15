package tournamentmanagerrepo

import (
	"database/sql"
	"fmt"
	"math/rand"
	"time"
)

func (r *tournamentManagerRepo) GenerateQuarterFinals(tournament_id int) (bool, error) {
	// 1️⃣ Fetch winners from Round of 16
	rows, err := r.db.Query(`
		SELECT 
			CASE 
				WHEN participant_a_score > participant_b_score THEN participant_a_id
				WHEN participant_b_score > participant_a_score THEN participant_b_id
				ELSE NULL
			END AS winner_id
		FROM matches
		WHERE tournament_id = $1 AND round = 'Round of 16' AND status = 'completed'
	`, tournament_id)
	if err != nil {
		return false, fmt.Errorf("failed to fetch round of 16 results: %w", err)
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

	// 2️⃣ Validate we have exactly 8 winners
	if len(winners) != 8 {
		return false, fmt.Errorf("expected 8 winners from Round of 16, got %d", len(winners))
	}

	// 3️⃣ Shuffle winners to randomize pairings
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(winners), func(i, j int) {
		winners[i], winners[j] = winners[j], winners[i]
	})

	// 4️⃣ Begin transaction
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

	// 5️⃣ Create 4 Quarterfinal matches
	for i := 0; i < len(winners); i += 2 {
		_, err := stmt.Exec(tournament_id, -1, "Quarterfinals", winners[i], winners[i+1])
		if err != nil {
			return false, fmt.Errorf("failed to insert quarterfinal match: %w", err)
		}
	}

	// 6️⃣ Commit
	if err := tx.Commit(); err != nil {
		return false, fmt.Errorf("failed to commit quarterfinal matches: %w", err)
	}

	return true, nil
}
