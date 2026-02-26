package participant

import (
	"context"
	"tournament-manager/internal/domain"
)

type Service interface {
	RequestToJoinTournament(context.Context,domain.ParticipantRequest) error
	IsApprovedParticipant(ctx context.Context,tournament_id int, user_id int) (bool, error)
	GetGroupDistribution(ctx context.Context,tournament_id int) ([]*domain.Group, error)
	SeeMatchSchedule(ctx context.Context,tournament_id int) ([]*domain.Match, error)
}

type ParticipantHandler struct {
	service Service
}

func NewParticipantHandler(service Service) *ParticipantHandler {
	return &ParticipantHandler{service: service}
}
