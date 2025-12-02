package tournamentmanagerrepo

import "tournament-manager/domain"

func (r *tournamentManagerRepo) GetLeaderboard(tournament_id int) (map[int][]domain.PlayerStat, error) {
	// get group id from group table by tournament id
	query := `SELECT id FROM groups WHERE tournament_id = $1`
	rows, err := r.db.Query(query, tournament_id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var groupIDs []int
	for rows.Next() {
		var groupID int
		if err := rows.Scan(&groupID); err != nil {
			return nil, err
		}
		groupIDs = append(groupIDs, groupID)
	}

	leaderboardMap := make(map[int][]domain.PlayerStat)

	for _, groupID := range groupIDs {
		query := `
			SELECT participant_id, matches_played, wins, draws, losses,
			       goals_scored, goals_conceded, goal_difference, points
			FROM player_stats
			WHERE  group_id = $1
			ORDER BY points DESC, goal_difference DESC, goals_scored DESC
		`
		rows, err := r.db.Query(query, groupID)
		if err != nil {
			return nil, err
		}
		defer rows.Close()

		var groupLeaderboard []domain.PlayerStat

		for rows.Next() {
			var stat domain.PlayerStat

			err := rows.Scan(
				&stat.ParticipantID,
				&stat.MatchesPlayed,
				&stat.Wins,
				&stat.Draws,
				&stat.Losses,
				&stat.GoalsScored,
				&stat.GoalsConceded,
				&stat.GoalDifference,
				&stat.Points,
			)
			if err != nil {
				return nil, err
			}
			groupLeaderboard = append(groupLeaderboard, stat)
		}
		leaderboardMap[groupID] = groupLeaderboard
	}

	return leaderboardMap, nil	

}

