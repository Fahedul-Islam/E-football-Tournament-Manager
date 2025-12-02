-- Revert group_id to NOT NULL
ALTER TABLE matches
    ALTER COLUMN group_id SET NOT NULL;

-- Revert round to nullable
ALTER TABLE matches
    ALTER COLUMN round DROP NOT NULL;
