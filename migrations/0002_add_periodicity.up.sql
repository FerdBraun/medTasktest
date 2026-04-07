-- Idempotent migration to add periodicity columns
DO $$ 
BEGIN 
    IF NOT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name='tasks' AND column_name='periodicity_type') THEN
        ALTER TABLE tasks ADD COLUMN periodicity_type TEXT NOT NULL DEFAULT 'no';
    END IF;

    IF NOT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name='tasks' AND column_name='periodicity_value') THEN
        ALTER TABLE tasks ADD COLUMN periodicity_value TEXT NOT NULL DEFAULT '';
    END IF;

    IF NOT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name='tasks' AND column_name='last_completed_at') THEN
        ALTER TABLE tasks ADD COLUMN last_completed_at TIMESTAMPTZ;
    END IF;
END $$;
