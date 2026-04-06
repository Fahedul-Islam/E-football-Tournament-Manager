package announcement

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"tournament-manager/internal/domain"
)

func (r *announcementRepo) ReactOnAnnouncement(ctx context.Context, tournamentID int, announcementID int, userID int, reaction string) (*domain.Announcement, error) {
	tx, err := r.db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		return nil, err
	}
	defer func() { _ = tx.Rollback() }()

	switch reaction {
	case "like":
		_, err = tx.ExecContext(ctx, "INSERT INTO announcement_reactions (announcement_id, user_id, reaction_type) VALUES ($1, $2, $3) ON CONFLICT (announcement_id, user_id) DO UPDATE SET reaction_type='like'", announcementID, userID, reaction)
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

func (r *announcementRepo) RemoveAnnouncementReaction(ctx context.Context, tournamentID int, announcementID int, userID int, reaction string) (*domain.Announcement, error) {
	tx, err := r.db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		return nil, err
	}
	defer func() { _ = tx.Rollback() }()

	// Delete the reaction
	_, err = tx.ExecContext(ctx, "DELETE FROM announcement_reactions WHERE announcement_id=$1 AND user_id=$2", announcementID, userID)
	if err != nil {
		return nil, err
	}

	// Decrease the likes or dislikes count in announcements table
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

func (r *announcementRepo) ReactToComment(ctx context.Context, commentID int, userID int, reactionType string) (*domain.AnnouncementComment, error) {
	log.Printf("Reacting to comment: commentID=%d, userID=%d, reactionType=%s", commentID, userID, reactionType)
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer func() { _ = tx.Rollback() }()
	_, err = tx.ExecContext(ctx, "INSERT INTO announcement_comment_reactions (comment_id, user_id, reaction_type) VALUES ($1, $2, $3) ON CONFLICT (comment_id, user_id) DO NOTHING", commentID, userID, reactionType)
	if err != nil {
		return nil, err
	}
	if reactionType == "like" {
		_, err = tx.ExecContext(ctx, "UPDATE announcement_comments SET likes_count = likes_count + 1 WHERE id=$1", commentID)
		if err != nil {
			return nil, err
		}
	}
	if reactionType == "dislike" {
		_, err = tx.ExecContext(ctx, "UPDATE announcement_comments SET dislikes_count = dislikes_count + 1 WHERE id=$1", commentID)
		if err != nil {
			return nil, err
		}
	}
	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	var c domain.AnnouncementComment
	err = r.db.QueryRowContext(ctx, "SELECT id, announcement_id, user_id, parent_comment_id, content, likes_count, dislikes_count, is_edited, created_at FROM announcement_comments WHERE id=$1", commentID).Scan(&c.ID, &c.AnnouncementID, &c.UserID, &c.ParentCommentID, &c.Content, &c.LikesCount, &c.DislikesCount, &c.IsEdited, &c.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func (r *announcementRepo) RemoveReactionFromComment(ctx context.Context, commentID int, userID int) (*domain.AnnouncementComment, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer func() { _ = tx.Rollback() }()
	var reactionType string
	err = tx.QueryRowContext(ctx, "SELECT reaction_type FROM announcement_comment_reactions WHERE comment_id=$1 AND user_id=$2", commentID, userID).Scan(&reactionType)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// No reaction to remove, just return the comment
			var c domain.AnnouncementComment
			err = r.db.QueryRowContext(ctx, "SELECT id, announcement_id, user_id, parent_comment_id, content, likes_count, dislikes_count, is_edited, created_at FROM announcement_comments WHERE id=$1", commentID).Scan(&c.ID, &c.AnnouncementID, &c.UserID, &c.ParentCommentID, &c.Content, &c.LikesCount, &c.DislikesCount, &c.IsEdited, &c.CreatedAt)
			if err != nil {
				return nil, err
			}
			return &c, nil
		}
		return nil, err
	}
	_, err = tx.ExecContext(ctx, "DELETE FROM announcement_comment_reactions WHERE comment_id=$1 AND user_id=$2", commentID, userID)
	if err != nil {
		return nil, err
	}
	if reactionType == "like" {
		_, err = tx.ExecContext(ctx, "UPDATE announcement_comments SET likes_count = GREATEST(likes_count - 1, 0) WHERE id=$1", commentID)
		if err != nil {
			return nil, err
		}
	}
	if reactionType == "dislike" {
		_, err = tx.ExecContext(ctx, "UPDATE announcement_comments SET dislikes_count = GREATEST(dislikes_count - 1, 0) WHERE id=$1", commentID)
		if err != nil {
			return nil, err
		}
	}
	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	var c domain.AnnouncementComment
	err = r.db.QueryRowContext(ctx, "SELECT id, announcement_id, user_id, parent_comment_id, content, likes_count, dislikes_count, is_edited, created_at FROM announcement_comments WHERE id=$1", commentID).Scan(&c.ID, &c.AnnouncementID, &c.UserID, &c.ParentCommentID, &c.Content, &c.LikesCount, &c.DislikesCount, &c.IsEdited, &c.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &c, nil
}
