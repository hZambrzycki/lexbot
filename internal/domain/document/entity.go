package document

import (
	"lexbox/internal/domain/shared"
)

type Document struct {
	ID           shared.ID
	CaseFileID   shared.ID
	OriginalName string
	StoragePath  string
	MimeType     string
	FileHash     string

	CreatedAt shared.Timestamp
	UpdatedAt shared.Timestamp

	ReviewStatus string
	ReviewedAt   string
	ReviewNote   string
}

func NewDocument(
	id shared.ID,
	caseFileID shared.ID,
	originalName string,
	storagePath string,
) (Document, error) {
	if caseFileID == "" {
		return Document{}, shared.ErrInvalidAssociation
	}

	if originalName == "" || storagePath == "" {
		return Document{}, shared.ErrEmptyField
	}

	now := shared.Now()

	return Document{
		ID:           id,
		CaseFileID:   caseFileID,
		OriginalName: originalName,
		StoragePath:  storagePath,
		CreatedAt:    now,
		UpdatedAt:    now,

		ReviewStatus: DocumentReviewStatusPending,
		ReviewedAt:   "",
		ReviewNote:   "",
	}, nil
}

func (d Document) WithUpdatedMetadata(mimeType string, fileHash string) Document {
	d.MimeType = mimeType
	d.FileHash = fileHash
	d.UpdatedAt = shared.Now()
	return d
}

const (
	DocumentReviewStatusPending  = "pending_review"
	DocumentReviewStatusReviewed = "reviewed"
	DocumentReviewStatusError    = "error"
)

func IsValidDocumentReviewStatus(status string) bool {
	switch status {
	case DocumentReviewStatusPending, DocumentReviewStatusReviewed, DocumentReviewStatusError:
		return true
	default:
		return false
	}
}
