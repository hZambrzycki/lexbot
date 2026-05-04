package documentapp

import (
	"context"
	"strings"
	"time"

	"lexbox/internal/application/ports"
	"lexbox/internal/domain/document"
	"lexbox/internal/domain/shared"
)

type UpdateDocumentReviewStateInput struct {
	DocumentID   string
	ReviewStatus string
	ReviewNote   string
}

type UpdateDocumentReviewStateResult struct {
	DocumentID   string `json:"document_id"`
	ReviewStatus string `json:"review_status"`
	ReviewedAt   string `json:"reviewed_at"`
	ReviewNote   string `json:"review_note"`
}

type UpdateDocumentReviewState struct {
	Documents ports.DocumentRepository
}

func (uc UpdateDocumentReviewState) Execute(ctx context.Context, in UpdateDocumentReviewStateInput) (UpdateDocumentReviewStateResult, error) {
	documentID := shared.NewID(strings.TrimSpace(in.DocumentID))
	if documentID == "" {
		return UpdateDocumentReviewStateResult{}, shared.ErrInvalidID
	}

	reviewStatus := strings.TrimSpace(in.ReviewStatus)
	if !document.IsValidDocumentReviewStatus(reviewStatus) {
		return UpdateDocumentReviewStateResult{}, shared.ErrInvalidAssociation
	}

	reviewNote := strings.TrimSpace(in.ReviewNote)

	reviewedAt := ""
	if reviewStatus == document.DocumentReviewStatusReviewed || reviewStatus == document.DocumentReviewStatusError {
		reviewedAt = time.Now().Format(time.RFC3339)
	}

	if reviewStatus == document.DocumentReviewStatusPending {
		reviewNote = ""
		reviewedAt = ""
	}

	if err := uc.Documents.UpdateReviewState(ctx, documentID, reviewStatus, reviewedAt, reviewNote); err != nil {
		return UpdateDocumentReviewStateResult{}, err
	}

	return UpdateDocumentReviewStateResult{
		DocumentID:   documentID.String(),
		ReviewStatus: reviewStatus,
		ReviewedAt:   reviewedAt,
		ReviewNote:   reviewNote,
	}, nil
}
