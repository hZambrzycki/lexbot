CREATE TABLE IF NOT EXISTS document_metadata (
    document_id TEXT PRIMARY KEY,
    document_type TEXT NOT NULL,
    legal_area TEXT NOT NULL,
    analyzed_at TEXT NOT NULL,
    FOREIGN KEY (document_id) REFERENCES documents(id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_document_metadata_legal_area
    ON document_metadata(legal_area);

CREATE INDEX IF NOT EXISTS idx_document_metadata_document_type
    ON document_metadata(document_type);