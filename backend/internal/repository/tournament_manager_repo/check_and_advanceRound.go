package tournamentmanagerrepo

import "context"

func (r *tournamentManagerRepo) CheckAndAdvanceRound(ctx context.Context, tournament_id int, round string) (bool, error) {

	var unplayedCount int
	err := r.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM matches WHERE tournament_id = $1 AND round = $2 AND status != 'completed'", tournament_id, round).Scan(&unplayedCount)
	if err != nil {
		return false, err
	}
	if unplayedCount > 0 {
		return false, nil
	}

	return true, nil
}
