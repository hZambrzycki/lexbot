package documentapp

import (
	"context"
	"errors"
	"os"
	"strings"

	"lexbox/internal/application/ports"
	"lexbox/internal/domain/document"
	"lexbox/internal/domain/shared"
)

type VerifyAllDocuments struct {
	Documents        ports.DocumentRepository
	DocumentContents ports.DocumentContentRepository
}

type VerifyAllDocumentsInput struct {
	CaseFileID string
}

type VerifyAllDocumentsResult struct {
	Total int
	OK    int

	MissingFile int
	MissingHash int
	MissingText int
	InvalidMime int

	MissingFileItems []VerifyAllDocumentsMissingFileItem
	MissingHashItems []VerifyAllDocumentsMissingHashItem
	MissingTextItems []VerifyAllDocumentsMissingTextItem
	InvalidMimeItems []VerifyAllDocumentsInvalidMimeItem
}

type VerifyAllDocumentsMissingFileItem struct {
	DocumentID   string
	OriginalName string
	StoragePath  string
}

type VerifyAllDocumentsMissingHashItem struct {
	DocumentID   string
	OriginalName string
	StoragePath  string
}

type VerifyAllDocumentsMissingTextItem struct {
	DocumentID   string
	OriginalName string
	StoragePath  string
}

type VerifyAllDocumentsInvalidMimeItem struct {
	DocumentID     string
	OriginalName   string
	StoragePath    string
	MimeType       string
	NormalizedMime string
}

func (uc VerifyAllDocuments) Execute(ctx context.Context, input VerifyAllDocumentsInput) (VerifyAllDocumentsResult, error) {
	var docs []document.Document

	if strings.TrimSpace(input.CaseFileID) != "" {
		caseFileID := shared.NewID(strings.TrimSpace(input.CaseFileID))
		if caseFileID == "" {
			return VerifyAllDocumentsResult{}, shared.ErrInvalidAssociation
		}

		caseDocs, err := uc.Documents.ListByCaseFileID(ctx, caseFileID)
		if err != nil {
			return VerifyAllDocumentsResult{}, err
		}
		docs = caseDocs
	} else {
		allDocs, err := uc.Documents.ListAll(ctx)
		if err != nil {
			return VerifyAllDocumentsResult{}, err
		}
		docs = allDocs
	}

	result := VerifyAllDocumentsResult{
		Total: len(docs),
	}

	for _, doc := range docs {
		isOK := true

		if _, err := os.Stat(doc.StoragePath); err != nil {
			result.MissingFile++
			result.MissingFileItems = append(result.MissingFileItems, VerifyAllDocumentsMissingFileItem{
				DocumentID:   doc.ID.String(),
				OriginalName: doc.OriginalName,
				StoragePath:  doc.StoragePath,
			})
			isOK = false
		}

		if strings.TrimSpace(doc.FileHash) == "" {
			result.MissingHash++
			result.MissingHashItems = append(result.MissingHashItems, VerifyAllDocumentsMissingHashItem{
				DocumentID:   doc.ID.String(),
				OriginalName: doc.OriginalName,
				StoragePath:  doc.StoragePath,
			})
			isOK = false
		}

		content, err := uc.DocumentContents.GetByDocumentID(ctx, doc.ID.String())
		if err != nil {
			if !errors.Is(err, shared.ErrNotFound) {
				return VerifyAllDocumentsResult{}, err
			}
			result.MissingText++
			result.MissingTextItems = append(result.MissingTextItems, VerifyAllDocumentsMissingTextItem{
				DocumentID:   doc.ID.String(),
				OriginalName: doc.OriginalName,
				StoragePath:  doc.StoragePath,
			})
			isOK = false
		} else if strings.TrimSpace(content) == "" {
			result.MissingText++
			result.MissingTextItems = append(result.MissingTextItems, VerifyAllDocumentsMissingTextItem{
				DocumentID:   doc.ID.String(),
				OriginalName: doc.OriginalName,
				StoragePath:  doc.StoragePath,
			})
			isOK = false
		}

		normalized := normalizeMimeType(doc.MimeType)
		if normalized != doc.MimeType {
			result.InvalidMime++
			result.InvalidMimeItems = append(result.InvalidMimeItems, VerifyAllDocumentsInvalidMimeItem{
				DocumentID:     doc.ID.String(),
				OriginalName:   doc.OriginalName,
				StoragePath:    doc.StoragePath,
				MimeType:       doc.MimeType,
				NormalizedMime: normalized,
			})
			isOK = false
		}

		if isOK {
			result.OK++
		}
	}

	return result, nil
}
