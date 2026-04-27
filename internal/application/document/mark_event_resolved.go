package documentapp

import (
	"context"
	"time"

	"lexbox/internal/application/ports"
	"lexbox/internal/domain/document"
	"lexbox/internal/domain/shared"
)

type MarkEventResolved struct {
	Events ports.DocumentEventRepository
}

type MarkEventResolvedInput struct {
	EventID        string
	ResolutionNote string
}

type MarkEventResolvedResult struct {
	EventID        string
	ReviewStatus   string
	ReviewedAt     string
	ResolvedAt     string
	ResolutionNote string
}

func (uc MarkEventResolved) Execute(ctx context.Context, input MarkEventResolvedInput) (MarkEventResolvedResult, error) {
	if input.EventID == "" {
		return MarkEventResolvedResult{}, shared.ErrInvalidID
	}

	eventID := shared.NewID(input.EventID)

	existing, err := uc.Events.GetByID(ctx, eventID)
	if err != nil {
		return MarkEventResolvedResult{}, err
	}

	now := time.Now().UTC().Format(time.RFC3339)

	reviewedAt := existing.ReviewedAt
	if reviewedAt == "" {
		reviewedAt = now
	}

	err = uc.Events.UpdateReviewState(
		ctx,
		eventID,
		document.ReviewStatusResolved,
		reviewedAt,
		now,
		input.ResolutionNote,
	)
	if err != nil {
		return MarkEventResolvedResult{}, err
	}

	return MarkEventResolvedResult{
		EventID:        existing.ID.String(),
		ReviewStatus:   document.ReviewStatusResolved,
		ReviewedAt:     reviewedAt,
		ResolvedAt:     now,
		ResolutionNote: input.ResolutionNote,
	}, nil
}
