package tournamentmanagerrepo

import "tournament-manager/domain"

func (r *tournamentManagerRepo) CreateTournament(created_by int, request domain.TournamentCreateRequest) error {
	query := `INSERT INTO tournaments (name, description, tournament_type, max_players,  created_by, start_date, end_date) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id`
	tournament := domain.Tournament{
		Name:           request.Name,
		Description:    request.Description,
		TournamentType: request.TournamentType,
		MaxPlayers:     request.MaxPlayers,
		StartDate:      request.StartDate,
		EndDate:        request.EndDate,
		CreatedBy:      created_by,
	}
	return r.db.QueryRow(query, tournament.Name, tournament.Description, tournament.TournamentType, tournament.MaxPlayers, tournament.CreatedBy, tournament.StartDate, tournament.EndDate).Scan(&tournament.ID)
}
