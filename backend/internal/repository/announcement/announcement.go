package announcement

import (
	"context"
	"tournament-manager/internal/domain"
)

func (r *announcementRepo) CreateAnnouncement(ctx context.Context, announcement *domain.Announcement) (*domain.Announcement, error) {
	query := `INSERT INTO announcements (tournament_id, author_id, title, content, announcement_type, is_pinned, is_commentable) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id`
	err := r.db.QueryRowContext(ctx, query,
		announcement.TournamentID,
		announcement.AuthorID,
		announcement.Title,
		announcement.Content,
		announcement.AnnouncementType,
		announcement.IsPinned,
		announcement.IsCommentable,
	).Scan(&announcement.ID)
	if err != nil {
		return nil, err
	}
	return announcement, nil
}

func (r *announcementRepo) GetAnnouncements(ctx context.Context, tournamentID int) ([]*domain.Announcement, error) {
	query := `SELECT id, tournament_id, author_id, title, content, announcement_type, is_pinned, is_commentable, likes_count, dislikes_count, comments_count, created_at, updated_at FROM announcements WHERE tournament_id = $1 ORDER BY created_at DESC`
	rows, err := r.db.QueryContext(ctx, query, tournamentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var announcements []*domain.Announcement
	for rows.Next() {
		var a domain.Announcement
		if err := rows.Scan(&a.ID, &a.TournamentID, &a.AuthorID, &a.Title, &a.Content, &a.AnnouncementType, &a.IsPinned, &a.IsCommentable, &a.LikesCount, &a.DislikesCount, &a.CommentsCount, &a.CreatedAt, &a.UpdatedAt); err != nil {
			return nil, err
		}
		announcements = append(announcements, &a)
	}
	return announcements, nil
}

func (r *announcementRepo) GetAnnouncementByID(ctx context.Context, tournamentID int, announcementID int, userID int) (*domain.Announcement, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	// Mark the announcement as seen by the user, if not already marked
	_, err = tx.ExecContext(ctx, "INSERT INTO announcement_seen (announcement_id, user_id, is_seen, seen_at) VALUES ($1, $2, true, NOW()) ON CONFLICT (announcement_id, user_id) DO NOTHING", announcementID, userID)
	if err != nil {
		return nil, err
	}

	query := `SELECT id, tournament_id, author_id, title, content, announcement_type, is_pinned, is_commentable, likes_count, dislikes_count, comments_count, created_at, updated_at FROM announcements WHERE tournament_id = $1 AND id = $2`
	var a domain.Announcement
	err = tx.QueryRowContext(ctx, query, tournamentID, announcementID).
		Scan(&a.ID, &a.TournamentID, &a.AuthorID, &a.Title, &a.Content,
			&a.AnnouncementType, &a.IsPinned, &a.IsCommentable,
			&a.LikesCount, &a.DislikesCount, &a.CommentsCount,
			&a.CreatedAt, &a.UpdatedAt)

	if err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return &a, nil
}

func (r *announcementRepo) UpdateAnnouncement(ctx context.Context, announcement *domain.Announcement) (*domain.Announcement, error) {
	query := `UPDATE announcements SET title = $1, content = $2, announcement_type = $3, is_pinned = $4, is_commentable = $5, updated_at = NOW() WHERE id = $6 AND tournament_id = $7 RETURNING created_at`
	err := r.db.QueryRowContext(ctx, query,
		announcement.Title,
		announcement.Content,
		announcement.AnnouncementType,
		announcement.IsPinned,
		announcement.IsCommentable,
		announcement.ID,
		announcement.TournamentID,
	).Scan(&announcement.CreatedAt)
	if err != nil {
		return nil, err
	}
	return announcement, nil
}

func (r *announcementRepo) DeleteAnnouncement(ctx context.Context, tournamentID int, announcementID int) error {
	query := `DELETE FROM announcements WHERE id = $1 AND tournament_id = $2`
	_, err := r.db.ExecContext(ctx, query, announcementID, tournamentID)
	return err
}

func (r *announcementRepo) GetParticipantsAnnouncementSeenStatus(ctx context.Context, tournamentID int, announcementID int, userID int) (*[]domain.Participant, error) {
	query := `SELECT 
			p.id,
			p.user_id,
			p.tournament_id,
			p.team_name,
			p.status,
			p.created_at
			FROM participants p
			JOIN announcement_seen asn
			ON p.user_id = asn.user_id
			WHERE p.tournament_id = $1
			AND asn.announcement_id = $2
			AND asn.is_seen = true;`

	rows, err := r.db.QueryContext(ctx, query, tournamentID, announcementID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var participants []domain.Participant
	for rows.Next() {
		var p domain.Participant
		if err := rows.Scan(&p.ID, &p.UserID, &p.TournamentID, &p.TeamName, &p.Status, &p.CreatedAt); err != nil {
			return nil, err
		}
		participants = append(participants, p)
	}
	return &participants, nil
}

func (r *announcementRepo) AddAnnouncementNotification(ctx context.Context, announcementID int, message string, participants []*domain.Participant) error {
	// Insert notification for each participant
	for _, p := range participants {
		_, err := r.db.ExecContext(ctx, "INSERT INTO notifications (user_id, notification_type, reference_id, message, is_read, created_at) VALUES ($1, $2, $3, $4, false, NOW())", p.UserID, "announcement", announcementID, message)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *announcementRepo) GetNotifications(ctx context.Context, userID int, page int) ([]*domain.Notification, error) {
	offset := (page - 1) * 20
	query := `SELECT id, user_id, notification_type, reference_id, message, is_read, created_at 
			FROM notifications 
			WHERE user_id = $1 
			ORDER BY created_at DESC 
			LIMIT 20 OFFSET $2`
	rows, err := r.db.QueryContext(ctx, query, userID, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var notifications []*domain.Notification
	for rows.Next() {
		var n domain.Notification
		if err := rows.Scan(&n.ID, &n.UserID, &n.NotificationType, &n.ReferenceID, &n.Message, &n.IsRead, &n.CreatedAt); err != nil {
			return nil, err
		}
		notifications = append(notifications, &n)
	}
	return notifications, nil
}

func (r *announcementRepo) MarkNotificationAsRead(ctx context.Context, notificationID int, userID int) error {
	query := `UPDATE notifications SET is_read = true WHERE id = $1 AND user_id = $2`
	_, err := r.db.ExecContext(ctx, query, notificationID, userID)
	return err
}

func (r *announcementRepo) MarkAllNotificationsAsRead(ctx context.Context, userID int) error {
	query := `UPDATE notifications SET is_read = true WHERE user_id = $1 and is_read = false`
	_, err := r.db.ExecContext(ctx, query, userID)
	return err
}

