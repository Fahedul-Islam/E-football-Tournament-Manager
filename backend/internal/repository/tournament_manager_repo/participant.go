package tournamentmanagerrepo

import (
	"context"
	"database/sql"
	"time"
	"tournament-manager/internal/domain"
)

func (r *tournamentManagerRepo) GetAllParticipant(ctx context.Context, tournament_id int) ([]*domain.Participant, error) {
	query := `SELECT * FROM participants WHERE tournament_id=$1`
	rows, err := r.db.QueryContext(ctx, query, tournament_id)
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

func (r *tournamentManagerRepo) AddParticipant(ctx context.Context, tournament_owner_id int, participant domain.ParticipantRequest) error {
	// Check if the tournament exists and is created by the tournament_owner_id
	var exists bool
	err := r.db.QueryRowContext(ctx, "SELECT EXISTS(SELECT 1 FROM tournaments WHERE id = $1 AND created_by = $2)", participant.TournamentID, tournament_owner_id).Scan(&exists)
	if err != nil {
		return err
	}
	if !exists {
		return sql.ErrNoRows
	}
	// check if maximum participants reached
	var count int
	err = r.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM participants WHERE tournament_id = $1 AND status = 'approved'", participant.TournamentID).Scan(&count)
	if err != nil {
		return err
	}

	var maxParticipants int
	err = r.db.QueryRowContext(ctx, "SELECT max_players FROM tournaments WHERE id = $1", participant.TournamentID).Scan(&maxParticipants)
	if err != nil {
		return err
	}

	if count >= maxParticipants {
		return sql.ErrNoRows
	}
	// Add participant with status 'approved'
	now := time.Now().Format(time.RFC3339)
	addedParticipant := domain.Participant{
		UserID:       participant.UserID,
		TournamentID: participant.TournamentID,
		TeamName:     participant.TeamName,
		Status:       "approved",
		CreatedAt:    now,
	}
	query := `INSERT INTO participants (user_id, tournament_id, team_name, status, created_at) VALUES ($1, $2, $3, $4, $5) RETURNING id`
	return r.db.QueryRowContext(ctx, query, addedParticipant.UserID, addedParticipant.TournamentID, addedParticipant.TeamName, addedParticipant.Status, addedParticipant.CreatedAt).Scan(&addedParticipant.ID)
}

func (r *tournamentManagerRepo) ApproveParticipant(ctx context.Context, tournament_owner_id int, req domain.ParticipantRequest) error {
	// Check if the tournament exists and is created by the tournament_owner_id
	var exists bool
	err := r.db.QueryRowContext(ctx, "SELECT EXISTS(SELECT 1 FROM tournaments WHERE id = $1 AND created_by = $2)", req.TournamentID, tournament_owner_id).Scan(&exists)
	if err != nil {
		return err
	}
	if !exists {
		return sql.ErrNoRows
	}
	// check if maximum participants reached
	var count int
	err = r.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM participants WHERE tournament_id = $1 AND status = 'approved'", req.TournamentID).Scan(&count)
	if err != nil {
		return err
	}

	var maxParticipants int
	err = r.db.QueryRowContext(ctx, "SELECT max_players FROM tournaments WHERE id = $1", req.TournamentID).Scan(&maxParticipants)
	if err != nil {
		return err
	}

	if count >= maxParticipants {
		return sql.ErrNoRows
	}
	_, err = r.db.ExecContext(ctx, "UPDATE participants SET status = 'approved' WHERE user_id = $1 AND tournament_id = $2", req.UserID, req.TournamentID)
	return err
}

func (r *tournamentManagerRepo) GetApprovedParticipants(ctx context.Context, tournament_id int) ([]*domain.Participant, error) {
	query := `SELECT * FROM participants WHERE tournament_id=$1 AND status='approved'`
	rows, err := r.db.QueryContext(ctx, query, tournament_id)
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

func (r *tournamentManagerRepo) RejectParticipant(ctx context.Context, tournament_owner_id int, req domain.ParticipantRequest) error {
	// Check if the tournament exists and is created by the tournament_owner_id
	var exists bool
	err := r.db.QueryRowContext(ctx, "SELECT EXISTS(SELECT 1 FROM tournaments WHERE id = $1 AND created_by = $2)", req.TournamentID, tournament_owner_id).Scan(&exists)
	if err != nil {
		return err
	}
	if !exists {
		return sql.ErrNoRows
	}
	_, err = r.db.ExecContext(ctx, "UPDATE participants SET status = 'rejected' WHERE user_id = $1 AND tournament_id = $2", req.UserID, req.TournamentID)
	return err
}

func (r *tournamentManagerRepo) RemoveParticipant(ctx context.Context, tournament_owner_id int, req domain.ParticipantRequest) error {
	// Check if the tournament exists and is created by the tournament_owner_id
	var exists bool
	err := r.db.QueryRowContext(ctx, "SELECT EXISTS(SELECT 1 FROM tournaments WHERE id = $1 AND created_by = $2)", req.TournamentID, tournament_owner_id).Scan(&exists)
	if err != nil {
		return err
	}
	if !exists {
		return sql.ErrNoRows
	}
	_, err = r.db.ExecContext(ctx, "DELETE FROM participants WHERE user_id = $1 AND tournament_id = $2", req.UserID, req.TournamentID)
	return err
}
