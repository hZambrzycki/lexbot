CREATE TABLE IF NOT EXISTS document_contents (
    document_id TEXT PRIMARY KEY,
    content TEXT NOT NULL,
    extracted_at TEXT NOT NULL,
    FOREIGN KEY (document_id) REFERENCES documents(id)
);