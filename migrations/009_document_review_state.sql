ALTER TABLE documents ADD COLUMN review_status TEXT NOT NULL DEFAULT 'pending_review';
ALTER TABLE documents ADD COLUMN reviewed_at TEXT NOT NULL DEFAULT '';
ALTER TABLE documents ADD COLUMN review_note TEXT NOT NULL DEFAULT '';

CREATE INDEX IF NOT EXISTS idx_documents_review_status
ON documents(review_status);