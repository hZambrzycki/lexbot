package documentapp

import (
	"context"
	"errors"
	"strings"

	"lexbox/internal/application/ports"
	"lexbox/internal/domain/document"
	"lexbox/internal/domain/shared"
)

type AnalyzeAllDocumentEventsInput struct {
	CaseFileID string
}

type AnalyzeAllDocumentEventsResult struct {
	Scanned  int
	Analyzed int
	Skipped  int
	Errors   int

	SkippedNoExtractedText int
}

type AnalyzeAllDocumentEvents struct {
	Documents        ports.DocumentRepository
	DocumentContents ports.DocumentContentRepository
	AnalyzeOne       AnalyzeDocumentEvents
}

func (uc AnalyzeAllDocumentEvents) Execute(ctx context.Context, input AnalyzeAllDocumentEventsInput) (AnalyzeAllDocumentEventsResult, error) {
	var docs []document.Document

	if strings.TrimSpace(input.CaseFileID) != "" {
		caseFileID := shared.NewID(strings.TrimSpace(input.CaseFileID))
		if caseFileID == "" {
			return AnalyzeAllDocumentEventsResult{}, shared.ErrInvalidAssociation
		}

		caseDocs, err := uc.Documents.ListByCaseFileID(ctx, caseFileID)
		if err != nil {
			return AnalyzeAllDocumentEventsResult{}, err
		}
		docs = caseDocs
	} else {
		allDocs, err := uc.Documents.ListAll(ctx)
		if err != nil {
			return AnalyzeAllDocumentEventsResult{}, err
		}
		docs = allDocs
	}

	result := AnalyzeAllDocumentEventsResult{
		Scanned: len(docs),
	}

	for _, doc := range docs {
		content, err := uc.DocumentContents.GetByDocumentID(ctx, doc.ID.String())
		if err != nil {
			if errors.Is(err, shared.ErrNotFound) {
				result.Skipped++
				result.SkippedNoExtractedText++
				continue
			}
			result.Errors++
			continue
		}

		if strings.TrimSpace(content) == "" {
			result.Skipped++
			result.SkippedNoExtractedText++
			continue
		}

		_, err = uc.AnalyzeOne.Execute(ctx, AnalyzeDocumentEventsInput{
			DocumentID: doc.ID.String(),
		})
		if err != nil {
			result.Errors++
			continue
		}

		result.Analyzed++
	}

	return result, nil
}
