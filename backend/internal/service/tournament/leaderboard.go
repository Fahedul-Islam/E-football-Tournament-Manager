package tournament

import (
	"context"
	"tournament-manager/internal/domain"
)

// GetLeaderboard returns the leaderboard for a tournament
func (s *service) GetLeaderboard(ctx context.Context, tournamentID int) (map[int][]domain.PlayerStat, error) {
	return s.tournamentRepo.GetLeaderboard(ctx, tournamentID)
}
