package announcement

import (
	"context"
	"tournament-manager/internal/domain"
)

func (r *announcementRepo) AddComment(ctx context.Context, announcementID int, userID int, parentCommentID *int, content *string) (*domain.AnnouncementComment, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()
	var commentID int
	err = tx.QueryRowContext(ctx, "INSERT INTO announcement_comments (announcement_id, user_id, parent_comment_id, content) VALUES ($1, $2, $3, $4) RETURNING id", announcementID, userID, parentCommentID, content).Scan(&commentID)
	if err != nil {
		return nil, err
	}

	// increase the comments count in announcements table
	_, err = tx.ExecContext(ctx, "UPDATE announcements SET comments_count = comments_count + 1 WHERE id=$1", announcementID)
	if err != nil {
		return nil, err
	}

	var c domain.AnnouncementComment
	err = tx.QueryRowContext(ctx, "SELECT id, announcement_id, user_id, parent_comment_id, content, created_at FROM announcement_comments WHERE id=$1", commentID).Scan(&c.ID, &c.AnnouncementID, &c.UserID, &c.ParentCommentID, &c.Content, &c.CreatedAt)
	if err != nil {
		return nil, err
	}
	err = tx.Commit()
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func (r *announcementRepo) GetComments(ctx context.Context, announcementID int, parentCommentID *int) ([]*domain.AnnouncementComment, error) {
	var query string
	var args []interface{}
	if parentCommentID != nil {
		query = "SELECT id, announcement_id, user_id, parent_comment_id, content, created_at FROM announcement_comments WHERE announcement_id=$1 AND parent_comment_id=$2"
		args = []interface{}{announcementID, *parentCommentID}
	} else {
		query = "SELECT id, announcement_id, user_id, parent_comment_id, content, created_at FROM announcement_comments WHERE announcement_id=$1 AND parent_comment_id IS NULL"
		args = []interface{}{announcementID}
	}
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []*domain.AnnouncementComment
	for rows.Next() {
		var c domain.AnnouncementComment
		err = rows.Scan(&c.ID, &c.AnnouncementID, &c.UserID, &c.ParentCommentID, &c.Content, &c.CreatedAt)
		if err != nil {
			return nil, err
		}
		comments = append(comments, &c)
	}
	return comments, nil
}

func (r *announcementRepo) DeleteComment(ctx context.Context, commentID int) error {
    tx, err := r.db.BeginTx(ctx, nil)
    if err != nil {
        return err
    }
    defer tx.Rollback()

    // 1. Get announcement_id from the comment itself
    var announcementID int
    err = tx.QueryRowContext(ctx,
        "SELECT announcement_id FROM announcement_comments WHERE id=$1",
        commentID,
    ).Scan(&announcementID)
    if err != nil {
        return err // returns ErrNoRows if comment doesn't exist
    }

    // 2. COUNT(*) without GROUP BY always returns exactly one row (0 if no matches)
    var replyCount int
    err = tx.QueryRowContext(ctx,
        "SELECT COUNT(*) FROM announcement_comments WHERE parent_comment_id=$1",
        commentID,
    ).Scan(&replyCount)
    if err != nil {
        return err
    }

    // 3. Delete the comment (replies are cascade deleted by FK)
    _, err = tx.ExecContext(ctx,
        "DELETE FROM announcement_comments WHERE id=$1",
        commentID,
    )
    if err != nil {
        return err
    }

    // 4. Decrement count for the comment + its replies
    _, err = tx.ExecContext(ctx,
        "UPDATE announcements SET comments_count = comments_count - 1 - $1 WHERE id=$2",
        replyCount, announcementID,
    )
    if err != nil {
        return err
    }

    return tx.Commit()
}

func (r *announcementRepo) EditComment(ctx context.Context, commentID int, content string) (*domain.AnnouncementComment, error) {
	_, err := r.db.ExecContext(ctx, "UPDATE announcement_comments SET content=$1, is_edited=true, updated_at=NOW() WHERE id=$2", content, commentID)
	if err != nil {
		return nil, err
	}
	var c domain.AnnouncementComment
	err = r.db.QueryRowContext(ctx, "SELECT id, announcement_id, user_id, parent_comment_id, content, created_at FROM announcement_comments WHERE id=$1", commentID).Scan(&c.ID, &c.AnnouncementID, &c.UserID, &c.ParentCommentID, &c.Content, &c.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &c, nil
}

