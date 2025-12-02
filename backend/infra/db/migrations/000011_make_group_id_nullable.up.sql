-- Make group_id nullable
ALTER TABLE matches
    ALTER COLUMN group_id DROP NOT NULL;

-- Ensure round always exists
ALTER TABLE matches
    ALTER COLUMN round SET NOT NULL;
