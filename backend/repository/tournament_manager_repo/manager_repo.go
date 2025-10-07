package tournamentmanagerrepo

import (
	"database/sql"
	"time"
	"tournament-manager/domain"
	tournamentmanager "tournament-manager/rest/handler/tournamentManager"
)

type TournamentManagerRepo interface {
	tournamentmanager.Service
}

type tournamentManagerRepo struct {
	db *sql.DB
}

func NewTournamentManagerRepo(db *sql.DB) TournamentManagerRepo {
	return &tournamentManagerRepo{db: db}
}

func (r *tournamentManagerRepo) CreateTournament(created_by int, request domain.TournamentCreateRequest) error {
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
	return r.db.QueryRow(query, tournament.Name, tournament.Description, tournament.TournamentType, tournament.MaxPlayers,  tournament.CreatedBy, tournament.StartDate, tournament.EndDate).Scan(&tournament.ID)
}

func (r *tournamentManagerRepo) GetTournamentByID(id int) (*domain.Tournament, error) {
	var tournament domain.Tournament
	query := `SELECT * FROM tournaments WHERE id = $1`
	if err := r.db.QueryRow(query, id).Scan(&tournament.ID, &tournament.Name, &tournament.Description, &tournament.StartDate, &tournament.EndDate, &tournament.CreatedBy); err != nil {
		return nil, err
	}
	return &tournament, nil
}

func (r *tournamentManagerRepo) GetAllTournaments() ([]*domain.Tournament, error) {
	rows, err := r.db.Query("SELECT id, name, description, tournament_type, max_players, start_date, end_date, created_by FROM tournaments")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tournaments []*domain.Tournament
	for rows.Next() {
		t := &domain.Tournament{}
		if err := rows.Scan(&t.ID, &t.Name, &t.Description, &t.TournamentType, &t.MaxPlayers, &t.StartDate, &t.EndDate, &t.CreatedBy); err != nil {
			return nil, err
		}
		tournaments = append(tournaments, t)
	}
	return tournaments, nil
}

func (r *tournamentManagerRepo) UpdateTournament(id int, tournament domain.TournamentCreateRequest) error {
	_, err := r.db.Exec("UPDATE tournaments SET name = $1, description = $2, start_date = $3, end_date = $4 WHERE id = $5",
		tournament.Name, tournament.Description, tournament.StartDate, tournament.EndDate, id)
	return err
}

func (r *tournamentManagerRepo) DeleteTournament(id int) error {
	_, err := r.db.Exec("DELETE FROM tournaments WHERE id = $1", id)
	return err
}

func (r *tournamentManagerRepo) ApproveParticipant(req domain.ParticipantRequest) error {
	// check if maximum participants reached
	var count int
	err := r.db.QueryRow("SELECT COUNT(*) FROM participants WHERE tournament_id = $1 AND status = 'approved'", req.TournamentID).Scan(&count)
	if err != nil {
		return err
	}

	var maxParticipants int
	err = r.db.QueryRow("SELECT max_participants FROM tournaments WHERE id = $1", req.TournamentID).Scan(&maxParticipants)
	if err != nil {
		return err
	}

	if count >= maxParticipants {
		return sql.ErrNoRows
	}
	_, err = r.db.Exec("UPDATE participants SET status = 'approved' WHERE user_id = $1 AND tournament_id = $2", req.UserID, req.TournamentID)
	return err
}

func (r *tournamentManagerRepo) RejectParticipant(req domain.ParticipantRequest) error {
	_, err := r.db.Exec("UPDATE participants SET status = 'rejected' WHERE user_id = $1 AND tournament_id = $2", req.UserID, req.TournamentID)
	return err
}

func (r *tournamentManagerRepo) AddParticipant(participant domain.ParticipantRequest) error {
	now := time.Now().Format(time.RFC3339)
	addedParticipant := domain.Participant{
		UserID:       participant.UserID,
		TournamentID: participant.TournamentID,
		TeamName:     participant.TeamName,
		Status:       "approved",
		CreatedAt:    now,
	}
	query := `INSERT INTO participants (user_id, tournament_id, team_name, status, created_at) VALUES ($1, $2, $3, $4, $5) RETURNING id`
	return r.db.QueryRow(query, addedParticipant.UserID, addedParticipant.TournamentID, addedParticipant.TeamName, addedParticipant.Status, addedParticipant.CreatedAt).Scan(&addedParticipant.ID)
}

func (r *tournamentManagerRepo) RemoveParticipant(req domain.ParticipantRequest) error {
	_, err := r.db.Exec("DELETE FROM participants WHERE user_id = $1 AND tournament_id = $2", req.UserID, req.TournamentID)
	return err
}
