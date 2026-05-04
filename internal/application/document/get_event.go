package documentapp

import (
	"context"
	"strings"

	"lexbox/internal/application/ports"
	"lexbox/internal/domain/shared"
)

type GetEventInput struct {
	EventID string
}

type GetEvent struct {
	Events ports.DocumentEventRepository
}

func (uc GetEvent) Execute(ctx context.Context, in GetEventInput) (UpcomingEvent, error) {
	eventID := shared.NewID(strings.TrimSpace(in.EventID))
	if eventID == "" {
		return UpcomingEvent{}, shared.ErrInvalidID
	}

	event, err := uc.Events.GetDetailByID(ctx, eventID)
	if err != nil {
		return UpcomingEvent{}, err
	}

	computation := strings.TrimSpace(event.Computation)
	if computation == "" {
		computation = BuildUpcomingComputation(
			event.DateKind,
			event.AnchorDate,
			event.RelativeDays,
			event.IsBusinessDays,
			event.TriggerText,
		)
	}

	return UpcomingEvent{
		EventID:      event.EventID,
		DocumentID:   event.DocumentID,
		OriginalName: event.OriginalName,

		CaseFileID:        event.CaseFileID,
		CaseFileReference: event.CaseFileReference,
		CaseFileTitle:     event.CaseFileTitle,

		EventType:      event.EventType,
		EventDate:      event.EventDate,
		SourceText:     event.SourceText,
		DaysRemaining:  0,
		Status:         "",
		Priority:       "",
		DuplicateCount: 1,
		DocumentNames:  []string{event.OriginalName},
		DocumentIDs:    []string{event.DocumentID},

		AnchorDate:     event.AnchorDate,
		DateKind:       event.DateKind,
		AnchorSource:   event.AnchorSource,
		RelativeDays:   event.RelativeDays,
		IsBusinessDays: event.IsBusinessDays,
		AddExtraDay:    event.AddExtraDay,
		CalendarScope:  event.CalendarScope,
		TriggerText:    event.TriggerText,
		Computation:    computation,

		ReviewStatus:   event.ReviewStatus,
		ReviewedAt:     event.ReviewedAt,
		ResolvedAt:     event.ResolvedAt,
		ResolutionNote: event.ResolutionNote,
	}, nil
}
