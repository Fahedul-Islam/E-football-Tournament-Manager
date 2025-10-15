ALTER TABLE player_stats
ADD CONSTRAINT unique_group_participant UNIQUE (group_id, participant_id);
