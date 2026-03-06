package announcement

import (
	"context"
	"tournament-manager/internal/domain"
)

type Service interface {
	// CRUD operations for announcements
	CreateAnnouncement(ctx context.Context, tournamentID int, userID int, req domain.AnnouncementCreateRequest) (*domain.Announcement, error)
	GetAnnouncements(ctx context.Context, tournamentID int, userID int) ([]*domain.Announcement, error)
	GetAnnouncementByID(ctx context.Context, tournamentID int, announcementID int, userID int) (*domain.Announcement, error)
	UpdateAnnouncement(ctx context.Context, tournamentID int, announcementID int, userID int, req domain.AnnouncementCreateRequest) (*domain.Announcement, error)
	DeleteAnnouncement(ctx context.Context, tournamentID int, announcementID int, userID int) error
	GetParticipantsAnnouncementSeenStatus(ctx context.Context, tournamentID int, announcementID int, userID int) (*[]domain.Participant, error)
	ReactOnAnnouncement(ctx context.Context, tournamentID int, announcementID int, userID int, reaction string) (*domain.Announcement, error)

	// Announcement Comment operations
	CreateComment(ctx context.Context, tournamentID int, announcementID int, userID int, req domain.CommentCreateRequest) (*domain.AnnouncementComment, error)
	GetComments(ctx context.Context, tournamentID int, userID int, parentCommentID *int, announcementID int) ([]*domain.AnnouncementComment, error)
	DeleteComment(ctx context.Context, tournamentID int, userID int, commentID int) error
	EditComment(ctx context.Context, tournamentID int, userID int, commentID int, req domain.CommentCreateRequest) (*domain.AnnouncementComment, error)
	ReactToComment(ctx context.Context, tournamentID int, commentID int, userID int, reaction string) (*domain.AnnouncementComment, error)
}

type AnnouncementHandler struct {
	announcementService Service
}

func NewAnnouncementHandler(announcementService Service) *AnnouncementHandler {
	return &AnnouncementHandler{
		announcementService: announcementService,
	}
}
