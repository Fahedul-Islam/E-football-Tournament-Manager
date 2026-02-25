package repository

import "tournament-manager/internal/domain"

// ParticipantRepository defines the interface for participant data access operations
type ParticipantRepository interface {
	RequestToJoinTournament(req domain.ParticipantRequest) error
	IsApprovedParticipant(tournamentID int, userID int) (bool, error)
	GetGroupDistribution(tournamentID int) ([]*domain.Group, error)
	SeeMatchSchedule(tournamentID int) ([]*domain.Match, error)
}
