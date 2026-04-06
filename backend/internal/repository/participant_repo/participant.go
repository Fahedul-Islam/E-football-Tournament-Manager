package participantrepo

import (
	"context"
	"database/sql"
	"time"
	"tournament-manager/internal/domain"
)

func (r *participantRepo) GetAllParticipant(ctx context.Context, tournament_id int) ([]*domain.Participant, error) {
	rows, err := r.db.QueryContext(ctx, "SELECT id, user_id, tournament_id, team_name, status, created_at FROM participants WHERE tournament_id=$1", tournament_id)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()

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


func (r *participantRepo) RequestToJoinTournament(ctx context.Context, req domain.ParticipantRequest) error {
	var total_current_participant int
	err := r.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM participants WHERE tournament_id=$1 AND status='approved'`, req.TournamentID).Scan(&total_current_participant)
	if err != nil {
		return err
	}
	var max_participants int
	err = r.db.QueryRowContext(ctx, `SELECT max_players FROM tournaments WHERE id=$1`, req.TournamentID).Scan(&max_participants)
	if err != nil {
		return err
	}
	if total_current_participant >= max_participants {
		return sql.ErrNoRows
	}
	query := `INSERT INTO participants (user_id, tournament_id, team_name, created_at) VALUES ($1, $2, $3, $4) RETURNING id`
	participant := domain.Participant{
		UserID:       req.UserID,
		TournamentID: req.TournamentID,
		TeamName:     req.TeamName,
		CreatedAt:    time.Now().Format(time.RFC3339),
	}
	return r.db.QueryRowContext(ctx, query, participant.UserID, participant.TournamentID, participant.TeamName, participant.CreatedAt).Scan(&participant.ID)
}

func (r *participantRepo) IsApprovedParticipant(ctx context.Context, tournament_id int, user_id int) (bool, error) {
	var exists bool
	err := r.db.QueryRowContext(ctx, "SELECT EXISTS(SELECT 1 FROM participants WHERE tournament_id = $1 AND user_id = $2 AND status = 'approved')", tournament_id, user_id).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}
