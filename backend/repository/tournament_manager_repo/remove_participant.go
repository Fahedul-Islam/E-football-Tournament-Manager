package tournamentmanagerrepo

import (
	"database/sql"
	"tournament-manager/domain"
)

func (r *tournamentManagerRepo) RemoveParticipant(tournament_owner_id int, req domain.ParticipantRequest) error {
	// Check if the tournament exists and is created by the tournament_owner_id
	var exists bool
	err := r.db.QueryRow("SELECT EXISTS(SELECT 1 FROM tournaments WHERE id = $1 AND created_by = $2)", req.TournamentID, tournament_owner_id).Scan(&exists)
	if err != nil {
		return err
	}
	if !exists {
		return sql.ErrNoRows
	}
	_, err = r.db.Exec("DELETE FROM participants WHERE user_id = $1 AND tournament_id = $2", req.UserID, req.TournamentID)
	return err
}
