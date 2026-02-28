CREATE TABLE participants (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    tournament_id INT NOT NULL REFERENCES tournaments(id) ON DELETE CASCADE,
    team_name VARCHAR(100) NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'pending' CHECK (status IN ('pending', 'approved', 'rejected')),
    created_at TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE(user_id, tournament_id)
);

CREATE INDEX idx_participants_tournament ON participants(tournament_id);
CREATE INDEX idx_participants_status ON participants(tournament_id, status);