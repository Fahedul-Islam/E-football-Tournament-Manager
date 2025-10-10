package participantrepo

import (
	"database/sql"
	"time"
	"tournament-manager/domain"
	"tournament-manager/rest/handler/participant"
)

type ParticipantRepo interface {
	participant.Service
}

type participantRepo struct {
	db *sql.DB
}

func NewParticipantRepo(db *sql.DB) ParticipantRepo {
	return &participantRepo{db: db}
}

func (r *participantRepo) RequestToJoinTournament(req domain.ParticipantRequest) error {
	var total_current_participant int
	err := r.db.QueryRow(`SELECT COUNT(*) FROM participants WHERE tournament_id=$1 AND status='approved'`, req.TournamentID).Scan(&total_current_participant)
	if err != nil {
		return err
	}
	var max_participants int
	err = r.db.QueryRow(`SELECT max_players FROM tournaments WHERE id=$1`, req.TournamentID).Scan(&max_participants)
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
	return r.db.QueryRow(query, participant.UserID, participant.TournamentID, participant.TeamName, participant.CreatedAt).Scan(&participant.ID)
}

func (r *participantRepo) IsApprovedParticipant(tournament_id int, user_id int) (bool, error) {
	var exists bool
	err := r.db.QueryRow("SELECT EXISTS(SELECT 1 FROM participants WHERE tournament_id = $1 AND user_id = $2 AND status = 'approved')", tournament_id, user_id).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

func (r *participantRepo) GetGroupDistribution(tournament_id int) ([]*domain.Group, error) {
	rows, err := r.db.Query(`
		SELECT g.id, g.name, p.id, p.user_id, p.team_name, p.status, p.created_at
		FROM groups g
		LEFT JOIN group_teams gt ON g.id = gt.group_id
		LEFT JOIN participants p ON gt.participant_id = p.id
		WHERE g.tournament_id = $1
		ORDER BY g.name, p.id
	`, tournament_id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	groupMap := make(map[int]*domain.Group)
	for rows.Next() {
		var groupID int
		var groupName string
		var participantID sql.NullInt64
		var userID sql.NullInt64
		var teamName sql.NullString
		var status sql.NullString
		var createdAt sql.NullString

		if err := rows.Scan(&groupID, &groupName, &participantID, &userID, &teamName, &status, &createdAt); err != nil {
			return nil, err
		}

		if _, exists := groupMap[groupID]; !exists {
			groupMap[groupID] = &domain.Group{
				GroupID:      groupID,
				GroupName:    groupName,
				Participants: []*domain.Participant{},
			}
		}

		if participantID.Valid {
			p := &domain.Participant{
				ID:           int(participantID.Int64),
				UserID:       int(userID.Int64),
				TournamentID: tournament_id,
				TeamName:     teamName.String,
				Status:       status.String,
				CreatedAt:    createdAt.String,
			}
			groupMap[groupID].Participants = append(groupMap[groupID].Participants, p)
		}
	}

	var groups []*domain.Group
	for _, group := range groupMap {
		groups = append(groups, group)
	}

	return groups, nil
}

func (r *participantRepo) SeeMatchSchedule(tournament_id int) ([]*domain.Match, error) {
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
