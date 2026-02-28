CREATE TABLE IF NOT EXISTS player_stats (
    id SERIAL PRIMARY KEY,
    tournament_id INT NOT NULL REFERENCES tournaments(id) ON DELETE CASCADE,
    group_id INT REFERENCES groups(id) ON DELETE CASCADE,  -- nullable for league
    participant_id INT NOT NULL REFERENCES participants(id) ON DELETE CASCADE,
    matches_played INT DEFAULT 0 CHECK (matches_played >= 0),
    wins INT DEFAULT 0 CHECK (wins >= 0),
    draws INT DEFAULT 0 CHECK (draws >= 0),
    losses INT DEFAULT 0 CHECK (losses >= 0),
    goals_scored INT DEFAULT 0 CHECK (goals_scored >= 0),
    goals_conceded INT DEFAULT 0 CHECK (goals_conceded >= 0),
    goal_difference INT DEFAULT 0,
    points INT DEFAULT 0 CHECK (points >= 0),
    UNIQUE(group_id, participant_id)  -- For ON CONFLICT in group stage
);

-- For league stats uniqueness (group_id IS NULL)
CREATE UNIQUE INDEX IF NOT EXISTS idx_player_stats_league_participant 
ON player_stats(tournament_id, participant_id) WHERE group_id IS NULL;

CREATE INDEX IF NOT EXISTS idx_player_stats_tournament ON player_stats(tournament_id);