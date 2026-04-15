package documentapp

import (
	"context"
	"strings"

	"lexbox/internal/application/ports"
	"lexbox/internal/domain/document"
	"lexbox/internal/domain/shared"
)

type AttachDocumentInput struct {
	CaseFileID   string
	OriginalName string
	StoragePath  string
	MimeType     string
	FileHash     string
}

type AttachDocument struct {
	Documents ports.DocumentRepository
	CaseFiles ports.CaseFileRepository
	IDs       ports.IDGenerator
}

func (uc AttachDocument) Execute(ctx context.Context, in AttachDocumentInput) (document.Document, error) {
	caseFileID := shared.NewID(strings.TrimSpace(in.CaseFileID))
	if caseFileID == "" {
		return document.Document{}, shared.ErrInvalidAssociation
	}

	_, err := uc.CaseFiles.GetByID(ctx, caseFileID)
	if err != nil {
		return document.Document{}, err
	}

	doc, err := document.NewDocument(
		uc.IDs.NewID(),
		caseFileID,
		strings.TrimSpace(in.OriginalName),
		strings.TrimSpace(in.StoragePath),
	)
	if err != nil {
		return document.Document{}, err
	}

	doc.MimeType = normalizeMimeType(in.MimeType)
	doc.FileHash = strings.TrimSpace(in.FileHash)

	if err := uc.Documents.Save(ctx, doc); err != nil {
		return document.Document{}, err
	}

	return doc, nil
}
