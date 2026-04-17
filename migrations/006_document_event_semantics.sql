ALTER TABLE document_events ADD COLUMN add_extra_day INTEGER NOT NULL DEFAULT 0;
ALTER TABLE document_events ADD COLUMN calendar_scope TEXT NOT NULL DEFAULT 'madrid';
ALTER TABLE document_events ADD COLUMN computation TEXT NOT NULL DEFAULT '';