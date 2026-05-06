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
	terms := splitSearchTerms(query)
	if len(terms) == 0 {
		return []querymodels.SearchDocumentResult{}, nil
	}

	sqlQuery, args := buildSearchQuery("", terms, limit)

	rows, err := r.db.QueryContext(ctx, sqlQuery, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return scanAndRankSearchResults(rows, query, terms, limit)
}

func (r *DocumentContentRepository) SearchByTextByCaseFile(ctx context.Context, caseFileID string, query string, limit int) ([]querymodels.SearchDocumentResult, error) {
	terms := splitSearchTerms(query)
	if len(terms) == 0 {
		return []querymodels.SearchDocumentResult{}, nil
	}

	sqlQuery, args := buildSearchQuery(caseFileID, terms, limit)

	rows, err := r.db.QueryContext(ctx, sqlQuery, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return scanAndRankSearchResults(rows, query, terms, limit)
}

func buildSearchQuery(caseFileID string, terms []string, limit int) (string, []any) {
	var sb strings.Builder
	args := make([]any, 0, len(terms)+2)

	sb.WriteString(`
		SELECT
			d.id,
			d.original_name,
			d.case_file_id,
			dc.content
		FROM document_contents dc
		INNER JOIN documents d ON d.id = dc.document_id
		WHERE
	`)

	if caseFileID != "" {
		sb.WriteString(` d.case_file_id = ? AND (`)
		args = append(args, caseFileID)
	} else {
		sb.WriteString(`(`)
	}

	for i, term := range terms {
		if i > 0 {
			sb.WriteString(` OR `)
		}
		sb.WriteString(`lower(dc.content) LIKE '%' || lower(?) || '%'`)
		args = append(args, term)
	}

	sb.WriteString(`)`)
	sb.WriteString(`
		ORDER BY d.id DESC
		LIMIT ?
	`)

	// Pedimos algo más de margen para luego reordenar en Go.
	fetchLimit := limit * 5
	if fetchLimit < 20 {
		fetchLimit = 20
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

		score := computeSearchScore(content, terms)
		if score <= 0 {
			continue
		}

		ranked = append(ranked, rankedSearchResult{
			result: querymodels.SearchDocumentResult{
				DocumentID:   documentID,
				OriginalName: originalName,
				CaseFileID:   caseFileID,
				Snippet:      buildSnippet(content, rawQuery, terms, 180),
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

func computeSearchScore(content string, terms []string) int {
	normalized := strings.ToLower(normalizeWhitespace(content))
	if normalized == "" {
		return 0
	}

	score := 0
	for _, term := range terms {
		term = strings.ToLower(strings.TrimSpace(term))
		if term == "" {
			continue
		}

		count := strings.Count(normalized, term)
		if count > 0 {
			score += count * len(term)
		}
	}

	return score
}

func buildSnippet(content string, rawQuery string, terms []string, maxLen int) string {
	normalized := normalizeWhitespace(content)
	if normalized == "" {
		return ""
	}

	bestTerm := bestMatchingTerm(normalized, terms)
	if bestTerm == "" {
		bestTerm = strings.TrimSpace(rawQuery)
	}

	if bestTerm == "" {
		return truncateSnippet(normalized, maxLen)
	}

	matchStart, matchEnd, ok := findCaseInsensitiveMatch(normalized, bestTerm)
	if !ok {
		return truncateSnippet(normalized, maxLen)
	}

	start := matchStart - (maxLen / 3)
	if start < 0 {
		start = 0
	}

	end := start + maxLen
	if end > len(normalized) {
		end = len(normalized)
		start = end - maxLen
		if start < 0 {
			start = 0
		}
	}

	snippet := normalized[start:end]
	relativeStart := matchStart - start
	relativeEnd := matchEnd - start

	if relativeStart < 0 || relativeEnd > len(snippet) || relativeStart >= relativeEnd {
		return decorateSnippetBounds(snippet, start, end, len(normalized))
	}

	highlighted := snippet[:relativeStart] + "[" + snippet[relativeStart:relativeEnd] + "]" + snippet[relativeEnd:]
	return decorateSnippetBounds(highlighted, start, end, len(normalized))
}

func bestMatchingTerm(content string, terms []string) string {
	lowerContent := strings.ToLower(content)

	bestTerm := ""
	bestCount := 0

	for _, term := range terms {
		term = strings.TrimSpace(term)
		if term == "" {
			continue
		}

		count := strings.Count(lowerContent, strings.ToLower(term))
		if count > bestCount {
			bestCount = count
			bestTerm = term
		}
	}

	return bestTerm
}

func splitSearchTerms(query string) []string {
	fields := strings.Fields(strings.TrimSpace(query))
	terms := make([]string, 0, len(fields))

	seen := make(map[string]struct{}, len(fields))
	for _, field := range fields {
		field = strings.TrimSpace(field)
		if field == "" {
			continue
		}

		lower := strings.ToLower(field)
		if _, exists := seen[lower]; exists {
			continue
		}
		seen[lower] = struct{}{}
		terms = append(terms, field)
	}

	return terms
}

func truncateSnippet(content string, maxLen int) string {
	if len(content) <= maxLen {
		return content
	}
	return content[:maxLen] + "..."
}

func decorateSnippetBounds(snippet string, start, end, totalLen int) string {
	if start > 0 {
		snippet = "..." + snippet
	}
	if end < totalLen {
		snippet = snippet + "..."
	}
	return snippet
}

func findCaseInsensitiveMatch(content string, query string) (int, int, bool) {
	lowerContent := strings.ToLower(content)
	lowerQuery := strings.ToLower(query)

	start := strings.Index(lowerContent, lowerQuery)
	if start == -1 {
		return 0, 0, false
	}

	return start, start + len(query), true
}

func normalizeWhitespace(value string) string {
	return strings.Join(strings.Fields(value), " ")
}
