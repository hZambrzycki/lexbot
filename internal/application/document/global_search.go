package documentapp

import (
	"context"
	"strings"

	"lexbox/internal/application/ports"
	"lexbox/internal/application/querymodels"
)

type GlobalSearchInput struct {
	Query string
	Limit int
}

type GlobalSearch struct {
	SearchDocuments SearchDocuments
	CaseFiles       ports.CaseFileSearchRepository
	Events          ports.EventSearchRepository
	Notes           ports.NoteSearchRepository
}

func (uc GlobalSearch) Execute(ctx context.Context, in GlobalSearchInput) ([]querymodels.GlobalSearchResult, error) {
	query := strings.TrimSpace(in.Query)
	if query == "" {
		return []querymodels.GlobalSearchResult{}, nil
	}

	limit := in.Limit
	if limit <= 0 {
		limit = 8
	}
	if limit > 25 {
		limit = 25
	}

	results := make([]querymodels.GlobalSearchResult, 0)

	if uc.CaseFiles != nil {
		items, err := uc.CaseFiles.SearchCaseFiles(ctx, query, limit)
		if err != nil {
			return nil, err
		}
		results = append(results, items...)
	}

	documentResults, err := uc.SearchDocuments.Execute(ctx, SearchDocumentsInput{
		Query: query,
		Limit: limit,
	})
	if err != nil {
		return nil, err
	}

	for _, item := range documentResults {
		results = append(results, querymodels.GlobalSearchResult{
			Type:     "document",
			ID:       item.DocumentID,
			Title:    item.OriginalName,
			Subtitle: "Documento",
			Href:     "/case-files/" + item.CaseFileID + "/documents/" + item.DocumentID,
			Snippet:  item.Snippet,
			Score:    item.Score,
		})
	}

	if uc.Events != nil {
		items, err := uc.Events.SearchEvents(ctx, query, limit)
		if err != nil {
			return nil, err
		}
		results = append(results, items...)
	}

	if uc.Notes != nil {
		items, err := uc.Notes.SearchNotes(ctx, query, limit)
		if err != nil {
			return nil, err
		}
		results = append(results, items...)
	}

	return results, nil
}
