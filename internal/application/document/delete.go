package documentapp

import (
	"context"
	"strings"

	"lexbox/internal/application/ports"
	"lexbox/internal/domain/shared"
)

type DeleteDocumentInput struct {
	DocumentID string
}

type DeleteDocumentResult struct {
	DocumentID  string
	CaseFileID  string
	StoragePath string
	Deleted     bool
}

type DeleteDocument struct {
	Documents ports.DocumentRepository
	Storage   ports.FileStorage
}

func (uc DeleteDocument) Execute(ctx context.Context, in DeleteDocumentInput) (DeleteDocumentResult, error) {
	documentID := shared.NewID(strings.TrimSpace(in.DocumentID))
	if documentID == "" {
		return DeleteDocumentResult{}, shared.ErrInvalidID
	}

	doc, err := uc.Documents.GetByID(ctx, documentID)
	if err != nil {
		return DeleteDocumentResult{}, err
	}

	if err := uc.Documents.Delete(ctx, documentID); err != nil {
		return DeleteDocumentResult{}, err
	}

	if uc.Storage != nil && strings.TrimSpace(doc.StoragePath) != "" {
		if err := uc.Storage.DeleteStoredFile(ctx, doc.StoragePath); err != nil {
			return DeleteDocumentResult{}, err
		}
	}

	return DeleteDocumentResult{
		DocumentID:  doc.ID.String(),
		CaseFileID:  doc.CaseFileID.String(),
		StoragePath: doc.StoragePath,
		Deleted:     true,
	}, nil
}
