package participantrepo

import (
	"context"
	"database/sql"
	"errors"
	"tournament-manager/internal/domain"
)

func (r *participantRepo) ReactOnAnnouncement(ctx context.Context, tournamentID int, announcementID int, userID int, reaction string) (*domain.Announcement, error) {
	tx, err := r.db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	switch reaction {
	case "like":
		_, err = tx.ExecContext(ctx, "INSERT INTO announcement_reactions ( announcement_id, user_id, reaction_type) VALUES ($1, $2, $3) ON CONFLICT (announcement_id, user_id) DO UPDATE SET reaction_type='like'", announcementID, userID, reaction)
		if err != nil {
			return nil, err
		}
		_, err = tx.ExecContext(ctx, "UPDATE announcements SET likes_count = likes_count + 1 WHERE id=$1 AND tournament_id=$2", announcementID, tournamentID)
	case "dislike":
		_, err = tx.ExecContext(ctx, "INSERT INTO announcement_reactions (announcement_id, user_id, reaction_type) VALUES ($1, $2, $3) ON CONFLICT (announcement_id, user_id) DO UPDATE SET reaction_type='dislike'", announcementID, userID, reaction)
		if err != nil {
			return nil, err
		}
		_, err = tx.ExecContext(ctx, "UPDATE announcements SET dislikes_count = dislikes_count + 1 WHERE id=$1 AND tournament_id=$2", announcementID, tournamentID)
	default:
		return nil, errors.New("invalid reaction type")
	}

	if err != nil {
		return nil, err
	}

	var a domain.Announcement
	err = tx.QueryRowContext(ctx, "SELECT id, tournament_id, title, content, likes_count, dislikes_count, created_at FROM announcements WHERE id=$1 AND tournament_id=$2", announcementID, tournamentID).Scan(&a.ID, &a.TournamentID, &a.Title, &a.Content, &a.LikesCount, &a.DislikesCount, &a.CreatedAt)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return &a, nil
}

// get previous reaction
func (r *participantRepo) GetAnnouncementPrevReaction(ctx context.Context, tournamentID int, announcementID int, userID int) (string, error) {
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

func (r *participantRepo) RemoveAnnouncementReaction(ctx context.Context, tournamentID int, announcementID int, userID int, reaction string) (*domain.Announcement, error) {
	tx, err := r.db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	// Delete the reaction
	_, err = tx.ExecContext(ctx, "DELETE FROM announcement_reactions WHERE announcement_id=$1 AND user_id=$2", announcementID, userID)
	if err != nil {
		return nil, err
	}

	// decrease the likes or dislikes count in announcements table
	switch reaction {
	case "like":
		_, err = tx.ExecContext(ctx, "UPDATE announcements SET likes_count = GREATEST(likes_count - 1, 0) WHERE id=$1 AND tournament_id=$2", announcementID, tournamentID)
	case "dislike":
		_, err = tx.ExecContext(ctx, "UPDATE announcements SET dislikes_count = GREATEST(dislikes_count - 1, 0) WHERE id=$1 AND tournament_id=$2", announcementID, tournamentID)
	default:
		return nil, errors.New("invalid reaction type")
	}

	if err != nil {
		return nil, err
	}

	var a domain.Announcement
	err = tx.QueryRowContext(ctx, "SELECT id, tournament_id, title, content, likes_count, dislikes_count, created_at FROM announcements WHERE id=$1 AND tournament_id=$2", announcementID, tournamentID).Scan(&a.ID, &a.TournamentID, &a.Title, &a.Content, &a.LikesCount, &a.DislikesCount, &a.CreatedAt)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return &a, nil
}
