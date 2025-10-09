CREATE TABLE IF NOT EXISTS tournament_stages (
    id SERIAL PRIMARY KEY,
    tournament_id INT REFERENCES tournaments(id) ON DELETE CASCADE,
    stage VARCHAR(30),  -- group, knockout, quarter, semi, final
    is_completed BOOLEAN DEFAULT FALSE
);
