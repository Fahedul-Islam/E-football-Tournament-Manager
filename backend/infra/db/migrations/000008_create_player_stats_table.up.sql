CREATE TABLE IF NOT EXISTS player_stats (
    id SERIAL PRIMARY KEY,
    group_id INT REFERENCES groups(id) ON DELETE CASCADE,
    participant_id INT REFERENCES participants(id) ON DELETE CASCADE,
    matches_played INT DEFAULT 0,
    wins INT DEFAULT 0,
    draws INT DEFAULT 0,
    losses INT DEFAULT 0,
    goals_scored INT DEFAULT 0,
    goals_conceded INT DEFAULT 0,
    goal_difference INT DEFAULT 0,
    points INT DEFAULT 0
);
