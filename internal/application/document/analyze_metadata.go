package documentapp

import (
	"context"
	"strings"
	"time"

	"lexbox/internal/application/ports"
	"lexbox/internal/domain/document"
	"lexbox/internal/domain/shared"
)

type AnalyzeDocumentMetadataInput struct {
	DocumentID string
}

type AnalyzeDocumentMetadataResult struct {
	DocumentID   string
	DocumentType string
	LegalArea    string
	AnalyzedAt   string
}

type AnalyzeDocumentMetadata struct {
	Documents        ports.DocumentRepository
	DocumentContents ports.DocumentContentRepository
	Metadata         ports.DocumentMetadataRepository
}

func (uc AnalyzeDocumentMetadata) Execute(ctx context.Context, input AnalyzeDocumentMetadataInput) (AnalyzeDocumentMetadataResult, error) {
	documentID := shared.NewID(strings.TrimSpace(input.DocumentID))
	if documentID == "" {
		return AnalyzeDocumentMetadataResult{}, shared.ErrInvalidID
	}

	_, err := uc.Documents.GetByID(ctx, documentID)
	if err != nil {
		return AnalyzeDocumentMetadataResult{}, err
	}

	content, err := uc.DocumentContents.GetByDocumentID(ctx, documentID.String())
	if err != nil {
		return AnalyzeDocumentMetadataResult{}, err
	}

	classification := classifyDocumentMetadata(content)

	metadata, err := document.NewMetadata(
		documentID,
		classification.DocumentType,
		classification.LegalArea,
		shared.Now().Time().Format(time.RFC3339),
	)
	if err != nil {
		return AnalyzeDocumentMetadataResult{}, err
	}

	if err := uc.Metadata.Save(ctx, metadata); err != nil {
		return AnalyzeDocumentMetadataResult{}, err
	}

	return AnalyzeDocumentMetadataResult{
		DocumentID:   metadata.DocumentID.String(),
		DocumentType: metadata.DocumentType,
		LegalArea:    metadata.LegalArea,
		AnalyzedAt:   metadata.AnalyzedAt,
	}, nil
}
