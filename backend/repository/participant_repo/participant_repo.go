package participantrepo

import (
	"database/sql"
	"time"
	"tournament-manager/domain"
	"tournament-manager/rest/handler/participant"
)

type ParticipantRepo interface {
	participant.Service
}

type participantRepo struct {
	db *sql.DB
}

func NewParticipantRepo(db *sql.DB) ParticipantRepo {
	return &participantRepo{db: db}
}

func (r *participantRepo) RequestToJoinTournament(req domain.ParticipantRequest) error {
	var total_current_participant int
	err := r.db.QueryRow(`SELECT COUNT(*) FROM participants WHERE tournament_id=$1 AND status='approved'`, req.TournamentID).Scan(&total_current_participant)
	if err!=nil {
		return  err
	}
	var max_participants int
	err = r.db.QueryRow(`SELECT max_players FROM tournaments WHERE id=$1`,req.TournamentID).Scan(&max_participants)
	if err!=nil {
		return  err
	}
	if total_current_participant >= max_participants {
		return  sql.ErrNoRows
	}
	query := `INSERT INTO participants (user_id, tournament_id, team_name, created_at) VALUES ($1, $2, $3, $4) RETURNING id`
	participant := domain.Participant{
		UserID:       req.UserID,
		TournamentID: req.TournamentID,
		TeamName:     req.TeamName,
		CreatedAt:    time.Now().Format(time.RFC3339),
	}
	return r.db.QueryRow(query, participant.UserID, participant.TournamentID, participant.TeamName, participant.CreatedAt).Scan(&participant.ID)
}
