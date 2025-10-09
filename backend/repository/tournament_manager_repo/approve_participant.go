package tournamentmanagerrepo

import (
	"database/sql"
	"tournament-manager/domain"
)

func (r *tournamentManagerRepo) ApproveParticipant(tournament_owner_id int, req domain.ParticipantRequest) error {
	// Check if the tournament exists and is created by the tournament_owner_id
	var exists bool
	err := r.db.QueryRow("SELECT EXISTS(SELECT 1 FROM tournaments WHERE id = $1 AND created_by = $2)", req.TournamentID, tournament_owner_id).Scan(&exists)
	if err != nil {
		return err
	}
	if !exists {
		return sql.ErrNoRows
	}
	// check if maximum participants reached
	var count int
	err = r.db.QueryRow("SELECT COUNT(*) FROM participants WHERE tournament_id = $1 AND status = 'approved'", req.TournamentID).Scan(&count)
	if err != nil {
		return err
	}

	var maxParticipants int
	err = r.db.QueryRow("SELECT max_players FROM tournaments WHERE id = $1", req.TournamentID).Scan(&maxParticipants)
	if err != nil {
		return err
	}

	if count >= maxParticipants {
		return sql.ErrNoRows
	}
	_, err = r.db.Exec("UPDATE participants SET status = 'approved' WHERE user_id = $1 AND tournament_id = $2", req.UserID, req.TournamentID)
	return err
}
