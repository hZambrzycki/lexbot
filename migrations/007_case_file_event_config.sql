ALTER TABLE case_files ADD COLUMN calendar_scope TEXT NOT NULL DEFAULT 'madrid';
ALTER TABLE case_files ADD COLUMN august_non_business INTEGER NOT NULL DEFAULT 1;