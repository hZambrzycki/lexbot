CREATE TABLE IF NOT EXISTS document_events (
    id TEXT PRIMARY KEY,
    document_id TEXT NOT NULL,
    event_type TEXT NOT NULL,
    event_date TEXT NOT NULL,
    source_text TEXT NOT NULL,
    created_at TEXT NOT NULL,
    FOREIGN KEY (document_id) REFERENCES documents(id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_document_events_document_id
    ON document_events(document_id);

CREATE INDEX IF NOT EXISTS idx_document_events_event_date
    ON document_events(event_date);

CREATE INDEX IF NOT EXISTS idx_document_events_event_type
    ON document_events(event_type);