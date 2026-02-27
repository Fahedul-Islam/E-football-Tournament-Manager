package tournament

import (
	tournamentmanager "tournament-manager/internal/delivery/http/handler/tournamentManager"
	"tournament-manager/internal/domain/repository"
)

type Service interface {
	tournamentmanager.Service
}

// TournamentService implements the tournament service layer
type TournamentRepo interface {
	repository.TournamentRepository
}

type service struct {
	tournamentRepo TournamentRepo
}

// NewTournamentService creates a new tournament service
func NewTournamentService(tournamentRepo TournamentRepo) Service {
	return &service{
		tournamentRepo: tournamentRepo,
	}
}
