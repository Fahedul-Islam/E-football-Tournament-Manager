package tournamentmanagerrepo

import (
	"database/sql"
	"tournament-manager/internal/domain/repository"
)

type tournamentManagerRepo struct {
	db *sql.DB
}

func NewTournamentManagerRepo(db *sql.DB) repository.TournamentRepository {
	return &tournamentManagerRepo{db: db}
}
