package sqlite

import (
	"context"
	"database/sql"
	"strings"

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
	filters querymodels.SearchDocumentFilters,
	limit int,
) ([]querymodels.SearchDocumentResult, error) {
	ftsQuery := searchinfra.BuildFTS5Query(query)

	where := make([]string, 0)
	args := make([]any, 0)

	if ftsQuery != "" {
		where = append(where, "document_search_index MATCH ?")
		args = append(args, ftsQuery)
	}

	if caseFileID != "" {
		where = append(where, "document_search_index.case_file_id = ?")
		args = append(args, caseFileID)
	}

	if filters.DocumentType != "" {
		where = append(where, "document_search_index.document_type = ?")
		args = append(args, filters.DocumentType)
	}

	if filters.LegalArea != "" {
		where = append(where, "document_search_index.legal_area = ?")
		args = append(args, filters.LegalArea)
	}

	if filters.ReviewStatus != "" {
		where = append(where, "documents.review_status = ?")
		args = append(args, filters.ReviewStatus)
	}

	if filters.DocType != "" {
		switch filters.DocType {
		case "pdf":
			where = append(where, "documents.mime_type = ?")
			args = append(args, "application/pdf")
		case "docx":
			where = append(where, "documents.mime_type = ?")
			args = append(args, "application/vnd.openxmlformats-officedocument.wordprocessingml.document")
		case "txt":
			where = append(where, "documents.mime_type = ?")
			args = append(args, "text/plain")
		case "md":
			where = append(where, "documents.mime_type IN (?, ?)")
			args = append(args, "text/markdown", "text/x-markdown")
		}
	}

	if filters.Has != "" {
		switch filters.Has {
		case "events":
			where = append(where, `EXISTS (
				SELECT 1
				FROM document_events
				WHERE document_events.document_id = document_search_index.document_id
			)`)
		case "no_events":
			where = append(where, `NOT EXISTS (
				SELECT 1
				FROM document_events
				WHERE document_events.document_id = document_search_index.document_id
			)`)
		case "text":
			where = append(where, `EXISTS (
				SELECT 1
				FROM document_contents
				WHERE document_contents.document_id = document_search_index.document_id
				AND TRIM(document_contents.content) <> ''
			)`)
		case "no_text":
			where = append(where, `NOT EXISTS (
				SELECT 1
				FROM document_contents
				WHERE document_contents.document_id = document_search_index.document_id
				AND TRIM(document_contents.content) <> ''
			)`)
		}
	}

	if len(where) == 0 {
		return []querymodels.SearchDocumentResult{}, nil
	}

	args = append(args, limit)

	sqlQuery := `
		SELECT
			document_search_index.document_id,
			document_search_index.original_name,
			document_search_index.case_file_id,
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
		JOIN documents
			ON documents.id = document_search_index.document_id
		WHERE ` + strings.Join(where, " AND ") + `
		ORDER BY score
		LIMIT ?
	`

	rows, err := r.db.QueryContext(ctx, sqlQuery, args...)
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
