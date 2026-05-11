package documentapp

import (
	"context"
	"strings"

	"lexbox/internal/application/ports"
	"lexbox/internal/application/querymodels"
	"lexbox/internal/domain/shared"
)

type SearchDocumentsInput struct {
	Query      string
	CaseFileID string
	Limit      int
}

type SearchDocuments struct {
	SearchIndex ports.DocumentSearchIndexRepository
}

func (uc SearchDocuments) Execute(ctx context.Context, in SearchDocumentsInput) ([]querymodels.SearchDocumentResult, error) {
	query := strings.TrimSpace(in.Query)
	if query == "" {
		return []querymodels.SearchDocumentResult{}, nil
	}

	limit := in.Limit
	if limit <= 0 {
		limit = 20
	}

	caseFileID := strings.TrimSpace(in.CaseFileID)
	if caseFileID != "" {
		parsedCaseFileID := shared.NewID(caseFileID)
		if parsedCaseFileID == "" {
			return nil, shared.ErrInvalidAssociation
		}

		caseFileID = parsedCaseFileID.String()
	}

	if uc.SearchIndex == nil {
		return []querymodels.SearchDocumentResult{}, nil
	}

	parsedQuery := ParseSearchDocumentsQuery(query)

	return uc.SearchIndex.Search(
		ctx,
		parsedQuery.Terms,
		caseFileID,
		parsedQuery.Filters,
		limit,
	)
}
