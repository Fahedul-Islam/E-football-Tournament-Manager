package tournamentmanagerrepo

import (
	"database/sql"
	"tournament-manager/domain"
)

func (r *tournamentManagerRepo) UpdateScore(tournament_owner_id int, req *domain.UpadateMatchScoreInput) (*domain.UpadateMatchScoreInput, error) {
	var vaildOwner bool
	err := r.db.QueryRow("SELECT EXISTS(SELECT 1 FROM tournaments WHERE id = $1 AND created_by = $2)", req.TournamentID, tournament_owner_id).Scan(&vaildOwner)
	if err != nil {
		return nil, err
	}
	if !vaildOwner {
		return nil, sql.ErrNoRows
	}
	_, err = r.db.Exec("UPDATE matches SET score_a = $1, score_b = $2, status= $3 WHERE tournament_id = $4 AND participant_a_id = $5 AND participant_b_id = $6 AND round = $7", req.ScoreA, req.ScoreB, "completed", req.TournamentID, req.ParticipantAID, req.ParticipantBID, req.Round)
	if err != nil {
		return nil, err
	}
	return req, nil
}

