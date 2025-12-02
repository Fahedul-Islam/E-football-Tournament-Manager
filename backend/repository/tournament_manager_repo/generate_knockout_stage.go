package tournamentmanagerrepo

import (
	"fmt"
	"math/rand"
	"time"
)

func (r *tournamentManagerRepo) GenerateKnockoutStage(tournament_id int) (bool, error) {
	// 1️⃣ Fetch all groups for this tournament
	groupRows, err := r.db.Query(`SELECT id FROM groups WHERE tournament_id = $1 ORDER BY id`, tournament_id)
	if err != nil {
		return false, fmt.Errorf("failed to fetch groups: %w", err)
	}
	defer groupRows.Close()

	var groupIDs []int
	for groupRows.Next() {
		var id int
		if err := groupRows.Scan(&id); err != nil {
			return false, fmt.Errorf("failed to scan group ID: %w", err)
		}
		groupIDs = append(groupIDs, id)
	}
	if len(groupIDs) == 0 {
		return false, fmt.Errorf("no groups found for tournament %d", tournament_id)
	}

	// 2️⃣ Get top 2 participants from each group (sorted by points, GD, goals scored)
	var qualifiedPlayers []int
	for _, gid := range groupIDs {
		rows, err := r.db.Query(`
			SELECT participant_id 
			FROM player_stats 
			WHERE group_id = $1 
			ORDER BY points DESC, goal_difference DESC, goals_scored DESC 
			LIMIT 2
		`, gid)
		if err != nil {
			return false, fmt.Errorf("failed to fetch top participants for group %d: %w", gid, err)
		}

		for rows.Next() {
			var pid int
			if err := rows.Scan(&pid); err != nil {
				rows.Close()
				return false, fmt.Errorf("failed to scan participant ID: %w", err)
			}
			qualifiedPlayers = append(qualifiedPlayers, pid)
		}
		rows.Close()
	}
	var next_round string
	if len(qualifiedPlayers) == 16 {
		next_round = "Round of 16"
	}else if len(qualifiedPlayers) == 8 {
		next_round = "Quarterfinals"
	} else if len(qualifiedPlayers) == 4 {
		next_round = "Semifinals"
	} else if len(qualifiedPlayers) == 2 {
		next_round = "Final"
	} else {
		return false, fmt.Errorf("expected 16, 8, 4 or 2 qualified participants, got %d", len(qualifiedPlayers))
	}

	// 3️⃣ Randomize the qualified participants for fairness
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(qualifiedPlayers), func(i, j int) {
		qualifiedPlayers[i], qualifiedPlayers[j] = qualifiedPlayers[j], qualifiedPlayers[i]
	})

	// 4️⃣ Begin transaction
	tx, err := r.db.Begin()
	if err != nil {
		return false, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() { _ = tx.Rollback() }()

	// 5️⃣ Prepare insert for Round of 16
	stmt, err := tx.Prepare(`
		INSERT INTO matches (tournament_id, group_id, round, participant_a_id, participant_b_id, status)
		VALUES ($1, $2, $3, $4, $5, 'scheduled')
	`)
	if err != nil {
		return false, fmt.Errorf("failed to prepare insert: %w", err)
	}
	defer stmt.Close()

	// 6️⃣ Create 8 Round of 16 matches (pairing players)
	for i := 0; i < len(qualifiedPlayers); i += 2 {
		a := qualifiedPlayers[i]
		b := qualifiedPlayers[i+1]
		if _, err := stmt.Exec(tournament_id, nil, next_round, a, b); err != nil {
			return false, fmt.Errorf("failed to insert knockout-stage match: %w", err)
		}
	}

	// 7️⃣ Commit
	if err := tx.Commit(); err != nil {
		return false, fmt.Errorf("failed to commit knockout-stage matches: %w", err)
	}

	return true, nil
}
