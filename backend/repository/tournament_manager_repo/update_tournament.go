package tournamentmanagerrepo

import "tournament-manager/domain"

func (r *tournamentManagerRepo) UpdateTournament(tournament_owner_id int, tournament_id int, tournament domain.TournamentCreateRequest) error {
	_, err := r.db.Exec("UPDATE tournaments SET name = $1, description = $2, start_date = $3, end_date = $4 WHERE id = $5 AND created_by = $6",
		tournament.Name, tournament.Description, tournament.StartDate, tournament.EndDate, tournament_id, tournament_owner_id)
	return err
}
