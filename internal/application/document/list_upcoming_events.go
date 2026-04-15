package documentapp

import (
	"context"
	"sort"
	"strings"
	"time"

	"lexbox/internal/application/ports"
	"lexbox/internal/domain/shared"
)

type ListUpcomingEventsInput struct {
	CaseFileID   string
	EventType    string
	RelativeOnly bool
}

type UpcomingEvent struct {
	EventID      string
	DocumentID   string
	OriginalName string
	EventType    string
	EventDate    string
	SourceText   string

	DaysRemaining int
	Status        string
	Priority      string

	DuplicateCount int
	DocumentNames  []string
	DocumentIDs    []string

	AnchorDate     string
	DateKind       string
	AnchorSource   string
	RelativeDays   int
	IsBusinessDays bool
	TriggerText    string
}

type ListUpcomingEvents struct {
	Events ports.DocumentEventRepository
}

func (uc ListUpcomingEvents) Execute(ctx context.Context, in ListUpcomingEventsInput) ([]UpcomingEvent, error) {
	var caseFileID shared.ID
	if strings.TrimSpace(in.CaseFileID) != "" {
		caseFileID = shared.NewID(strings.TrimSpace(in.CaseFileID))
		if caseFileID == "" {
			return nil, shared.ErrInvalidID
		}
	}

	eventType := strings.TrimSpace(strings.ToLower(in.EventType))

	results, err := uc.Events.ListUpcoming(ctx, "", caseFileID, eventType)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	enriched := make([]UpcomingEvent, 0, len(results))

	for _, e := range results {
		if in.RelativeOnly && e.DateKind != "relative" {
			continue
		}

		eventDate, err := time.Parse("2006-01-02", e.EventDate)
		if err != nil {
			continue
		}

		eventDay := time.Date(eventDate.Year(), eventDate.Month(), eventDate.Day(), 0, 0, 0, 0, now.Location())
		diff := int(eventDay.Sub(today).Hours() / 24)

		status := ""
		switch {
		case diff < 0:
			status = "overdue"
		case diff == 0:
			status = "today"
		default:
			status = "upcoming"
		}

		priority := classifyUpcomingPriority(e.EventType, diff)

		enriched = append(enriched, UpcomingEvent{
			EventID:        e.EventID,
			DocumentID:     e.DocumentID,
			OriginalName:   e.OriginalName,
			EventType:      e.EventType,
			EventDate:      e.EventDate,
			SourceText:     e.SourceText,
			DaysRemaining:  diff,
			Status:         status,
			Priority:       priority,
			DuplicateCount: 1,
			DocumentNames:  []string{e.OriginalName},
			DocumentIDs:    []string{e.DocumentID},
			AnchorDate:     e.AnchorDate,
			DateKind:       e.DateKind,
			AnchorSource:   e.AnchorSource,
			RelativeDays:   e.RelativeDays,
			IsBusinessDays: e.IsBusinessDays,
			TriggerText:    e.TriggerText,
		})
	}

	return deduplicateByDateAndType(enriched), nil
}

func classifyUpcomingPriority(eventType string, daysRemaining int) string {
	switch eventType {
	case "deadline":
		if daysRemaining <= 0 {
			return "critical"
		}
		return "high"

	case "hearing":
		if daysRemaining <= 0 {
			return "critical"
		}
		return "high"

	case "appearance":
		if daysRemaining <= 0 {
			return "high"
		}
		return "medium"

	case "requirement":
		return "medium"

	case "filing":
		if daysRemaining <= 0 {
			return "high"
		}
		return "medium"

	case "notification":
		return "low"

	default:
		if daysRemaining < 0 {
			return "high"
		}
		return "low"
	}
}

func deduplicateByDateAndType(items []UpcomingEvent) []UpcomingEvent {
	type key struct {
		EventDate string
		EventType string
	}

	seen := make(map[key]int, len(items))
	result := make([]UpcomingEvent, 0, len(items))

	for _, item := range items {
		k := key{
			EventDate: item.EventDate,
			EventType: item.EventType,
		}

		if idx, exists := seen[k]; exists {
			result[idx].DuplicateCount++

			if !strings.Contains(result[idx].SourceText, item.SourceText) {
				result[idx].SourceText += " | " + item.SourceText
			}

			result[idx].DocumentNames = appendUniqueString(result[idx].DocumentNames, item.OriginalName)
			result[idx].DocumentIDs = appendUniqueString(result[idx].DocumentIDs, item.DocumentID)

			if result[idx].DateKind == "" && item.DateKind != "" {
				result[idx].DateKind = item.DateKind
				result[idx].AnchorDate = item.AnchorDate
				result[idx].AnchorSource = item.AnchorSource
				result[idx].RelativeDays = item.RelativeDays
				result[idx].IsBusinessDays = item.IsBusinessDays
				result[idx].TriggerText = item.TriggerText
			}

			if priorityRank(item.Priority) < priorityRank(result[idx].Priority) {
				result[idx].Priority = item.Priority
			}

			continue
		}

		seen[k] = len(result)
		result = append(result, item)
	}

	for i := range result {
		sort.Strings(result[i].DocumentNames)
		sort.Strings(result[i].DocumentIDs)
	}

	sort.SliceStable(result, func(i, j int) bool {
		si := statusRank(result[i].Status)
		sj := statusRank(result[j].Status)
		if si != sj {
			return si < sj
		}

		if result[i].EventDate != result[j].EventDate {
			return result[i].EventDate < result[j].EventDate
		}

		pi := priorityRank(result[i].Priority)
		pj := priorityRank(result[j].Priority)
		if pi != pj {
			return pi < pj
		}

		return result[i].EventType < result[j].EventType
	})

	return result
}

func statusRank(status string) int {
	switch status {
	case "overdue":
		return 0
	case "today":
		return 1
	case "upcoming":
		return 2
	default:
		return 3
	}
}

func priorityRank(priority string) int {
	switch priority {
	case "critical":
		return 0
	case "high":
		return 1
	case "medium":
		return 2
	case "low":
		return 3
	default:
		return 4
	}
}

func appendUniqueString(values []string, candidate string) []string {
	candidate = strings.TrimSpace(candidate)
	if candidate == "" {
		return values
	}

	for _, value := range values {
		if value == candidate {
			return values
		}
	}

	return append(values, candidate)
}
