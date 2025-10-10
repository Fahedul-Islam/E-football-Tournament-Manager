package tournamentmanagerrepo

import "tournament-manager/domain"

func (r *tournamentManagerRepo) GetAllMatches(tournament_id int) ([]*domain.Match, error) {
	rows, err := r.db.Query("SELECT * FROM matches WHERE tournament_id=$1", tournament_id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var matches []*domain.Match
	for rows.Next() {
		m := &domain.Match{}
		if err := rows.Scan(&m.ID, &m.TournamentID, &m.GroupID, &m.Round, &m.ParticipantAID, &m.ParticipantBID, &m.ScoreA, &m.ScoreB, &m.MatchDate, &m.CreatedAt, &m.Status); err != nil {
			return nil, err
		}
		matches = append(matches, m)
	}
	return matches, nil
}
