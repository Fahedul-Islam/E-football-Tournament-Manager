CREATE TABLE IF NOT EXISTS group_teams (
    id SERIAL PRIMARY KEY,
    group_id INT NOT NULL REFERENCES groups(id) ON DELETE CASCADE,
    participant_id INT NOT NULL REFERENCES participants(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE(group_id, participant_id),
    UNIQUE(participant_id)  -- each participant can only be in one group
);

CREATE INDEX IF NOT EXISTS idx_group_teams_group ON group_teams(group_id);