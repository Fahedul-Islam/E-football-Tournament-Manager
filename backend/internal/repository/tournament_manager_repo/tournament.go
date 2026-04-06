package tournamentmanagerrepo

import (
	"context"
	"database/sql"
	"tournament-manager/internal/domain"
)

func (r *tournamentManagerRepo) CreateTournament(ctx context.Context, created_by int, request domain.TournamentCreateRequest) error {
	query := `INSERT INTO tournaments (name, description, tournament_type, max_players,  created_by, start_date, end_date) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id`
	tournament := domain.Tournament{
		Name:           request.Name,
		Description:    request.Description,
		TournamentType: request.TournamentType,
		MaxPlayers:     request.MaxPlayers,
		StartDate:      request.StartDate,
		EndDate:        request.EndDate,
		CreatedBy:      created_by,
	}
	return r.db.QueryRowContext(ctx, query, tournament.Name, tournament.Description, tournament.TournamentType, tournament.MaxPlayers, tournament.CreatedBy, tournament.StartDate, tournament.EndDate).Scan(&tournament.ID)
}

func (r *tournamentManagerRepo) DeleteTournament(ctx context.Context, tournament_owner_id int, id int) error {
	var exists bool
	err := r.db.QueryRowContext(ctx, "SELECT EXISTS(SELECT 1 FROM tournaments WHERE id = $1 AND created_by = $2)", id, tournament_owner_id).Scan(&exists)
	if err != nil {
		return err
	}
	if !exists {
		return sql.ErrNoRows
	}
	_, err = r.db.ExecContext(ctx, "DELETE FROM tournaments WHERE id = $1", id)
	return err
}

func (r *tournamentManagerRepo) GetAllTournaments(ctx context.Context, tournament_owner_id int) ([]*domain.Tournament, error) {
	rows, err := r.db.QueryContext(ctx, "SELECT id, name, description, tournament_type, max_players, created_by, start_date, end_date, created_at FROM tournaments WHERE created_by = $1", tournament_owner_id)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	var tournaments []*domain.Tournament
	for rows.Next() {
		t := &domain.Tournament{}
		if err := rows.Scan(&t.ID, &t.Name, &t.Description, &t.TournamentType, &t.MaxPlayers, &t.CreatedBy, &t.StartDate, &t.EndDate, &t.CreatedAt); err != nil {
			return nil, err
		}
		tournaments = append(tournaments, t)
	}
	return tournaments, nil
}

func (r *tournamentManagerRepo) GetTournamentByID(ctx context.Context, id int) (*domain.Tournament, error) {
	var tournament domain.Tournament
	query := `SELECT id, name, description, tournament_type, max_players, created_by, start_date, end_date, created_at FROM tournaments WHERE id = $1`
	if err := r.db.QueryRowContext(ctx, query, id).Scan(
		&tournament.ID, &tournament.Name, &tournament.Description,
		&tournament.TournamentType, &tournament.MaxPlayers, &tournament.CreatedBy,
		&tournament.StartDate, &tournament.EndDate, &tournament.CreatedAt,
	); err != nil {
		return nil, err
	}
	return &tournament, nil
}

func (r *tournamentManagerRepo) UpdateTournament(ctx context.Context, tournament_owner_id int, tournament_id int, tournament domain.TournamentCreateRequest) error {
	_, err := r.db.ExecContext(ctx, "UPDATE tournaments SET name = $1, description = $2, start_date = $3, end_date = $4 WHERE id = $5 AND created_by = $6",
		tournament.Name, tournament.Description, tournament.StartDate, tournament.EndDate, tournament_id, tournament_owner_id)
	return err
}

func (r *tournamentManagerRepo) GetTournamentType(ctx context.Context, tournament_id int) (string, error) {
	var tournament_type string
	err := r.db.QueryRowContext(ctx, "SELECT tournament_type FROM tournaments WHERE id = $1", tournament_id).Scan(&tournament_type)
	if err != nil {
		return "", err
	}
	return tournament_type, nil
}
