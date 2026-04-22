package casefileapp

import (
	"context"
	"fmt"
	documentapp "lexbox/internal/application/document"
	"lexbox/internal/application/ports"
	"lexbox/internal/domain/casefile"
	"lexbox/internal/domain/shared"
	"strings"
)

type GetCaseFileDashboardInput struct {
	CaseFileID string
}

type CaseFileDashboardResult struct {
	CaseFile       casefile.CaseFile
	NoteCount      int
	DocumentCount  int
	UpcomingEvents []documentapp.UpcomingEvent

	DocumentsWithoutText     int
	DocumentsWithoutTextList []string

	DocumentsWithUnknownMetadata     int
	DocumentsWithUnknownMetadataList []string

	DocumentsWithoutEvents     int
	DocumentsWithoutEventsList []string

	OverdueCount  int
	TodayCount    int
	UpcomingCount int

	CriticalCount int
	HighCount     int
	MediumCount   int
	LowCount      int

	NeedsAttention        bool
	TopAlert              string
	RecommendedNextAction string
	ProceduralHint        string
}

type GetCaseFileDashboard struct {
	GetCaseFileDetail  GetCaseFileDetail
	ListUpcomingEvents documentapp.ListUpcomingEvents
	DocumentContents   ports.DocumentContentRepository
	Metadata           ports.DocumentMetadataRepository
	Events             ports.DocumentEventRepository
}

func (uc GetCaseFileDashboard) Execute(ctx context.Context, in GetCaseFileDashboardInput) (CaseFileDashboardResult, error) {
	detail, err := uc.GetCaseFileDetail.Execute(ctx, GetCaseFileDetailInput{
		ID: in.CaseFileID,
	})
	if err != nil {
		return CaseFileDashboardResult{}, err
	}

	events, err := uc.ListUpcomingEvents.Execute(ctx, documentapp.ListUpcomingEventsInput{
		CaseFileID: in.CaseFileID,
	})
	if err != nil {
		return CaseFileDashboardResult{}, err
	}

	result := CaseFileDashboardResult{
		CaseFile:                         detail.CaseFile,
		NoteCount:                        len(detail.Notes),
		DocumentCount:                    len(detail.Documents),
		UpcomingEvents:                   events,
		DocumentsWithoutTextList:         []string{},
		DocumentsWithUnknownMetadataList: []string{},
		DocumentsWithoutEventsList:       []string{},
	}

	for _, doc := range detail.Documents {
		if uc.DocumentContents != nil {
			content, err := uc.DocumentContents.GetByDocumentID(ctx, doc.ID.String())
			if err != nil {
				if err == shared.ErrNotFound {
					result.DocumentsWithoutText++
					result.DocumentsWithoutTextList = append(result.DocumentsWithoutTextList, doc.OriginalName)
				} else {
					return CaseFileDashboardResult{}, err
				}
			} else if strings.TrimSpace(content) == "" {
				result.DocumentsWithoutText++
				result.DocumentsWithoutTextList = append(result.DocumentsWithoutTextList, doc.OriginalName)
			}
		}

		if uc.Metadata != nil {
			meta, err := uc.Metadata.GetByDocumentID(ctx, doc.ID)
			if err != nil {
				if err == shared.ErrNotFound {
					result.DocumentsWithUnknownMetadata++
					result.DocumentsWithUnknownMetadataList = append(result.DocumentsWithUnknownMetadataList, doc.OriginalName)
				} else {
					return CaseFileDashboardResult{}, err
				}
			} else {
				if strings.TrimSpace(meta.DocumentType) == "" || meta.DocumentType == "unknown" ||
					strings.TrimSpace(meta.LegalArea) == "" || meta.LegalArea == "unknown" {
					result.DocumentsWithUnknownMetadata++
					result.DocumentsWithUnknownMetadataList = append(result.DocumentsWithUnknownMetadataList, doc.OriginalName)
				}
			}
		}

		if uc.Events != nil {
			docEvents, err := uc.Events.ListByDocumentID(ctx, doc.ID)
			if err != nil {
				return CaseFileDashboardResult{}, err
			}
			if len(docEvents) == 0 {
				result.DocumentsWithoutEvents++
				result.DocumentsWithoutEventsList = append(result.DocumentsWithoutEventsList, doc.OriginalName)
			}
		}
	}

	for _, e := range events {
		switch e.Status {
		case "overdue":
			result.OverdueCount++
		case "today":
			result.TodayCount++
		case "upcoming":
			result.UpcomingCount++
		}

		switch e.Priority {
		case "critical":
			result.CriticalCount++
		case "high":
			result.HighCount++
		case "medium":
			result.MediumCount++
		case "low":
			result.LowCount++
		}
	}
	result.RecommendedNextAction = buildDashboardRecommendedAction(events)
	result.ProceduralHint = buildDashboardProceduralHint(events)
	result.NeedsAttention, result.TopAlert = buildDashboardAlert(result)
	return result, nil
}

