package tournament

import (
	"context"
	"tournament-manager/internal/domain"
)

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
