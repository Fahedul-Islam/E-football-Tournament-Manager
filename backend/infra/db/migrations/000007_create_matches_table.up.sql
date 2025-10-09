CREATE TABLE IF NOT EXISTS matches (
    id SERIAL PRIMARY KEY,
    tournament_id INT NOT NULL REFERENCES tournaments(id) ON DELETE CASCADE,
    group_id INT NOT NULL REFERENCES groups(id) ON DELETE CASCADE,
    round VARCHAR(50), -- e.g., "Group Stage", "Round of 16", "Quarterfinals","Semifinals", "Final"
    participant_a_id INT NOT NULL REFERENCES participants(id) ON DELETE CASCADE,
    participant_b_id INT NOT NULL REFERENCES participants(id) ON DELETE CASCADE,
    participant_a_score INT DEFAULT 0,
    participant_b_score INT DEFAULT 0,
    match_date TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    status VARCHAR(20) NOT NULL DEFAULT 'scheduled',
    UNIQUE(group_id, participant_a_id, participant_b_id)
);