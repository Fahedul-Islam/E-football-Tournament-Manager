package announcement

import (
	"context"
	"database/sql"
	"errors"
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

func (r *announcementRepo) VerifyCommentOwner(ctx context.Context, commentID int, userID int) (bool, error) {
	var ownerID int
	err := r.db.QueryRowContext(ctx, "SELECT user_id FROM announcement_comments WHERE id = $1", commentID).Scan(&ownerID)
	if err != nil {
		return false, err
	}
	return ownerID == userID, nil
}

func (r *announcementRepo) GetAnnouncementPrevReaction(ctx context.Context, tournamentID int, announcementID int, userID int) (string, error) {
	var reaction string
	err := r.db.QueryRowContext(ctx, "SELECT reaction_type FROM announcement_reactions WHERE announcement_id=$1 AND user_id=$2", announcementID, userID).Scan(&reaction)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", nil // No previous reaction
		}
		return "", err
	}
	return reaction, nil
}

func (r *announcementRepo) GetCommentPrevReaction(ctx context.Context, tournamentID int, commentID int, userID int) (string, error) {
	var reaction string
	err := r.db.QueryRowContext(ctx, "SELECT reaction_type FROM announcement_comment_reactions WHERE comment_id=$1 AND user_id=$2", commentID, userID).Scan(&reaction)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", nil // No previous reaction
		}
		return "", err
	}
	return reaction, nil
}

func (r *announcementRepo) VerifyCommentBelongsToTournament(ctx context.Context, tournamentID int, commentID int) (bool, error) {
	var exists bool
	err := r.db.QueryRowContext(ctx, `
		SELECT EXISTS(
			SELECT 1 FROM announcement_comments ac
			INNER JOIN announcements a ON a.id = ac.announcement_id
			WHERE ac.id = $1 AND a.tournament_id = $2
		)
	`, commentID, tournamentID).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}
