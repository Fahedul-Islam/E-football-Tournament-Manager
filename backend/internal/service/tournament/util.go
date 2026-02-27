package tournament

import "context"

// CheckAndAdvanceRound checks if a round is complete and advances to the next round
func (s *service) CheckAndAdvanceRound(ctx context.Context, tournamentID int, round string) (bool, error) {
	return s.tournamentRepo.CheckAndAdvanceRound(ctx, tournamentID, round)
}

// GetGroupCount returns the number of groups for a tournament
func (s *service) GetGroupCount(ctx context.Context, tournamentID int) (int, error) {
	return s.tournamentRepo.GetGroupCount(ctx, tournamentID)
}

// GetTournamentType returns the tournament type
func (s *service) GetTournamentType(ctx context.Context, tournamentID int) (string, error) {
	return s.tournamentRepo.GetTournamentType(ctx, tournamentID)
}

// VerifyTournamentOwner verifies if a user is the owner of a tournament
func (s *service) VerifyTournamentOwner(ctx context.Context, tournamentID int, userID int) (bool, error) {
	return s.tournamentRepo.VerifyTournamentOwner(ctx, tournamentID, userID)
}
