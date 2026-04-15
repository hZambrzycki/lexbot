package documentapp

import (
	"context"
	"errors"
	"os"
	"strings"

	"lexbox/internal/application/ports"
	"lexbox/internal/domain/document"
	"lexbox/internal/domain/shared"
)

type ReindexAllDocumentsInput struct {
	CaseFileID string
}

type ReindexAllDocumentsResult struct {
	Scanned int

	Reindexed int
	Skipped   int
	Errors    int

	SkippedMissingFile    int
	SkippedAlreadyIndexed int
	SkippedUnsupported    int
	SkippedEmptyContent   int
}

type ReindexAllDocuments struct {
	Documents        ports.DocumentRepository
	DocumentContents ports.DocumentContentRepository
	ReindexDocument  ReindexDocument
}

func (uc ReindexAllDocuments) Execute(ctx context.Context, input ReindexAllDocumentsInput) (ReindexAllDocumentsResult, error) {
	var docs []document.Document

	if strings.TrimSpace(input.CaseFileID) != "" {
		caseFileID := shared.NewID(strings.TrimSpace(input.CaseFileID))
		if caseFileID == "" {
			return ReindexAllDocumentsResult{}, shared.ErrInvalidAssociation
		}

		caseDocs, err := uc.Documents.ListByCaseFileID(ctx, caseFileID)
		if err != nil {
			return ReindexAllDocumentsResult{}, err
		}
		docs = caseDocs
	} else {
		allDocs, err := uc.Documents.ListAll(ctx)
		if err != nil {
			return ReindexAllDocumentsResult{}, err
		}
		docs = allDocs
	}

	result := ReindexAllDocumentsResult{
		Scanned: len(docs),
	}

	for _, doc := range docs {
		if _, err := os.Stat(doc.StoragePath); err != nil {
			if os.IsNotExist(err) {
				result.Skipped++
				result.SkippedMissingFile++
				continue
			}
			return ReindexAllDocumentsResult{}, err
		}

		content, err := uc.DocumentContents.GetByDocumentID(ctx, doc.ID.String())
		if err == nil && strings.TrimSpace(content) != "" {
			result.Skipped++
			result.SkippedAlreadyIndexed++
			continue
		}
		if err != nil && !errors.Is(err, shared.ErrNotFound) {
			return ReindexAllDocumentsResult{}, err
		}

		_, err = uc.ReindexDocument.Execute(ctx, ReindexDocumentInput{
			DocumentID: doc.ID.String(),
		})
		if err != nil {
			switch {
			case errors.Is(err, ErrReindexUnsupported):
				result.Skipped++
				result.SkippedUnsupported++
			case errors.Is(err, ErrReindexEmptyContent):
				result.Skipped++
				result.SkippedEmptyContent++
			default:
				result.Errors++
			}
			continue
		}

		result.Reindexed++
	}

	return result, nil
}
