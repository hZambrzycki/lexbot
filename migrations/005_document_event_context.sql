ALTER TABLE document_events ADD COLUMN anchor_date TEXT;
ALTER TABLE document_events ADD COLUMN date_kind TEXT;
ALTER TABLE document_events ADD COLUMN anchor_source TEXT;
ALTER TABLE document_events ADD COLUMN relative_days INTEGER NOT NULL DEFAULT 0;
ALTER TABLE document_events ADD COLUMN is_business_days INTEGER NOT NULL DEFAULT 0;
ALTER TABLE document_events ADD COLUMN trigger_text TEXT;