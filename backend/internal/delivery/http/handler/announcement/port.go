package announcement

import (
	"context"
	"tournament-manager/internal/domain"
)

type Service interface {
	CreateAnnouncement(ctx context.Context, tournamentID int, userID int, req domain.AnnouncementCreateRequest) (*domain.Announcement, error)
	GetAnnouncements(ctx context.Context, tournamentID int, userID int) ([]*domain.Announcement, error)
	GetAnnouncementByID(ctx context.Context, tournamentID int, announcementID int, userID int) (*domain.Announcement, error)
	UpdateAnnouncement(ctx context.Context, tournamentID int, announcementID int, userID int, req domain.AnnouncementCreateRequest) (*domain.Announcement, error)
	DeleteAnnouncement(ctx context.Context, tournamentID int, announcementID int, userID int) error
	GetParticipantsAnnouncementSeenStatus(ctx context.Context, tournamentID int, announcementID int, userID int) (*[]domain.Participant, error)
	ReactOnAnnouncement(ctx context.Context, tournamentID int, announcementID int, userID int, reaction string) (*domain.Announcement, error)
}

type AnnouncementHandler struct {
	announcementService Service
}

func NewAnnouncementHandler(announcementService Service) *AnnouncementHandler {
	return &AnnouncementHandler{
		announcementService: announcementService,
	}
}
