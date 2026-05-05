package documentapp

import (
	"context"
	"errors"
	"strings"
	"time"

	"lexbox/internal/application/ports"
	domaindoc "lexbox/internal/domain/document"
	"lexbox/internal/domain/shared"
)

var (
	ErrInvalidDocumentReviewStatus = errors.New("invalid document review status")
	ErrDocumentIDRequired          = errors.New("document_id is required")
	ErrReviewNoteRequired          = errors.New("review_note is required when marking document as error")
)

type ReviewDocument struct {
	Documents ports.DocumentRepository
}

type ReviewDocumentInput struct {
	DocumentID   string
	ReviewStatus string
	ReviewNote   string
}

type ReviewDocumentOutput struct {
	DocumentID   string `json:"document_id"`
	ReviewStatus string `json:"review_status"`
	ReviewedAt   string `json:"reviewed_at"`
	ReviewNote   string `json:"review_note"`
}

func (uc ReviewDocument) Execute(ctx context.Context, input ReviewDocumentInput) (ReviewDocumentOutput, error) {
	documentID := shared.ID(strings.TrimSpace(input.DocumentID))
	reviewStatus := strings.TrimSpace(input.ReviewStatus)
	reviewNote := strings.TrimSpace(input.ReviewNote)

	if documentID == "" {
		return ReviewDocumentOutput{}, ErrDocumentIDRequired
	}

	if !domaindoc.IsValidDocumentReviewStatus(reviewStatus) {
		return ReviewDocumentOutput{}, ErrInvalidDocumentReviewStatus
	}

	if reviewStatus == domaindoc.DocumentReviewStatusError && reviewNote == "" {
		return ReviewDocumentOutput{}, ErrReviewNoteRequired
	}

	if _, err := uc.Documents.GetByID(ctx, documentID); err != nil {
		return ReviewDocumentOutput{}, err
	}

	reviewedAt := ""

	switch reviewStatus {
	case domaindoc.DocumentReviewStatusPending:
		reviewNote = ""
		reviewedAt = ""

	case domaindoc.DocumentReviewStatusReviewed, domaindoc.DocumentReviewStatusError:
		reviewedAt = time.Now().UTC().Format(time.RFC3339)
	}

	if err := uc.Documents.UpdateReviewState(
		ctx,
		documentID,
		reviewStatus,
		reviewedAt,
		reviewNote,
	); err != nil {
		return ReviewDocumentOutput{}, err
	}

	return ReviewDocumentOutput{
		DocumentID:   string(documentID),
		ReviewStatus: reviewStatus,
		ReviewedAt:   reviewedAt,
		ReviewNote:   reviewNote,
	}, nil
}
