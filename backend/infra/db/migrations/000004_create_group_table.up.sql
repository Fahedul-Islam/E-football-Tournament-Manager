CREATE TABLE IF NOT EXISTS groups (
    id SERIAL PRIMARY KEY,
    name VARCHAR(10) NOT NULL,
    tournament_id INT NOT NULL REFERENCES tournaments(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE(name, tournament_id)
);

CREATE INDEX IF NOT EXISTS idx_groups_tournament ON groups(tournament_id);