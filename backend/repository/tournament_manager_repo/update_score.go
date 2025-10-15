package tournamentmanagerrepo

import (
	"database/sql"
	"fmt"
	"tournament-manager/domain"
)

func (r *tournamentManagerRepo) UpdateScore(tournament_owner_id int, req *domain.UpadateMatchScoreInput) (*domain.UpadateMatchScoreInput, error) {
	// 1️⃣ Validate tournament ownership
	var validOwner bool
	err := r.db.QueryRow(`
		SELECT EXISTS(SELECT 1 FROM tournaments WHERE id = $1 AND created_by = $2)
	`, req.TournamentID, tournament_owner_id).Scan(&validOwner)
	if err != nil {
		return nil, fmt.Errorf("failed to verify tournament ownership: %w", err)
	}
	if !validOwner {
		return nil, sql.ErrNoRows
	}

	// 2️⃣ Update the match record
	_, err = r.db.Exec(`
		UPDATE matches 
		SET participant_a_score = $1, participant_b_score = $2, status = $3
		WHERE tournament_id = $4 AND participant_a_id = $5 AND participant_b_id = $6 AND round = $7
	`, req.ScoreA, req.ScoreB, "completed", req.TournamentID, req.ParticipantAID, req.ParticipantBID, req.Round)
	if err != nil {
		return nil, fmt.Errorf("failed to update match: %w", err)
	}

	// 3️⃣ Get group_id of this match
	var groupID int
	err = r.db.QueryRow(`
		SELECT group_id FROM matches
		WHERE tournament_id = $1 AND participant_a_id = $2 AND participant_b_id = $3 AND round = $4
	`, req.TournamentID, req.ParticipantAID, req.ParticipantBID, req.Round).Scan(&groupID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch group ID: %w", err)
	}

	// 4️⃣ Begin transaction for player_stats update
	tx, err := r.db.Begin()
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() { _ = tx.Rollback() }()

	// 5️⃣ Update stats for both participants
	aGoals, bGoals := req.ScoreA, req.ScoreB
	aPoints, bPoints := 0, 0
	aWin, bWin, aDraw, bDraw, aLoss, bLoss := 0, 0, 0, 0, 0, 0

	switch {
	case aGoals > bGoals:
		aPoints, bPoints = 3, 0
		aWin, bLoss = 1, 1
	case aGoals < bGoals:
		aPoints, bPoints = 0, 3
		aLoss, bWin = 1, 1
	default:
		aPoints, bPoints = 1, 1
		aDraw, bDraw = 1, 1
	}

	// Participant A stats update
	_, err = tx.Exec(`
		UPDATE player_stats
		SET 
			matches_played = matches_played + 1,
			wins = wins + $1,
			draws = draws + $2,
			losses = losses + $3,
			goals_scored = goals_scored + $4,
			goals_conceded = goals_conceded + $5,
			goal_difference = goal_difference + ($4 - $5),
			points = points + $6
		WHERE group_id = $7 AND participant_id = $8
	`, aWin, aDraw, aLoss, aGoals, bGoals, aPoints, groupID, req.ParticipantAID)
	if err != nil {
		return nil, fmt.Errorf("failed to update player A stats: %w", err)
	}

	// Participant B stats update
	_, err = tx.Exec(`
		UPDATE player_stats
		SET 
			matches_played = matches_played + 1,
			wins = wins + $1,
			draws = draws + $2,
			losses = losses + $3,
			goals_scored = goals_scored + $4,
			goals_conceded = goals_conceded + $5,
			goal_difference = goal_difference + ($4 - $5),
			points = points + $6
		WHERE group_id = $7 AND participant_id = $8
	`, bWin, bDraw, bLoss, bGoals, aGoals, bPoints, groupID, req.ParticipantBID)
	if err != nil {
		return nil, fmt.Errorf("failed to update player B stats: %w", err)
	}

	// 6️⃣ Commit transaction
	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit player stats update: %w", err)
	}

	return req, nil
}
