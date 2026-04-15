package documentapp

import (
	"context"
	"errors"
	"strings"

	"lexbox/internal/application/ports"
	"lexbox/internal/domain/document"
	"lexbox/internal/domain/shared"
)

type ListDocumentsByCaseFileInput struct {
	CaseFileID string
}

type ListDocumentsByCaseFileItem struct {
	Document     document.Document
	DocumentType string
	LegalArea    string
	HasMetadata  bool
}

type ListDocumentsByCaseFile struct {
	Documents ports.DocumentRepository
	Metadata  ports.DocumentMetadataRepository
}

func (uc ListDocumentsByCaseFile) Execute(ctx context.Context, in ListDocumentsByCaseFileInput) ([]ListDocumentsByCaseFileItem, error) {
	caseFileID := shared.NewID(strings.TrimSpace(in.CaseFileID))
	if caseFileID == "" {
		return nil, shared.ErrInvalidID
	}

	docs, err := uc.Documents.ListByCaseFileID(ctx, caseFileID)
	if err != nil {
		return nil, err
	}

	items := make([]ListDocumentsByCaseFileItem, 0, len(docs))

	for _, doc := range docs {
		item := ListDocumentsByCaseFileItem{
			Document:     doc,
			DocumentType: "unknown",
			LegalArea:    "unknown",
			HasMetadata:  false,
		}

		if uc.Metadata != nil {
			metadata, err := uc.Metadata.GetByDocumentID(ctx, doc.ID)
			if err != nil {
				if !errors.Is(err, shared.ErrNotFound) {
					return nil, err
				}
			} else {
				item.DocumentType = metadata.DocumentType
				item.LegalArea = metadata.LegalArea
				item.HasMetadata = true
			}
		}

		items = append(items, item)
	}

	return items, nil
}
