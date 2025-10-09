package tournamentmanagerrepo

import (
	"database/sql"
	"tournament-manager/domain"
	tournamentmanager "tournament-manager/rest/handler/tournamentManager"
)

type TournamentManagerRepo interface {
	tournamentmanager.Service
}

type tournamentManagerRepo struct {
	db *sql.DB
}

func NewTournamentManagerRepo(db *sql.DB) TournamentManagerRepo {
	return &tournamentManagerRepo{db: db}
}

func (r *tournamentManagerRepo) GetTournamentByID(id int) (*domain.Tournament, error) {
	var tournament domain.Tournament
	query := `SELECT * FROM tournaments WHERE id = $1`
	if err := r.db.QueryRow(query, id).Scan(&tournament.ID, &tournament.Name, &tournament.Description, &tournament.StartDate, &tournament.EndDate, &tournament.CreatedBy); err != nil {
		return nil, err
	}
	return &tournament, nil
}
