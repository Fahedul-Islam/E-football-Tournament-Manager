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
	query := `INSERT INTO participants (user_id, tournament_id, team_name, created_at) VALUES ($1, $2, $3, $4) RETURNING id`
	participant := domain.Participant{
		UserID:       req.UserID,
		TournamentID: req.TournamentID,
		TeamName:     req.TeamName,
		CreatedAt:    time.Now().Format(time.RFC3339),
	}
	return r.db.QueryRow(query, participant.UserID, participant.TournamentID, participant.TeamName, participant.CreatedAt).Scan(&participant.ID)
}
