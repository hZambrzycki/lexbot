package sqlite

import (
	"context"
	"database/sql"

	"lexbox/internal/application/querymodels"
	searchinfra "lexbox/internal/infrastructure/search"
)

type DocumentSearchIndexRepository struct {
	db *sql.DB
}

func NewDocumentSearchIndexRepository(db *sql.DB) *DocumentSearchIndexRepository {
	return &DocumentSearchIndexRepository{
		db: db,
	}
}

func (r *DocumentSearchIndexRepository) UpsertDocument(
	ctx context.Context,
	documentID string,
	caseFileID string,
	originalName string,
	content string,
	documentType string,
	legalArea string,
) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if _, err := tx.ExecContext(ctx, `
		DELETE FROM document_search_index
		WHERE document_id = ?
	`, documentID); err != nil {
		return err
	}

	if _, err := tx.ExecContext(ctx, `
		INSERT INTO document_search_index (
			document_id,
			case_file_id,
			original_name,
			content,
			document_type,
			legal_area
		)
		VALUES (?, ?, ?, ?, ?, ?)
	`,
		documentID,
		caseFileID,
		originalName,
		content,
		documentType,
		legalArea,
	); err != nil {
		return err
	}

	return tx.Commit()
}

func (r *DocumentSearchIndexRepository) DeleteDocument(
	ctx context.Context,
	documentID string,
) error {
	_, err := r.db.ExecContext(ctx, `
		DELETE FROM document_search_index
		WHERE document_id = ?
	`, documentID)

	return err
}

func (r *DocumentSearchIndexRepository) Search(
	ctx context.Context,
	query string,
	caseFileID string,
	limit int,
) ([]querymodels.SearchDocumentResult, error) {
	ftsQuery := searchinfra.BuildFTS5Query(query)
	if ftsQuery == "" {
		return []querymodels.SearchDocumentResult{}, nil
	}

	var (
		rows *sql.Rows
		err  error
	)

	if caseFileID != "" {
		rows, err = r.db.QueryContext(ctx, `
			SELECT
				document_id,
				original_name,
				case_file_id,
				snippet(
					document_search_index,
					3,
					'[',
					']',
					'...',
					24
				) AS snippet,
				bm25(document_search_index) AS score
			FROM document_search_index
			WHERE document_search_index MATCH ?
			AND case_file_id = ?
			ORDER BY score
			LIMIT ?
		`, ftsQuery, caseFileID, limit)
	} else {
		rows, err = r.db.QueryContext(ctx, `
			SELECT
				document_id,
				original_name,
				case_file_id,
				snippet(
					document_search_index,
					3,
					'[',
					']',
					'...',
					24
				) AS snippet,
				bm25(document_search_index) AS score
			FROM document_search_index
			WHERE document_search_index MATCH ?
			ORDER BY score
			LIMIT ?
		`, ftsQuery, limit)
	}

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	results := make([]querymodels.SearchDocumentResult, 0)

	for rows.Next() {
		var (
			result querymodels.SearchDocumentResult
			score  float64
		)

		if err := rows.Scan(
			&result.DocumentID,
			&result.OriginalName,
			&result.CaseFileID,
			&result.Snippet,
			&score,
		); err != nil {
			return nil, err
		}

		result.Score = int(score * -1000)

		results = append(results, result)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return results, nil
}
