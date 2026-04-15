package documentapp

import (
	"context"
	"errors"
	"os"
	"strings"

	"lexbox/internal/application/ports"
	"lexbox/internal/domain/shared"
)

type VerifyDocumentInput struct {
	DocumentID string
}

type VerifyDocumentResult struct {
	DocumentID string

	ExistsInDB bool
	FileExists bool

	MimeType         string
	MimeNormalized   string
	IsMimeNormalized bool

	HasHash bool

	HasExtractedText bool
}

type VerifyDocument struct {
	Documents        ports.DocumentRepository
	DocumentContents ports.DocumentContentRepository
}

func (uc VerifyDocument) Execute(ctx context.Context, input VerifyDocumentInput) (VerifyDocumentResult, error) {
	documentID := shared.NewID(strings.TrimSpace(input.DocumentID))
	if documentID == "" {
		return VerifyDocumentResult{}, shared.ErrInvalidID
	}

	doc, err := uc.Documents.GetByID(ctx, documentID)
	if err != nil {
		return VerifyDocumentResult{}, err
	}

	result := VerifyDocumentResult{
		DocumentID:       doc.ID.String(),
		ExistsInDB:       true,
		MimeType:         doc.MimeType,
		MimeNormalized:   normalizeMimeType(doc.MimeType),
		IsMimeNormalized: normalizeMimeType(doc.MimeType) == doc.MimeType,
		HasHash:          strings.TrimSpace(doc.FileHash) != "",
		HasExtractedText: false,
	}

	if _, err := os.Stat(doc.StoragePath); err == nil {
		result.FileExists = true
	}

	content, err := uc.DocumentContents.GetByDocumentID(ctx, doc.ID.String())
	if err == nil && strings.TrimSpace(content) != "" {
		result.HasExtractedText = true
	} else if err != nil && !errors.Is(err, shared.ErrNotFound) {
		return VerifyDocumentResult{}, err
	}

	return result, nil
}
