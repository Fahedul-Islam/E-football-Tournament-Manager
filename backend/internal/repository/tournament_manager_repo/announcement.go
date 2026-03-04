package tournamentmanagerrepo

import (
	"context"
	"tournament-manager/internal/domain"
)

func (r *tournamentManagerRepo) CreateAnnouncement(ctx context.Context, announcement *domain.Announcement) (*domain.Announcement, error) {
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

func (r *tournamentManagerRepo) GetAnnouncements(ctx context.Context, tournamentID int) ([]*domain.Announcement, error) {
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

func (r *tournamentManagerRepo) GetAnnouncementByID(ctx context.Context, tournamentID int, announcementID int) (*domain.Announcement, error) {
	query := `SELECT id, tournament_id, author_id, title, content, announcement_type, is_pinned, is_commentable, likes_count, dislikes_count, comments_count, created_at, updated_at FROM announcements WHERE tournament_id = $1 AND id = $2`
	row := r.db.QueryRowContext(ctx, query, tournamentID, announcementID)

	var a domain.Announcement
	if err := row.Scan(&a.ID, &a.TournamentID, &a.AuthorID, &a.Title, &a.Content, &a.AnnouncementType, &a.IsPinned, &a.IsCommentable, &a.LikesCount, &a.DislikesCount, &a.CommentsCount, &a.CreatedAt, &a.UpdatedAt); err != nil {
		return nil, err
	}
	return &a, nil
}

func (r *tournamentManagerRepo) UpdateAnnouncement(ctx context.Context, announcement *domain.Announcement) (*domain.Announcement, error) {
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

func (r *tournamentManagerRepo) DeleteAnnouncement(ctx context.Context, tournamentID int, announcementID int) error {
	query := `DELETE FROM announcements WHERE id = $1 AND tournament_id = $2`
	_, err := r.db.ExecContext(ctx, query, announcementID, tournamentID)
	return err
}