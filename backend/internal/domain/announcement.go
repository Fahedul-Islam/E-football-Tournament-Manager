package domain

import "time"

type Announcement struct {
	ID               int       `json:"id"`
	TournamentID     int       `json:"tournament_id"`
	AuthorID         int       `json:"author_id"`
	Title            string    `json:"title"`
	Content          string    `json:"content"`
	AnnouncementType string    `json:"announcement_type"`
	IsPinned         bool      `json:"is_pinned"`
	IsCommentable    bool      `json:"is_commentable"`
	LikesCount       int       `json:"likes_count"`
	DislikesCount    int       `json:"dislikes_count"`
	CommentsCount    int       `json:"comments_count"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

type AnnouncementReaction struct {
	ID             int       `json:"id"`
	AnnouncementID int       `json:"announcement_id"`
	UserID         int       `json:"user_id"`
	ReactionType   string    `json:"reaction_type"` // "like" or "dislike"
	CreatedAt      time.Time `json:"created_at"`
}

type AnnouncementComment struct {
	ID              int       `json:"id"`
	AnnouncementID  int       `json:"announcement_id"`
	UserID          int       `json:"user_id"`
	ParentCommentID *int      `json:"parent_comment_id,omitempty"`
	Content         string    `json:"content"`
	LikesCount      int       `json:"likes_count"`
	DislikesCount   int       `json:"dislikes_count"`
	IsEdited        bool      `json:"is_edited"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

type AnnouncementCommentReaction struct {
	ID           int       `json:"id"`
	CommentID    int       `json:"comment_id"`
	UserID       int       `json:"user_id"`
	ReactionType string    `json:"reaction_type"` // "like" or "dislike"
	CreatedAt    time.Time `json:"created_at"`
}

type AnnouncementSeen struct {
	ID             int       `json:"id"`
	AnnouncementID int       `json:"announcement_id"`
	UserID         int       `json:"user_id"`
	IsSeen         bool      `json:"is_seen"`
	SeenAt         time.Time `json:"seen_at"`
}

type AnnouncementCreateRequest struct {
	Title            string `json:"title" validate:"required"`
	Content          string `json:"content" validate:"required"`
	AnnouncementType string `json:"announcement_type" validate:"required,oneof=general match update"`
	IsPinned         bool   `json:"is_pinned"`
	IsCommentable    bool   `json:"is_commentable"`
}

type CommentCreateRequest struct {
	Content         string `json:"content" validate:"required"`
	ParentCommentID *int   `json:"parent_comment_id,omitempty"`
}
