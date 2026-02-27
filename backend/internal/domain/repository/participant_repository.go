package repository

import (
	"context"
	"tournament-manager/internal/domain"
)

// ParticipantRepository defines the interface for participant data access operations
type ParticipantRepository interface {
	RequestToJoinTournament(ctx context.Context, req domain.ParticipantRequest) error
	IsApprovedParticipant(ctx context.Context, tournamentID int, userID int) (bool, error)
	GetGroupDistribution(ctx context.Context, tournamentID int) ([]*domain.Group, error)
	GetMatchSchedule(ctx context.Context, tournamentID int) ([]*domain.Match, error)
}