func buildDashboardAlert(result CaseFileDashboardResult) (bool, string) {
	best := selectTopEvent(result.UpcomingEvents)

	if best != nil {
		statusText := ""

		if best.Status == "overdue" {
			statusText = fmt.Sprintf("%d days ago", -best.DaysRemaining)
		} else if best.Status == "today" {
			statusText = "today"
		} else {
			statusText = fmt.Sprintf("in %d days", best.DaysRemaining)
		}

		return true, fmt.Sprintf(
			"%s %s on %s (%s)",
			best.Priority,
			best.EventType,
			best.EventDate,
			statusText,
		)
	}

	switch {
	case result.DocumentsWithoutText > 0:
		return true, "some documents have no extracted text"

	case result.DocumentsWithUnknownMetadata > 0:
		return true, "some documents still have unknown metadata"

	case result.DocumentsWithoutEvents > 0:
		return true, "some documents have no detected events"

	default:
		return false, "no immediate alerts"
	}
}

func selectTopEvent(events []documentapp.UpcomingEvent) *documentapp.UpcomingEvent {
	if len(events) == 0 {
		return nil
	}

	best := &events[0]

	for i := range events {
		e := &events[i]

		if isBetterEvent(e, best) {
			best = e
		}
	}

	return best
}

func isBetterEvent(a, b *documentapp.UpcomingEvent) bool {
	// 1. status priority: overdue > today > upcoming
	if statusRank(a.Status) != statusRank(b.Status) {
		return statusRank(a.Status) < statusRank(b.Status)
	}

	// 2. priority: critical > high > medium > low
	if priorityRank(a.Priority) != priorityRank(b.Priority) {
		return priorityRank(a.Priority) < priorityRank(b.Priority)
	}

	// 3. closer date
	return a.DaysRemaining < b.DaysRemaining
}

func statusRank(s string) int {
	switch s {
	case "overdue":
		return 0
	case "today":
		return 1
	default:
		return 2
	}
}

func priorityRank(p string) int {
	switch p {
	case "critical":
		return 0
	case "high":
		return 1
	case "medium":
		return 2
	default:
		return 3
	}
}

func shortDocumentName(names []string) string {
	if len(names) == 0 {
		return ""
	}
	return names[0]
}

func buildDashboardRecommendedAction(events []documentapp.UpcomingEvent) string {
	best := selectTopEvent(events)
	if best == nil {
		return "no immediate procedural action"
	}

	verb := actionVerbForEventType(best.EventType)

	doc := shortDocumentName(best.DocumentNames)
	docPart := ""
	if doc != "" {
		docPart = fmt.Sprintf(" (%s)", doc)
	}

	switch best.Status {
	case "overdue":
		if best.Priority == "critical" {
			return fmt.Sprintf(
				"%s overdue critical %s%s from %s immediately",
				verb,
				best.EventType,
				docPart,
				best.EventDate,
			)
		}
		return fmt.Sprintf(
			"%s overdue %s%s from %s",
			verb,
			best.EventType,
			docPart,
			best.EventDate,
		)

	case "today":
		return fmt.Sprintf(
			"%s %s%s due today (%s)",
			verb,
			best.EventType,
			docPart,
			best.EventDate,
		)

	case "upcoming":
		if best.Priority == "high" || best.Priority == "critical" {
			return fmt.Sprintf(
				"%s upcoming %s%s for %s",
				verb,
				best.EventType,
				docPart,
				best.EventDate,
			)
		}
		return fmt.Sprintf(
			"monitor upcoming %s%s for %s",
			best.EventType,
			docPart,
			best.EventDate,
		)

	default:
		return "no immediate procedural action"
	}
}

func actionVerbForEventType(eventType string) string {
	switch eventType {
	case "deadline":
		return "review"
	case "hearing":
		return "prepare"
	case "appearance":
		return "confirm"
	case "notification":
		return "review"
	case "requirement":
		return "comply with"
	case "filing":
		return "prepare"
	default:
		return "review"
	}
}

func buildDashboardProceduralHint(events []documentapp.UpcomingEvent) string {
	best := selectTopEvent(events)
	if best == nil {
		return "no immediate procedural concerns"
	}

	return proceduralHintForEvent(best)
}

func proceduralHintForEvent(e *documentapp.UpcomingEvent) string {
	if e == nil {
		return "no immediate procedural concerns"
	}

	switch e.EventType {
	case "deadline":
		switch e.Status {
		case "overdue":
			return "possible deadline breach"
		case "today":
			return "deadline expires today"
		case "upcoming":
			return "review time remaining and filing readiness"
		}

	case "hearing":
		switch e.Status {
		case "overdue":
			return "review hearing outcome or non-appearance consequences"
		case "today":
			return "prepare hearing strategy and appearance"
		case "upcoming":
			return "prepare hearing strategy"
		}

	case "appearance":
		switch e.Status {
		case "overdue":
			return "review attendance status and possible consequences"
		case "today":
			return "confirm attendance and required documents"
		case "upcoming":
			return "confirm attendance and documents"
		}

	case "requirement":
		switch e.Status {
		case "overdue":
			return "review compliance status"
		case "today":
			return "complete the requirement today"
		case "upcoming":
			return "prepare compliance response"
		}

	case "filing":
		switch e.Status {
		case "overdue":
			return "review missed filing risk"
		case "today":
			return "finalize filing today"
		case "upcoming":
			return "prepare filing materials"
		}

	case "notification":
		switch e.Status {
		case "overdue":
			return "review whether the notification triggered any expired term"
		case "today":
			return "review consequences of the notification"
		case "upcoming":
			return "monitor notification-related deadlines"
		}
	}

	return "review procedural context"
}
