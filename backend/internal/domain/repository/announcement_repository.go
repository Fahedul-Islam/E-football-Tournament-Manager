package repository

import (
	"context"
	"tournament-manager/internal/domain"
)

// AnnouncementRepository defines the interface for announcement data access operations
type AnnouncementRepository interface {
	// CRUD operations for announcements
	CreateAnnouncement(ctx context.Context, announcement *domain.Announcement) (*domain.Announcement, error)
	GetAnnouncements(ctx context.Context, tournamentID int) ([]*domain.Announcement, error)
	GetAnnouncementByID(ctx context.Context, tournamentID int, announcementID int, userID int) (*domain.Announcement, error)
	UpdateAnnouncement(ctx context.Context, announcement *domain.Announcement) (*domain.Announcement, error)
	DeleteAnnouncement(ctx context.Context, tournamentID int, announcementID int) error
	GetParticipantsAnnouncementSeenStatus(ctx context.Context, tournamentID int, announcementID int, userID int) (*[]domain.Participant, error)

	// Reaction On Announcement operations
	ReactOnAnnouncement(ctx context.Context, tournamentID int, announcementID int, userID int, reaction string) (*domain.Announcement, error)
	RemoveAnnouncementReaction(ctx context.Context, tournamentID int, announcementID int, userID int, reaction string) (*domain.Announcement, error)

	// Announcement Comment operations
	AddComment(ctx context.Context, announcementID int, userID int, parentCommentID *int, content *string) (*domain.AnnouncementComment, error)
	GetComments(ctx context.Context, announcementID int, parentCommentID *int) ([]*domain.AnnouncementComment, error)
	DeleteComment(ctx context.Context, commentID int) error
	EditComment(ctx context.Context, commentID int, content string) (*domain.AnnouncementComment, error)
	ReactToComment(ctx context.Context, commentID int, userID int, reactionType string) (*domain.AnnouncementComment, error)
	RemoveReactionFromComment(ctx context.Context, commentID int, userID int) (*domain.AnnouncementComment, error)
	GetParentCommentUserID(ctx context.Context, parentCommentID *int) (int, error)

	//notification operations
	AddAnnouncementNotification(ctx context.Context, announcementID int, message string, participants []*domain.Participant) error
	GetNotifications(ctx context.Context, userID int, page int) ([]*domain.Notification, error)
	MarkNotificationAsRead(ctx context.Context, notificationID int, userID int) error
	MarkAllNotificationsAsRead(ctx context.Context, userID int) error
	AddCommentNotification(ctx context.Context, userID int, notification_type string, commentID int, message string) error

	// Helper operations
	VerifyTournamentOwner(ctx context.Context, tournamentID int, userID int) (bool, error)
	GetAllParticipant(ctx context.Context, tournamentID int) ([]*domain.Participant, error)
	VerifyCommentOwner(ctx context.Context, commentID int, userID int) (bool, error)
	VerifyCommentBelongsToTournament(ctx context.Context, tournamentID int, commentID int) (bool, error)
	GetCommentPrevReaction(ctx context.Context, tournamentID int, commentID int, userID int) (string, error)
	GetAnnouncementPrevReaction(ctx context.Context, tournamentID int, announcementID int, userID int) (string, error)
}
