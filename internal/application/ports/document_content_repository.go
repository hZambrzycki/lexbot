package ports

import (
	"context"

	"lexbox/internal/application/querymodels"
)

type DocumentContentRepository interface {
	Save(ctx context.Context, documentID string, content string) error
	GetByDocumentID(ctx context.Context, documentID string) (string, error)
	SearchByText(ctx context.Context, query string, limit int) ([]querymodels.SearchDocumentResult, error)
	SearchByTextByCaseFile(ctx context.Context, caseFileID string, query string, limit int) ([]querymodels.SearchDocumentResult, error)
}
