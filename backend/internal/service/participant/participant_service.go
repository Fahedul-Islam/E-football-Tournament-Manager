package participant

import (
	"context"
	"errors"
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

// GetGroupDistribution returns the group distribution for a tournament
func (s *service) DistributeGroup(ctx context.Context, tournamentID int, user_id int) ([]*domain.Group, error) {
	isApproved, err := s.participantRepo.IsApprovedParticipant(ctx, tournamentID, user_id)
	if err != nil {
		return nil, err
	}
	if !isApproved {
		return nil, errors.New("user is not an approved participant")
	}
	return s.participantRepo.GetGroupDistribution(ctx, tournamentID)
}

// SeeMatchSchedule returns the match schedule for a tournament
func (s *service) MatchSchedule(ctx context.Context, tournamentID int, user_id int) ([]*domain.Match, error) {
	isApproved, err := s.participantRepo.IsApprovedParticipant(ctx, tournamentID, user_id)
	if err != nil {
		return nil, err
	}
	if !isApproved {
		return nil, errors.New("user is not an approved participant")
	}
	return s.participantRepo.GetMatchSchedule(ctx, tournamentID)
}
