package documentapp

import (
	"context"
	"strings"

	"lexbox/internal/application/ports"
	"lexbox/internal/domain/document"
	"lexbox/internal/domain/shared"
)

type ListDocumentEventsInput struct {
	DocumentID string
}

type ListDocumentEvents struct {
	Events ports.DocumentEventRepository
}

func (uc ListDocumentEvents) Execute(ctx context.Context, in ListDocumentEventsInput) ([]document.Event, error) {
	documentID := shared.NewID(strings.TrimSpace(in.DocumentID))
	if documentID == "" {
		return nil, shared.ErrInvalidID
	}

	return uc.Events.ListByDocumentID(ctx, documentID)
}
