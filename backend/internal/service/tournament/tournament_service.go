package tournament

import (
	"context"
	tournamentmanager "tournament-manager/internal/delivery/http/handler/tournamentManager"
	"tournament-manager/internal/domain"
	"tournament-manager/internal/domain/repository"
)

type Service interface {
	tournamentmanager.Service
}

// TournamentService implements the tournament service layer
type TournamentRepo interface {
	repository.TournamentRepository
}

type service struct {
	tournamentRepo TournamentRepo
}

// NewTournamentService creates a new tournament service
func NewTournamentService(tournamentRepo TournamentRepo) Service {
	return &service{
		tournamentRepo: tournamentRepo,
	}
}

// CreateTournament creates a new tournament
func (s *service) CreateTournament(ctx context.Context, createdBy int, request domain.TournamentCreateRequest) error {
	return s.tournamentRepo.CreateTournament(ctx, createdBy, request)
}

// GetTournamentByID returns a tournament by ID
func (s *service) GetTournamentByID(ctx context.Context, id int) (*domain.Tournament, error) {
	return s.tournamentRepo.GetTournamentByID(ctx, id)
}

// GetAllTournaments returns all tournaments for a given owner
func (s *service) GetAllTournaments(ctx context.Context, tournamentOwnerID int) ([]*domain.Tournament, error) {
	return s.tournamentRepo.GetAllTournaments(ctx, tournamentOwnerID)
}

// UpdateTournament updates a tournament
func (s *service) UpdateTournament(ctx context.Context, tournamentOwnerID int, tournamentID int, tournament domain.TournamentCreateRequest) error {
	return s.tournamentRepo.UpdateTournament(ctx, tournamentOwnerID, tournamentID, tournament)
}

// DeleteTournament deletes a tournament
func (s *service) DeleteTournament(ctx context.Context, tournamentOwnerID int, tournamentID int) error {
	return s.tournamentRepo.DeleteTournament(ctx, tournamentOwnerID, tournamentID)
}

// ApproveParticipant approves a participant for a tournament
func (s *service) ApproveParticipant(ctx context.Context, tournamentOwnerID int, req domain.ParticipantRequest) error {
	return s.tournamentRepo.ApproveParticipant(ctx, tournamentOwnerID, req)
}

// RejectParticipant rejects a participant from a tournament
func (s *service) RejectParticipant(ctx context.Context, tournamentOwnerID int, req domain.ParticipantRequest) error {
	return s.tournamentRepo.RejectParticipant(ctx, tournamentOwnerID, req)
}

// AddParticipant adds a participant to a tournament
func (s *service) AddParticipant(ctx context.Context, tournamentOwnerID int, participant domain.ParticipantRequest) error {
	return s.tournamentRepo.AddParticipant(ctx, tournamentOwnerID, participant)
}

// RemoveParticipant removes a participant from a tournament
func (s *service) RemoveParticipant(ctx context.Context, tournamentOwnerID int, req domain.ParticipantRequest) error {
	return s.tournamentRepo.RemoveParticipant(ctx, tournamentOwnerID, req)
}

// GetAllParticipant returns all participants for a tournament
func (s *service) GetAllParticipant(ctx context.Context, tournamentID int) ([]*domain.Participant, error) {
	return s.tournamentRepo.GetAllParticipant(ctx, tournamentID)
}

// GetApprovedParticipants returns all approved participants for a tournament
func (s *service) GetApprovedParticipants(ctx context.Context, tournamentID int) ([]*domain.Participant, error) {
	return s.tournamentRepo.GetApprovedParticipants(ctx, tournamentID)
}

// CreateMatchSchedules creates match schedules for a tournament
func (s *service) CreateMatchSchedules(ctx context.Context, tournamentID int, groupCount int, approvedParticipants []*domain.Participant) error {
	return s.tournamentRepo.CreateMatchSchedules(ctx, tournamentID, groupCount, approvedParticipants)
}

// GenerateGroups generates groups for a tournament
func (s *service) GenerateGroups(ctx context.Context, tournamentID int, groupCount int, approvedParticipants []*domain.Participant) error {
	return s.tournamentRepo.GenerateGroups(ctx, tournamentID, groupCount, approvedParticipants)
}

// GetAllMatches returns all matches for a tournament
func (s *service) GetAllMatches(ctx context.Context, tournamentID int) ([]*domain.Match, error) {
	return s.tournamentRepo.GetAllMatches(ctx, tournamentID)
}

// UpdateScore updates the score for a match
func (s *service) UpdateScore(ctx context.Context, tournamentOwnerID int, req *domain.UpadateMatchScoreInput) (*domain.UpadateMatchScoreInput, error) {
	return s.tournamentRepo.UpdateScore(ctx, tournamentOwnerID, req)
}

// CheckAndAdvanceRound checks if a round is complete and advances to the next round
func (s *service) CheckAndAdvanceRound(ctx context.Context, tournamentID int, round string) (bool, error) {
	return s.tournamentRepo.CheckAndAdvanceRound(ctx, tournamentID, round)
}

// GetGroupCount returns the number of groups for a tournament
func (s *service) GetGroupCount(ctx context.Context, tournamentID int) (int, error) {
	return s.tournamentRepo.GetGroupCount(ctx, tournamentID)
}

// GenerateKnockoutStage generates the knockout stage matches
func (s *service) GenerateKnockoutStage(ctx context.Context, tournamentID int) (bool, error) {
	return s.tournamentRepo.GenerateKnockoutStage(ctx, tournamentID)
}

// GenerateQuarterFinals generates quarter-final matches
func (s *service) GenerateQuarterFinals(ctx context.Context, tournamentID int) (bool, error) {
	return s.tournamentRepo.GenerateQuarterFinals(ctx, tournamentID)
}

// GenerateSemiFinals generates semi-final matches
func (s *service) GenerateSemiFinals(ctx context.Context, tournamentID int) (bool, error) {
	return s.tournamentRepo.GenerateSemiFinals(ctx, tournamentID)
}

// GenerateFinal generates the final match
func (s *service) GenerateFinal(ctx context.Context, tournamentID int) (bool, error) {
	return s.tournamentRepo.GenerateFinal(ctx, tournamentID)
}

// GetLeaderboard returns the leaderboard for a tournament
func (s *service) GetLeaderboard(ctx context.Context, tournamentID int) (map[int][]domain.PlayerStat, error) {
	return s.tournamentRepo.GetLeaderboard(ctx, tournamentID)
}

// GetTournamentType returns the tournament type
func (s *service) GetTournamentType(ctx context.Context, tournamentID int) (string, error) {
	return s.tournamentRepo.GetTournamentType(ctx, tournamentID)
}

// LeagueStyleSchedule creates a league style schedule
func (s *service) LeagueStyleSchedule(ctx context.Context, tournamentID int, approvedParticipants []*domain.Participant) error {
	return s.tournamentRepo.LeagueStyleSchedule(ctx, tournamentID, approvedParticipants)
}

// VerifyTournamentOwner verifies if a user is the owner of a tournament
func (s *service) VerifyTournamentOwner(ctx context.Context, tournamentID int, userID int) (bool, error) {
	return s.tournamentRepo.VerifyTournamentOwner(ctx, tournamentID, userID)
}
