package documentapp

import (
	"context"
	"strings"
	"time"

	"lexbox/internal/application/ports"
	"lexbox/internal/domain/document"
	"lexbox/internal/domain/shared"
)

type AnalyzeDocumentEventsInput struct {
	DocumentID string
}

type AnalyzeDocumentEventsResult struct {
	DocumentID string
	Detected   int
}

type AnalyzeDocumentEvents struct {
	Documents        ports.DocumentRepository
	DocumentContents ports.DocumentContentRepository
	Events           ports.DocumentEventRepository
	IDs              ports.IDGenerator
}

func (uc AnalyzeDocumentEvents) Execute(ctx context.Context, input AnalyzeDocumentEventsInput) (AnalyzeDocumentEventsResult, error) {
	documentID := shared.NewID(strings.TrimSpace(input.DocumentID))
	if documentID == "" {
		return AnalyzeDocumentEventsResult{}, shared.ErrInvalidID
	}

	_, err := uc.Documents.GetByID(ctx, documentID)
	if err != nil {
		return AnalyzeDocumentEventsResult{}, err
	}

	content, err := uc.DocumentContents.GetByDocumentID(ctx, documentID.String())
	if err != nil {
		return AnalyzeDocumentEventsResult{}, err
	}

	candidates := extractDocumentEvents(content)

	events := make([]document.Event, 0, len(candidates))
	createdAt := time.Now().Format(time.RFC3339)

	for _, candidate := range candidates {
		event, err := document.NewEvent(
			uc.IDs.NewID(),
			documentID,
			candidate.EventType,
			candidate.EventDate,
			candidate.SourceText,
			createdAt,
			candidate.AnchorDate,
			candidate.DateKind,
			candidate.AnchorSource,
			candidate.RelativeDays,
			candidate.IsBusinessDays,
			candidate.TriggerText,
		)
		if err != nil {
			return AnalyzeDocumentEventsResult{}, err
		}
		events = append(events, event)
	}

	if err := uc.Events.ReplaceByDocumentID(ctx, documentID, events); err != nil {
		return AnalyzeDocumentEventsResult{}, err
	}

	return AnalyzeDocumentEventsResult{
		DocumentID: documentID.String(),
		Detected:   len(events),
	}, nil
}
