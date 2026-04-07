-- Only include Up migration commands here for psql execution
ALTER TABLE tasks ADD COLUMN periodicity_type TEXT NOT NULL DEFAULT 'no';
ALTER TABLE tasks ADD COLUMN periodicity_value TEXT NOT NULL DEFAULT '';
ALTER TABLE tasks ADD COLUMN last_completed_at TIMESTAMPTZ;
