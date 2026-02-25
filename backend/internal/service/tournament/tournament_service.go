package tournament

import (
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
func (s *service) CreateTournament(createdBy int, request domain.TournamentCreateRequest) error {
	return s.tournamentRepo.CreateTournament(createdBy, request)
}

// GetTournamentByID returns a tournament by ID
func (s *service) GetTournamentByID(id int) (*domain.Tournament, error) {
	return s.tournamentRepo.GetTournamentByID(id)
}

// GetAllTournaments returns all tournaments for a given owner
func (s *service) GetAllTournaments(tournamentOwnerID int) ([]*domain.Tournament, error) {
	return s.tournamentRepo.GetAllTournaments(tournamentOwnerID)
}

// UpdateTournament updates a tournament
func (s *service) UpdateTournament(tournamentOwnerID int, tournamentID int, tournament domain.TournamentCreateRequest) error {
	return s.tournamentRepo.UpdateTournament(tournamentOwnerID, tournamentID, tournament)
}

// DeleteTournament deletes a tournament
func (s *service) DeleteTournament(tournamentOwnerID int, tournamentID int) error {
	return s.tournamentRepo.DeleteTournament(tournamentOwnerID, tournamentID)
}

// ApproveParticipant approves a participant for a tournament
func (s *service) ApproveParticipant(tournamentOwnerID int, req domain.ParticipantRequest) error {
	return s.tournamentRepo.ApproveParticipant(tournamentOwnerID, req)
}

// RejectParticipant rejects a participant from a tournament
func (s *service) RejectParticipant(tournamentOwnerID int, req domain.ParticipantRequest) error {
	return s.tournamentRepo.RejectParticipant(tournamentOwnerID, req)
}

// AddParticipant adds a participant to a tournament
func (s *service) AddParticipant(tournamentOwnerID int, participant domain.ParticipantRequest) error {
	return s.tournamentRepo.AddParticipant(tournamentOwnerID, participant)
}

// RemoveParticipant removes a participant from a tournament
func (s *service) RemoveParticipant(tournamentOwnerID int, req domain.ParticipantRequest) error {
	return s.tournamentRepo.RemoveParticipant(tournamentOwnerID, req)
}

// GetAllParticipant returns all participants for a tournament
func (s *service) GetAllParticipant(tournamentID int) ([]*domain.Participant, error) {
	return s.tournamentRepo.GetAllParticipant(tournamentID)
}

// GetApprovedParticipants returns all approved participants for a tournament
func (s *service) GetApprovedParticipants(tournamentID int) ([]*domain.Participant, error) {
	return s.tournamentRepo.GetApprovedParticipants(tournamentID)
}

// CreateMatchSchedules creates match schedules for a tournament
func (s *service) CreateMatchSchedules(tournamentID int, groupCount int, approvedParticipants []*domain.Participant) error {
	return s.tournamentRepo.CreateMatchSchedules(tournamentID, groupCount, approvedParticipants)
}

// GenerateGroups generates groups for a tournament
func (s *service) GenerateGroups(tournamentID int, groupCount int, approvedParticipants []*domain.Participant) error {
	return s.tournamentRepo.GenerateGroups(tournamentID, groupCount, approvedParticipants)
}

// GetAllMatches returns all matches for a tournament
func (s *service) GetAllMatches(tournamentID int) ([]*domain.Match, error) {
	return s.tournamentRepo.GetAllMatches(tournamentID)
}

// UpdateScore updates the score for a match
func (s *service) UpdateScore(tournamentOwnerID int, req *domain.UpadateMatchScoreInput) (*domain.UpadateMatchScoreInput, error) {
	return s.tournamentRepo.UpdateScore(tournamentOwnerID, req)
}

// CheckAndAdvanceRound checks if a round is complete and advances to the next round
func (s *service) CheckAndAdvanceRound(tournamentID int, round string) (bool, error) {
	return s.tournamentRepo.CheckAndAdvanceRound(tournamentID, round)
}

// GetGroupCount returns the number of groups for a tournament
func (s *service) GetGroupCount(tournamentID int) (int, error) {
	return s.tournamentRepo.GetGroupCount(tournamentID)
}

// GenerateKnockoutStage generates the knockout stage matches
func (s *service) GenerateKnockoutStage(tournamentID int) (bool, error) {
	return s.tournamentRepo.GenerateKnockoutStage(tournamentID)
}

// GenerateQuarterFinals generates quarter-final matches
func (s *service) GenerateQuarterFinals(tournamentID int) (bool, error) {
	return s.tournamentRepo.GenerateQuarterFinals(tournamentID)
}

// GenerateSemiFinals generates semi-final matches
func (s *service) GenerateSemiFinals(tournamentID int) (bool, error) {
	return s.tournamentRepo.GenerateSemiFinals(tournamentID)
}

// GenerateFinal generates the final match
func (s *service) GenerateFinal(tournamentID int) (bool, error) {
	return s.tournamentRepo.GenerateFinal(tournamentID)
}

// GetLeaderboard returns the leaderboard for a tournament
func (s *service) GetLeaderboard(tournamentID int) (map[int][]domain.PlayerStat, error) {
	return s.tournamentRepo.GetLeaderboard(tournamentID)
}

// GetTournamentType returns the tournament type
func (s *service) GetTournamentType(tournamentID int) (string, error) {
	return s.tournamentRepo.GetTournamentType(tournamentID)
}

// LeagueStyleSchedule creates a league style schedule
func (s *service) LeagueStyleSchedule(tournamentID int, approvedParticipants []*domain.Participant) error {
	return s.tournamentRepo.LeagueStyleSchedule(tournamentID, approvedParticipants)
}

// VerifyTournamentOwner verifies if a user is the owner of a tournament
func (s *service) VerifyTournamentOwner(tournamentID int, userID int) (bool, error) {
	return s.tournamentRepo.VerifyTournamentOwner(tournamentID, userID)
}
