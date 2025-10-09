package tournamentmanagerrepo

import (
	"database/sql"
	"time"
	"tournament-manager/domain"
)

func (r *tournamentManagerRepo) AddParticipant(tournament_owner_id int, participant domain.ParticipantRequest) error {
	// Check if the tournament exists and is created by the tournament_owner_id
	var exists bool
	err := r.db.QueryRow("SELECT EXISTS(SELECT 1 FROM tournaments WHERE id = $1 AND created_by = $2)", participant.TournamentID, tournament_owner_id).Scan(&exists)
	if err != nil {
		return err
	}
	if !exists {
		return sql.ErrNoRows
	}
	// check if maximum participants reached
	var count int
	err = r.db.QueryRow("SELECT COUNT(*) FROM participants WHERE tournament_id = $1 AND status = 'approved'", participant.TournamentID).Scan(&count)
	if err != nil {
		return err
	}

	var maxParticipants int
	err = r.db.QueryRow("SELECT max_players FROM tournaments WHERE id = $1", participant.TournamentID).Scan(&maxParticipants)
	if err != nil {
		return err
	}

	if count >= maxParticipants {
		return sql.ErrNoRows
	}
	// Add participant with status 'approved'
	now := time.Now().Format(time.RFC3339)
	addedParticipant := domain.Participant{
		UserID:       participant.UserID,
		TournamentID: participant.TournamentID,
		TeamName:     participant.TeamName,
		Status:       "approved",
		CreatedAt:    now,
	}
	query := `INSERT INTO participants (user_id, tournament_id, team_name, status, created_at) VALUES ($1, $2, $3, $4, $5) RETURNING id`
	return r.db.QueryRow(query, addedParticipant.UserID, addedParticipant.TournamentID, addedParticipant.TeamName, addedParticipant.Status, addedParticipant.CreatedAt).Scan(&addedParticipant.ID)
}
