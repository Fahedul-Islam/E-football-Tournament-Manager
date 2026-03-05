package participant

import (
	"context"
	"tournament-manager/internal/domain"
)

type Service interface {
	RequestToJoinTournament(context.Context, domain.ParticipantRequest) error
	DistributeGroup(ctx context.Context, tournament_id int, user_id int) ([]*domain.Group, error)
	MatchSchedule(ctx context.Context, tournament_id int, user_id int) ([]*domain.Match, error)
	ReactOnAnnouncement(ctx context.Context, tournamentID int, announcementID int, userID int, reaction string) (*domain.Announcement, error)
}

type ParticipantHandler struct {
	service Service
}

func NewParticipantHandler(service Service) *ParticipantHandler {
	return &ParticipantHandler{service: service}
}
