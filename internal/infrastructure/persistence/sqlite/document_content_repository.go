package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"sort"
	"strings"
	"time"

	"lexbox/internal/application/querymodels"
	"lexbox/internal/domain/shared"
	"lexbox/internal/infrastructure/search"
)

type DocumentContentRepository struct {
	db *sql.DB
}

func NewDocumentContentRepository(db *sql.DB) *DocumentContentRepository {
	return &DocumentContentRepository{db: db}
}

func (r *DocumentContentRepository) Save(ctx context.Context, documentID string, content string) error {
	const query = `
		INSERT INTO document_contents (document_id, content, extracted_at)
		VALUES (?, ?, ?)
		ON CONFLICT(document_id) DO UPDATE SET
			content = excluded.content,
			extracted_at = excluded.extracted_at
	`

	_, err := r.db.ExecContext(
		ctx,
		query,
		documentID,
		content,
		time.Now().Format(time.RFC3339),
	)

	return err
}

func (r *DocumentContentRepository) GetByDocumentID(ctx context.Context, documentID string) (string, error) {
	const query = `
		SELECT content
		FROM document_contents
		WHERE document_id = ?
	`

	var content string
	err := r.db.QueryRowContext(ctx, query, documentID).Scan(&content)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", shared.ErrNotFound
		}
		return "", err
	}

	return content, nil
}

func (r *DocumentContentRepository) SearchByText(ctx context.Context, query string, limit int) ([]querymodels.SearchDocumentResult, error) {
	terms := search.SplitTerms(query)
	if len(terms) == 0 {
		return []querymodels.SearchDocumentResult{}, nil
	}

	sqlQuery, args := buildSearchQuery("", limit)

	rows, err := r.db.QueryContext(ctx, sqlQuery, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return scanAndRankSearchResults(rows, query, terms, limit)
}

func (r *DocumentContentRepository) SearchByTextByCaseFile(ctx context.Context, caseFileID string, query string, limit int) ([]querymodels.SearchDocumentResult, error) {
	terms := search.SplitTerms(query)
	if len(terms) == 0 {
		return []querymodels.SearchDocumentResult{}, nil
	}

	sqlQuery, args := buildSearchQuery(caseFileID, limit)

	rows, err := r.db.QueryContext(ctx, sqlQuery, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return scanAndRankSearchResults(rows, query, terms, limit)
}

func buildSearchQuery(caseFileID string, limit int) (string, []any) {
	var sb strings.Builder
	args := make([]any, 0, 2)

	sb.WriteString(`
		SELECT
			d.id,
			d.original_name,
			d.case_file_id,
			dc.content
		FROM document_contents dc
		INNER JOIN documents d ON d.id = dc.document_id
	`)

	if caseFileID != "" {
		sb.WriteString(`
			WHERE d.case_file_id = ?
		`)
		args = append(args, caseFileID)
	}

	sb.WriteString(`
		ORDER BY d.id DESC
		LIMIT ?
	`)

	fetchLimit := limit * 20
	if fetchLimit < 100 {
		fetchLimit = 100
	}

	args = append(args, fetchLimit)

	return sb.String(), args
}

type rankedSearchResult struct {
	result querymodels.SearchDocumentResult
	score  int
}

func scanAndRankSearchResults(rows *sql.Rows, rawQuery string, terms []string, limit int) ([]querymodels.SearchDocumentResult, error) {
	ranked := make([]rankedSearchResult, 0)

	for rows.Next() {
		var (
			documentID   string
			originalName string
			caseFileID   string
			content      string
		)

		if err := rows.Scan(&documentID, &originalName, &caseFileID, &content); err != nil {
			return nil, err
		}

		score := search.ComputeScore(content, terms)
		if score <= 0 {
			continue
		}

		ranked = append(ranked, rankedSearchResult{
			result: querymodels.SearchDocumentResult{
				DocumentID:   documentID,
				OriginalName: originalName,
				CaseFileID:   caseFileID,
				Snippet:      search.BuildSnippet(content, rawQuery, terms, 180),
				Score:        score,
			},
			score: score,
		})
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	sort.SliceStable(ranked, func(i, j int) bool {
		if ranked[i].score == ranked[j].score {
			return ranked[i].result.DocumentID > ranked[j].result.DocumentID
		}
		return ranked[i].score > ranked[j].score
	})

	if len(ranked) > limit {
		ranked = ranked[:limit]
	}

	results := make([]querymodels.SearchDocumentResult, 0, len(ranked))
	for _, item := range ranked {
		results = append(results, item.result)
	}

	return results, nil
}
