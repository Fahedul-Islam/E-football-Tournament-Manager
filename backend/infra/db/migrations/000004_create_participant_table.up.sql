CREATE TABLE IF NOT EXISTS participants (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    tournament_id INT NOT NULL REFERENCES tournaments(id) ON DELETE CASCADE,
    team_name VARCHAR(100) NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'pending',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(user_id, tournament_id)
);