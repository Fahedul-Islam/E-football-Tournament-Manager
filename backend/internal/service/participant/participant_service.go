package participant

import (
	"context"
	"tournament-manager/internal/delivery/http/handler/participant"
	"tournament-manager/internal/domain"
	"tournament-manager/internal/domain/repository"
)

type Service interface {
	participant.Service
}

// ParticipantRepo implements the participant service layer
type ParticipantRepo interface {
	repository.ParticipantRepository
}

type service struct {
	participantRepo ParticipantRepo
}

// NewParticipantService creates a new participant service
func NewParticipantService(participantRepo ParticipantRepo) Service {
	return &service{
		participantRepo: participantRepo,
	}
}

// RequestToJoinTournament requests to join a tournament
func (s *service) RequestToJoinTournament(ctx context.Context, req domain.ParticipantRequest) error {
	return s.participantRepo.RequestToJoinTournament(ctx, req)
}

// IsApprovedParticipant checks if a user is an approved participant
func (s *service) IsApprovedParticipant(ctx context.Context, tournamentID int, userID int) (bool, error) {
	return s.participantRepo.IsApprovedParticipant(ctx, tournamentID, userID)
}

// GetGroupDistribution returns the group distribution for a tournament
func (s *service) GetGroupDistribution(ctx context.Context, tournamentID int) ([]*domain.Group, error) {
	return s.participantRepo.GetGroupDistribution(ctx, tournamentID)
}

// SeeMatchSchedule returns the match schedule for a tournament
func (s *service) SeeMatchSchedule(ctx context.Context, tournamentID int) ([]*domain.Match, error) {
	return s.participantRepo.SeeMatchSchedule(ctx, tournamentID)
}
