package documentapp

import (
	"context"
	"strings"

	"lexbox/internal/application/ports"
	"lexbox/internal/domain/shared"
)

type ExtractDocumentTextInput struct {
	DocumentID string
}

type ExtractDocumentText struct {
	Documents        ports.DocumentRepository
	DocumentContents ports.DocumentContentRepository
	Extractor        ports.TextExtractor
}

func (uc ExtractDocumentText) Execute(ctx context.Context, in ExtractDocumentTextInput) (string, error) {
	documentID := shared.NewID(strings.TrimSpace(in.DocumentID))
	if documentID == "" {
		return "", shared.ErrInvalidID
	}

	doc, err := uc.Documents.GetByID(ctx, documentID)
	if err != nil {
		return "", err
	}

	content, err := uc.Extractor.ExtractText(ctx, doc.StoragePath)
	if err != nil {
		return "", err
	}

	if err := uc.DocumentContents.Save(ctx, documentID.String(), content); err != nil {
		return "", err
	}

	return content, nil
}
