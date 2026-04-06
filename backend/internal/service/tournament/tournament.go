package tournament

import (
	"context"
	"errors"
	"tournament-manager/internal/domain"
)

// CreateTournament creates a new tournament
func (s *service) CreateTournament(ctx context.Context, createdBy int, req domain.TournamentCreateRequest) error {
	if req.TournamentType != "knockout" && req.TournamentType != "league" && req.TournamentType != "group+knockout" {
		return errors.New("invalid tournament type")
	}
	if req.MaxPlayers < 3 {
		return errors.New("not enough players")
	}
	if req.MaxPlayers > 64 {
		return errors.New("too many players")
	}
	if req.Description == "" || req.Name == "" || len(req.Description) < 50 {
		return errors.New("name and description must be filled, and description must be at least 50 characters long")
	}
	return s.tournamentRepo.CreateTournament(ctx, createdBy, req)
}

// GetTournamentByID returns a tournament by ID
func (s *service) GetTournamentByID(ctx context.Context, id int) (*domain.Tournament, error) {
	return s.tournamentRepo.GetTournamentByID(ctx, id)
}

// GetAllTournaments returns all tournaments for a given owner
func (s *service) GetAllTournaments(ctx context.Context, tournamentOwnerID int) ([]*domain.Tournament, error) {
	return s.tournamentRepo.GetAllTournaments(ctx, tournamentOwnerID)
}

// UpdateTournament updates a tournament
func (s *service) UpdateTournament(ctx context.Context, tournamentOwnerID int, tournamentID int, req domain.TournamentCreateRequest) error {
	if req.TournamentType != "knockout" && req.TournamentType != "league" && req.TournamentType != "group+knockout" {
		return errors.New("invalid tournament type")
	}
	if req.MaxPlayers < 3 {
		return errors.New("not enough players")
	}
	if req.MaxPlayers > 64 {
		return errors.New("too many players")
	}
	if req.Description == "" || req.Name == "" || len(req.Description) < 50 {
		return errors.New("name and description must be filled, and description must be at least 50 characters long")
	}
	return s.tournamentRepo.UpdateTournament(ctx, tournamentOwnerID, tournamentID, req)
}

// DeleteTournament deletes a tournament
func (s *service) DeleteTournament(ctx context.Context, tournamentOwnerID int, tournamentID int) error {
	return s.tournamentRepo.DeleteTournament(ctx, tournamentOwnerID, tournamentID)
}
