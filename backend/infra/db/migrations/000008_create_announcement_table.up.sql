CREATE TABLE IF NOT EXISTS announcements (
    id SERIAL PRIMARY KEY,
    tournament_id INTEGER NOT NULL REFERENCES tournaments(id) ON DELETE CASCADE,
    author_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    title VARCHAR(255) NOT NULL,
    content TEXT NOT NULL,
    announcement_type VARCHAR(50) NOT NULL DEFAULT 'general' CHECK (announcement_type IN ('general', 'update', 'reminder', 'result','urgent', 'other')),
    is_pinned BOOLEAN NOT NULL DEFAULT FALSE,
    is_commentable BOOLEAN NOT NULL DEFAULT TRUE,
    likes_count INTEGER NOT NULL DEFAULT 0,
    dislikes_count INTEGER NOT NULL DEFAULT 0,
    comments_count INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_announcement_tournament_id ON announcements(tournament_id);
CREATE INDEX idx_announcement_author_id ON announcements(author_id);