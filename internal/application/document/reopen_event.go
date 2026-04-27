package documentapp

import (
	"context"

	"lexbox/internal/application/ports"
	"lexbox/internal/domain/document"
	"lexbox/internal/domain/shared"
)

type ReopenEvent struct {
	Events ports.DocumentEventRepository
}

type ReopenEventInput struct {
	EventID string
}

type ReopenEventResult struct {
	EventID        string
	ReviewStatus   string
	ReviewedAt     string
	ResolvedAt     string
	ResolutionNote string
}

func (uc ReopenEvent) Execute(ctx context.Context, input ReopenEventInput) (ReopenEventResult, error) {
	if input.EventID == "" {
		return ReopenEventResult{}, shared.ErrInvalidID
	}

	eventID := shared.NewID(input.EventID)

	existing, err := uc.Events.GetByID(ctx, eventID)
	if err != nil {
		return ReopenEventResult{}, err
	}

	err = uc.Events.UpdateReviewState(
		ctx,
		eventID,
		document.ReviewStatusPending,
		"",
		"",
		"",
	)
	if err != nil {
		return ReopenEventResult{}, err
	}

	return ReopenEventResult{
		EventID:        existing.ID.String(),
		ReviewStatus:   document.ReviewStatusPending,
		ReviewedAt:     "",
		ResolvedAt:     "",
		ResolutionNote: "",
	}, nil
}
