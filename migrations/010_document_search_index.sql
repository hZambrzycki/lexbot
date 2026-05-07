CREATE VIRTUAL TABLE IF NOT EXISTS document_search_index USING fts5(
    document_id UNINDEXED,
    case_file_id UNINDEXED,
    original_name,
    content,
    document_type,
    legal_area,
    tokenize = 'unicode61 remove_diacritics 2'
);