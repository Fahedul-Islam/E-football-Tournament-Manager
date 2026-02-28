CREATE TABLE IF NOT EXISTS tournaments (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    description TEXT,
    tournament_type VARCHAR(50) NOT NULL 
        CHECK (tournament_type IN ('knockout', 'league', 'group+knockout')),
    max_players INT NOT NULL CHECK (max_players > 0),
    created_by INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    start_date TIMESTAMPTZ NOT NULL,
    end_date   TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    CHECK (end_date >= start_date)
);

CREATE INDEX IF NOT EXISTS idx_tournaments_created_by 
ON tournaments(created_by);