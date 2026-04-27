ALTER TABLE document_events ADD COLUMN review_status TEXT NOT NULL DEFAULT 'pending';
ALTER TABLE document_events ADD COLUMN reviewed_at TEXT NOT NULL DEFAULT '';
ALTER TABLE document_events ADD COLUMN resolved_at TEXT NOT NULL DEFAULT '';
ALTER TABLE document_events ADD COLUMN resolution_note TEXT NOT NULL DEFAULT '';

CREATE INDEX IF NOT EXISTS idx_document_events_review_status
ON document_events(review_status);