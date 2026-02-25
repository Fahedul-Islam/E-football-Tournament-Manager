package repository

import "tournament-manager/internal/domain"

// TournamentRepository defines the interface for tournament data access operations
type TournamentRepository interface {
	CreateTournament(createdBy int, request domain.TournamentCreateRequest) error
	GetTournamentByID(id int) (*domain.Tournament, error)
	GetAllTournaments(tournamentOwnerID int) ([]*domain.Tournament, error)
	UpdateTournament(tournamentOwnerID int, tournamentID int, tournament domain.TournamentCreateRequest) error
	DeleteTournament(tournamentOwnerID int, tournamentID int) error
	ApproveParticipant(tournamentOwnerID int, req domain.ParticipantRequest) error
	RejectParticipant(tournamentOwnerID int, req domain.ParticipantRequest) error
	AddParticipant(tournamentOwnerID int, participant domain.ParticipantRequest) error
	RemoveParticipant(tournamentOwnerID int, req domain.ParticipantRequest) error
	GetAllParticipant(tournamentID int) ([]*domain.Participant, error)
	GetApprovedParticipants(tournamentID int) ([]*domain.Participant, error)
	CreateMatchSchedules(tournamentID int, groupCount int, approvedParticipants []*domain.Participant) error
	GenerateGroups(tournamentID int, groupCount int, approvedParticipants []*domain.Participant) error
	GetAllMatches(tournamentID int) ([]*domain.Match, error)
	UpdateScore(tournamentOwnerID int, req *domain.UpadateMatchScoreInput) (*domain.UpadateMatchScoreInput, error)
	CheckAndAdvanceRound(tournamentID int, round string) (bool, error)
	GetGroupCount(tournamentID int) (int, error)
	GenerateKnockoutStage(tournamentID int) (bool, error)
	GenerateQuarterFinals(tournamentID int) (bool, error)
	GenerateSemiFinals(tournamentID int) (bool, error)
	GenerateFinal(tournamentID int) (bool, error)
	GetLeaderboard(tournamentID int) (map[int][]domain.PlayerStat, error)
	GetTournamentType(tournamentID int) (string, error)
	LeagueStyleSchedule(tournamentID int, approvedParticipants []*domain.Participant) error
	VerifyTournamentOwner(tournamentID int, userID int) (bool, error)
}
