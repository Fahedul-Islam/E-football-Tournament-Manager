package tournamentmanagerrepo

import "tournament-manager/domain"

func (r *tournamentManagerRepo) GetApprovedParticipants(tournament_id int) ([]*domain.Participant, error) {
	query := `SELECT * FROM participants WHERE tournament_id=$1 AND status='approved'`
	rows, err := r.db.Query(query, tournament_id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var participants []*domain.Participant
	for rows.Next() {
		p := &domain.Participant{}
		if err := rows.Scan(&p.ID, &p.UserID, &p.TournamentID, &p.TeamName, &p.Status, &p.CreatedAt); err != nil {
			return nil, err
		}
		participants = append(participants, p)
	}
	return participants, nil
}
