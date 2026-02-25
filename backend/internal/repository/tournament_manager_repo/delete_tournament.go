package tournamentmanagerrepo

import "database/sql"

func (r *tournamentManagerRepo) DeleteTournament(tournament_owner_id int, id int) error {
	var exists bool
	err := r.db.QueryRow("SELECT EXISTS(SELECT 1 FROM tournaments WHERE id = $1 AND created_by = $2)", id, tournament_owner_id).Scan(&exists)
	if err != nil {
		return err
	}
	if !exists {
		return sql.ErrNoRows
	}
	_, err = r.db.Exec("DELETE FROM tournaments WHERE id = $1", id)
	return err
}
