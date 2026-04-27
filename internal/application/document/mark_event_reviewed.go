package documentapp

import (
	"context"
	"time"

	"lexbox/internal/application/ports"
	"lexbox/internal/domain/document"
	"lexbox/internal/domain/shared"
)

type MarkEventReviewed struct {
	Events ports.DocumentEventRepository
}

type MarkEventReviewedInput struct {
	EventID string
}

type MarkEventReviewedResult struct {
	EventID        string
	ReviewStatus   string
	ReviewedAt     string
	ResolvedAt     string
	ResolutionNote string
}

func (uc MarkEventReviewed) Execute(ctx context.Context, input MarkEventReviewedInput) (MarkEventReviewedResult, error) {
	if input.EventID == "" {
		return MarkEventReviewedResult{}, shared.ErrInvalidID
	}

	eventID := shared.NewID(input.EventID)

	existing, err := uc.Events.GetByID(ctx, eventID)
	if err != nil {
		return MarkEventReviewedResult{}, err
	}

	reviewedAt := time.Now().UTC().Format(time.RFC3339)

	err = uc.Events.UpdateReviewState(
		ctx,
		eventID,
		document.ReviewStatusReviewed,
		reviewedAt,
		"",
		"",
	)
	if err != nil {
		return MarkEventReviewedResult{}, err
	}

	return MarkEventReviewedResult{
		EventID:        existing.ID.String(),
		ReviewStatus:   document.ReviewStatusReviewed,
		ReviewedAt:     reviewedAt,
		ResolvedAt:     "",
		ResolutionNote: "",
	}, nil
}
