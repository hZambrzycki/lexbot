package documentapp

import (
	"context"
	"errors"
	"strings"

	"lexbox/internal/application/ports"
	"lexbox/internal/domain/document"
	"lexbox/internal/domain/shared"
)

type BackfillDocumentSearchIndexInput struct {
	CaseFileID string
}

type BackfillDocumentSearchIndexResult struct {
	Scanned int
	Indexed int
	Skipped int
	Errors  int

	SkippedNoExtractedText int
	SkippedNoMetadata      int
}

type BackfillDocumentSearchIndex struct {
	Documents        ports.DocumentRepository
	DocumentContents ports.DocumentContentRepository
	Metadata         ports.DocumentMetadataRepository
	SearchIndex      ports.DocumentSearchIndexRepository
}

func (uc BackfillDocumentSearchIndex) Execute(
	ctx context.Context,
	input BackfillDocumentSearchIndexInput,
) (BackfillDocumentSearchIndexResult, error) {
	var docs []document.Document

	if strings.TrimSpace(input.CaseFileID) != "" {
		caseFileID := shared.NewID(strings.TrimSpace(input.CaseFileID))
		if caseFileID == "" {
			return BackfillDocumentSearchIndexResult{}, shared.ErrInvalidAssociation
		}

		caseDocs, err := uc.Documents.ListByCaseFileID(ctx, caseFileID)
		if err != nil {
			return BackfillDocumentSearchIndexResult{}, err
		}

		docs = caseDocs
	} else {
		allDocs, err := uc.Documents.ListAll(ctx)
		if err != nil {
			return BackfillDocumentSearchIndexResult{}, err
		}

		docs = allDocs
	}

	result := BackfillDocumentSearchIndexResult{
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

		documentType := "unknown"
		legalArea := "unknown"

		metadata, err := uc.Metadata.GetByDocumentID(ctx, doc.ID)
		if err != nil {
			if errors.Is(err, shared.ErrNotFound) {
				result.SkippedNoMetadata++
			} else {
				result.Errors++
				continue
			}
		} else {
			documentType = metadata.DocumentType
			legalArea = metadata.LegalArea
		}

		if err := uc.SearchIndex.UpsertDocument(
			ctx,
			doc.ID.String(),
			doc.CaseFileID.String(),
			doc.OriginalName,
			content,
			documentType,
			legalArea,
		); err != nil {
			result.Errors++
			continue
		}

		result.Indexed++
	}

	return result, nil
}
