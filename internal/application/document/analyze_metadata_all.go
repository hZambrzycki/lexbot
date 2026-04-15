package documentapp

import (
	"context"
	"errors"
	"strings"

	"lexbox/internal/application/ports"
	"lexbox/internal/domain/document"
	"lexbox/internal/domain/shared"
)

type AnalyzeAllDocumentMetadataInput struct {
	CaseFileID string
}

type AnalyzeAllDocumentMetadataResult struct {
	Scanned  int
	Analyzed int
	Skipped  int
	Errors   int

	SkippedNoExtractedText int
	SkippedInvalidCase     int
}

type AnalyzeAllDocumentMetadata struct {
	Documents        ports.DocumentRepository
	DocumentContents ports.DocumentContentRepository
	AnalyzeOne       AnalyzeDocumentMetadata
}

func (uc AnalyzeAllDocumentMetadata) Execute(ctx context.Context, input AnalyzeAllDocumentMetadataInput) (AnalyzeAllDocumentMetadataResult, error) {
	var docs []document.Document

	if strings.TrimSpace(input.CaseFileID) != "" {
		caseFileID := shared.NewID(strings.TrimSpace(input.CaseFileID))
		if caseFileID == "" {
			return AnalyzeAllDocumentMetadataResult{}, shared.ErrInvalidAssociation
		}

		caseDocs, err := uc.Documents.ListByCaseFileID(ctx, caseFileID)
		if err != nil {
			return AnalyzeAllDocumentMetadataResult{}, err
		}
		docs = caseDocs
	} else {
		allDocs, err := uc.Documents.ListAll(ctx)
		if err != nil {
			return AnalyzeAllDocumentMetadataResult{}, err
		}
		docs = allDocs
	}

	result := AnalyzeAllDocumentMetadataResult{
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

		_, err = uc.AnalyzeOne.Execute(ctx, AnalyzeDocumentMetadataInput{
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
