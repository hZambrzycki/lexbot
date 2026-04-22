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
	CaseFiles        ports.CaseFileRepository
	Events           ports.DocumentEventRepository
	IDs              ports.IDGenerator
}

func (uc AnalyzeDocumentEvents) Execute(ctx context.Context, input AnalyzeDocumentEventsInput) (AnalyzeDocumentEventsResult, error) {
	documentID := shared.NewID(strings.TrimSpace(input.DocumentID))
	if documentID == "" {
		return AnalyzeDocumentEventsResult{}, shared.ErrInvalidID
	}

	doc, err := uc.Documents.GetByID(ctx, documentID)
	if err != nil {
		return AnalyzeDocumentEventsResult{}, err
	}

	content, err := uc.DocumentContents.GetByDocumentID(ctx, documentID.String())
	if err != nil {
		return AnalyzeDocumentEventsResult{}, err
	}

	cfg := DefaultEventComputationConfig()

	if uc.CaseFiles != nil && doc.CaseFileID != "" {
		cf, err := uc.CaseFiles.GetByID(ctx, doc.CaseFileID)
		if err != nil {
			return AnalyzeDocumentEventsResult{}, err
		}
		cfg = EventComputationConfigFromCaseFile(cf)
	}

	candidates := extractDocumentEvents(content, cfg)

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
			candidate.AddExtraDay,
			cfg.CalendarScope,
			candidate.TriggerText,
			BuildUpcomingComputation(
				candidate.DateKind,
				candidate.AnchorDate,
				candidate.RelativeDays,
				candidate.IsBusinessDays,
				candidate.TriggerText,
			),
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
