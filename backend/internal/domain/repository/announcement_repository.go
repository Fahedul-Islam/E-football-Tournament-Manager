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

	// Reaction operations
	ReactOnAnnouncement(ctx context.Context, tournamentID int, announcementID int, userID int, reaction string) (*domain.Announcement, error)
	GetAnnouncementPrevReaction(ctx context.Context, tournamentID int, announcementID int, userID int) (string, error)
	RemoveAnnouncementReaction(ctx context.Context, tournamentID int, announcementID int, userID int, reaction string) (*domain.Announcement, error)

	// Helper operations
	VerifyTournamentOwner(ctx context.Context, tournamentID int, userID int) (bool, error)
	GetAllParticipant(ctx context.Context, tournamentID int) ([]*domain.Participant, error)
}
