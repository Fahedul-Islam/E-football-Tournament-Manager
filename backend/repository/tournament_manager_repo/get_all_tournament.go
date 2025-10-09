package tournamentmanagerrepo

import (
	"tournament-manager/domain"
)

func (r *tournamentManagerRepo) GetAllTournaments(tournament_owner_id int) ([]*domain.Tournament, error) {
	rows, err := r.db.Query("SELECT * FROM tournaments WHERE created_by = $1", tournament_owner_id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tournaments []*domain.Tournament
	for rows.Next() {
		t := &domain.Tournament{}
		if err := rows.Scan(&t.ID, &t.Name, &t.Description, &t.TournamentType, &t.MaxPlayers, &t.CreatedBy, &t.StartDate, &t.EndDate, &t.CreatedAt); err != nil {
			return nil, err
		}
		tournaments = append(tournaments, t)
	}
	return tournaments, nil
}
