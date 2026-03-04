package repository

import (
	"context"
	"tournament-manager/internal/domain"
)

// TournamentRepository defines the interface for tournament data access operations
type TournamentRepository interface {
	CreateTournament(ctx context.Context, createdBy int, request domain.TournamentCreateRequest) error
	GetTournamentByID(ctx context.Context, id int) (*domain.Tournament, error)
	GetAllTournaments(ctx context.Context, tournamentOwnerID int) ([]*domain.Tournament, error)
	UpdateTournament(ctx context.Context, tournamentOwnerID int, tournamentID int, tournament domain.TournamentCreateRequest) error
	DeleteTournament(ctx context.Context, tournamentOwnerID int, tournamentID int) error
	ApproveParticipant(ctx context.Context, tournamentOwnerID int, req domain.ParticipantRequest) error
	RejectParticipant(ctx context.Context, tournamentOwnerID int, req domain.ParticipantRequest) error
	AddParticipant(ctx context.Context, tournamentOwnerID int, participant domain.ParticipantRequest) error
	RemoveParticipant(ctx context.Context, tournamentOwnerID int, req domain.ParticipantRequest) error
	GetAllParticipant(ctx context.Context, tournamentID int) ([]*domain.Participant, error)
	GetApprovedParticipants(ctx context.Context, tournamentID int) ([]*domain.Participant, error)
	CreateMatchSchedules(ctx context.Context, tournamentID int, groupCount int, approvedParticipants []*domain.Participant) error
	GenerateGroups(ctx context.Context, tournamentID int, groupCount int, approvedParticipants []*domain.Participant) error
	GetAllMatches(ctx context.Context, tournamentID int) ([]*domain.Match, error)
	UpdateScore(ctx context.Context, tournamentOwnerID int, req *domain.UpdateMatchScoreInput) (*domain.UpdateMatchScoreInput, error)
	CheckAndAdvanceRound(ctx context.Context, tournamentID int, round string) (bool, error)
	GetGroupCount(ctx context.Context, tournamentID int) (int, error)
	GenerateKnockoutStage(ctx context.Context, tournamentID int) (bool, error)
	GenerateQuarterFinals(ctx context.Context, tournamentID int) (bool, error)
	GenerateSemiFinals(ctx context.Context, tournamentID int) (bool, error)
	GenerateFinal(ctx context.Context, tournamentID int) (bool, error)
	GetLeaderboard(ctx context.Context, tournamentID int) (map[int][]domain.PlayerStat, error)
	GetTournamentType(ctx context.Context, tournamentID int) (string, error)
	LeagueStyleSchedule(ctx context.Context, tournamentID int, approvedParticipants []*domain.Participant) error
	VerifyTournamentOwner(ctx context.Context, tournamentID int, userID int) (bool, error)

	//All about announcement
	CreateAnnouncement(ctx context.Context, announcement *domain.Announcement) (*domain.Announcement, error)
	GetAnnouncements(ctx context.Context, tournamentID int) ([]*domain.Announcement, error)
	GetAnnouncementByID(ctx context.Context, tournamentID int, announcementID int) (*domain.Announcement, error)
	UpdateAnnouncement(ctx context.Context, announcement *domain.Announcement) (*domain.Announcement, error)
	DeleteAnnouncement(ctx context.Context, tournamentID int, announcementID int) error

}
