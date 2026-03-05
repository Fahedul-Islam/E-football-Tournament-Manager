package announcement

import (
	"context"
	"tournament-manager/internal/domain"
)

func (r *announcementRepo) VerifyTournamentOwner(ctx context.Context, tournamentID int, userID int) (bool, error) {
	var ownerID int
	err := r.db.QueryRowContext(ctx, "SELECT created_by FROM tournaments WHERE id = $1", tournamentID).Scan(&ownerID)
	if err != nil {
		return false, err
	}
	return ownerID == userID, nil
}

func (r *announcementRepo) GetAllParticipant(ctx context.Context, tournamentID int) ([]*domain.Participant, error) {
	rows, err := r.db.QueryContext(ctx, "SELECT id, user_id, tournament_id, team_name, status, created_at FROM participants WHERE tournament_id=$1", tournamentID)
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
