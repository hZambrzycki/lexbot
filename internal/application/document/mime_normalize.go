package documentapp

import (
	"context"

	"lexbox/internal/application/ports"
)

type MimeNormalizeChangedDocument struct {
	DocumentID  string
	OldMimeType string
	NewMimeType string
}

type MimeNormalizeResult struct {
	Scanned int
	Updated int
	Changed []MimeNormalizeChangedDocument
}

type MimeNormalize struct {
	Documents ports.DocumentRepository
}

func (uc MimeNormalize) Execute(ctx context.Context) (MimeNormalizeResult, error) {
	documents, err := uc.Documents.ListAll(ctx)
	if err != nil {
		return MimeNormalizeResult{}, err
	}

	result := MimeNormalizeResult{
		Scanned: len(documents),
		Updated: 0,
		Changed: make([]MimeNormalizeChangedDocument, 0),
	}

	for _, doc := range documents {
		normalized := normalizeMimeType(doc.MimeType)
		if normalized == doc.MimeType {
			continue
		}

		updatedDoc := doc.WithUpdatedMetadata(normalized, doc.FileHash)
		if err := uc.Documents.Update(ctx, updatedDoc); err != nil {
			return MimeNormalizeResult{}, err
		}

		result.Updated++
		result.Changed = append(result.Changed, MimeNormalizeChangedDocument{
			DocumentID:  doc.ID.String(),
			OldMimeType: doc.MimeType,
			NewMimeType: normalized,
		})
	}

	return result, nil
}
