CREATE TABLE IF NOT EXISTS announcement_reactions (
    id SERIAL PRIMARY KEY,
    announcement_id INTEGER NOT NULL REFERENCES announcements(id) ON DELETE CASCADE,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    reaction_type VARCHAR(50) NOT NULL CHECK (reaction_type IN ('like', 'dislike')),
    created_at TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE(announcement_id, user_id)
);

CREATE INDEX idx_reaction_announcement_id ON announcement_reactions(announcement_id);
CREATE INDEX idx_reaction_user_id ON announcement_reactions(user_id);