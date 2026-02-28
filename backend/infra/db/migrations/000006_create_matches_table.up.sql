CREATE TABLE IF NOT EXISTS matches (
    id SERIAL PRIMARY KEY,
    tournament_id INT NOT NULL REFERENCES tournaments(id) ON DELETE CASCADE,
    group_id INT REFERENCES groups(id) ON DELETE CASCADE,  -- nullable for knockout/league
    round VARCHAR(50) NOT NULL,
    participant_a_id INT NOT NULL REFERENCES participants(id) ON DELETE CASCADE,
    participant_b_id INT NOT NULL REFERENCES participants(id) ON DELETE CASCADE,
    participant_a_score INT DEFAULT 0 CHECK (participant_a_score >= 0),
    participant_b_score INT DEFAULT 0 CHECK (participant_b_score >= 0),
    match_date TIMESTAMPTZ,
    status VARCHAR(20) NOT NULL DEFAULT 'scheduled' CHECK (status IN ('scheduled', 'completed', 'cancelled')),
    created_at TIMESTAMPTZ DEFAULT NOW(),
    CHECK (participant_a_id <> participant_b_id)
);

-- Unique index to prevent duplicate matchups (order-independent)
CREATE UNIQUE INDEX idx_matches_unique_pair 
ON matches(tournament_id, round, LEAST(participant_a_id, participant_b_id), GREATEST(participant_a_id, participant_b_id));

CREATE INDEX IF NOT EXISTS idx_matches_tournament ON matches(tournament_id);
CREATE INDEX IF NOT EXISTS idx_matches_group ON matches(group_id);
CREATE INDEX IF NOT EXISTS idx_matches_round ON matches(tournament_id, round);
CREATE INDEX IF NOT EXISTS idx_matches_status ON matches(tournament_id, status);