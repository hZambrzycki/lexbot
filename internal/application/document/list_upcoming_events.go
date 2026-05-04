package documentapp

import (
	"context"
	"sort"
	"strconv"
	"strings"
	"time"

	"lexbox/internal/application/ports"
	"lexbox/internal/domain/document"
	"lexbox/internal/domain/shared"
)

type ListUpcomingEventsInput struct {
	CaseFileID   string
	EventType    string
	RelativeOnly bool
	ReviewStatus string
}

type UpcomingEvent struct {
	EventID      string
	DocumentID   string
	OriginalName string

	CaseFileID        string
	CaseFileReference string
	CaseFileTitle     string

	EventType      string
	EventDate      string
	SourceText     string
	DaysRemaining  int
	Status         string
	Priority       string
	DuplicateCount int
	DocumentNames  []string
	DocumentIDs    []string

	AnchorDate     string
	DateKind       string
	AnchorSource   string
	RelativeDays   int
	IsBusinessDays bool
	AddExtraDay    bool
	CalendarScope  string
	TriggerText    string
	Computation    string

	ReviewStatus   string
	ReviewedAt     string
	ResolvedAt     string
	ResolutionNote string
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
	reviewStatus := strings.TrimSpace(strings.ToLower(in.ReviewStatus))

	if reviewStatus != "" && !document.IsValidReviewStatus(reviewStatus) {
		return nil, shared.ErrInvalidAssociation
	}

	results, err := uc.Events.ListUpcoming(ctx, "", caseFileID, eventType, reviewStatus)
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

		computation := strings.TrimSpace(e.Computation)
		if computation == "" {
			computation = BuildUpcomingComputation(
				e.DateKind,
				e.AnchorDate,
				e.RelativeDays,
				e.IsBusinessDays,
				e.TriggerText,
			)
		}

		enriched = append(enriched, UpcomingEvent{
			EventID:      e.EventID,
			DocumentID:   e.DocumentID,
			OriginalName: e.OriginalName,

			CaseFileID:        e.CaseFileID,
			CaseFileReference: e.CaseFileReference,
			CaseFileTitle:     e.CaseFileTitle,

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
			AddExtraDay:    e.AddExtraDay,
			CalendarScope:  e.CalendarScope,
			TriggerText:    e.TriggerText,
			Computation:    computation,

			ReviewStatus:   e.ReviewStatus,
			ReviewedAt:     e.ReviewedAt,
			ResolvedAt:     e.ResolvedAt,
			ResolutionNote: e.ResolutionNote,
		})
	}

	return deduplicateUpcomingEvents(enriched), nil
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

func deduplicateUpcomingEvents(items []UpcomingEvent) []UpcomingEvent {
	type key struct {
		CaseFileID     string
		EventDate      string
		EventType      string
		DateKind       string
		AnchorDate     string
		AnchorSource   string
		RelativeDays   int
		IsBusinessDays bool
		TriggerText    string
	}

	seen := make(map[key]int, len(items))
	result := make([]UpcomingEvent, 0, len(items))

	for _, item := range items {
		k := key{
			CaseFileID:     strings.TrimSpace(item.CaseFileID),
			EventDate:      item.EventDate,
			EventType:      item.EventType,
			DateKind:       strings.TrimSpace(item.DateKind),
			AnchorDate:     strings.TrimSpace(item.AnchorDate),
			AnchorSource:   strings.TrimSpace(item.AnchorSource),
			RelativeDays:   item.RelativeDays,
			IsBusinessDays: item.IsBusinessDays,
			TriggerText:    strings.TrimSpace(item.TriggerText),
		}

		if idx, exists := seen[k]; exists {
			result[idx].DuplicateCount++

			if !strings.Contains(result[idx].SourceText, item.SourceText) {
				result[idx].SourceText += " | " + item.SourceText
			}

			result[idx].DocumentNames = appendUniqueString(result[idx].DocumentNames, item.OriginalName)
			result[idx].DocumentIDs = appendUniqueString(result[idx].DocumentIDs, item.DocumentID)

			if priorityRank(item.Priority) < priorityRank(result[idx].Priority) {
				result[idx].Priority = item.Priority
			}

			if reviewPriorityRank(item.ReviewStatus) < reviewPriorityRank(result[idx].ReviewStatus) {
				result[idx].ReviewStatus = item.ReviewStatus
				result[idx].ReviewedAt = item.ReviewedAt
				result[idx].ResolvedAt = item.ResolvedAt
				result[idx].ResolutionNote = item.ResolutionNote
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

		ri := reviewPriorityRank(result[i].ReviewStatus)
		rj := reviewPriorityRank(result[j].ReviewStatus)
		if ri != rj {
			return ri < rj
		}

		if result[i].CaseFileReference != result[j].CaseFileReference {
			return result[i].CaseFileReference < result[j].CaseFileReference
		}

		if result[i].EventType != result[j].EventType {
			return result[i].EventType < result[j].EventType
		}

		if result[i].DateKind != result[j].DateKind {
			return result[i].DateKind < result[j].DateKind
		}

		if result[i].AnchorDate != result[j].AnchorDate {
			return result[i].AnchorDate < result[j].AnchorDate
		}

		if result[i].RelativeDays != result[j].RelativeDays {
			return result[i].RelativeDays < result[j].RelativeDays
		}

		if result[i].IsBusinessDays != result[j].IsBusinessDays {
			return !result[i].IsBusinessDays && result[j].IsBusinessDays
		}

		return result[i].TriggerText < result[j].TriggerText
	})

	return result
}

func BuildUpcomingComputation(dateKind, anchorDate string, relativeDays int, isBusinessDays bool, triggerText string) string {
	if strings.TrimSpace(dateKind) == "" {
		return ""
	}

	if dateKind != "relative" {
		return "absolute date"
	}

	if strings.TrimSpace(anchorDate) == "" || relativeDays <= 0 {
		return "relative date"
	}

	parts := []string{"anchor"}

	if triggerAddsExtraDay(triggerText) {
		parts = append(parts, "next_day")
	}

	dayLabel := "natural days"
	if isBusinessDays {
		dayLabel = "business days"
	}

	parts = append(parts, strconv.Itoa(relativeDays)+" "+dayLabel)

	return strings.Join(parts, " + ")
}

func triggerAddsExtraDay(triggerText string) bool {
	normalized := normalizeASCIIText(triggerText)

	return containsAny(normalized,
		"a contar desde el dia siguiente",
		"desde el dia siguiente",
		"a contar desde el siguiente",
		"desde el siguiente",
		"a partir del dia siguiente",
		"a partir del siguiente",
		"a contar desde el siguiente dia habil",
		"desde el siguiente dia habil",
		"a partir del siguiente dia habil",
		"a contar desde el dia habil siguiente",
		"desde el dia habil siguiente",
		"a partir del dia habil siguiente",
	)
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

func reviewPriorityRank(reviewStatus string) int {
	switch strings.TrimSpace(strings.ToLower(reviewStatus)) {
	case "pending", "":
		return 0
	case "reviewed":
		return 1
	case "resolved":
		return 2
	default:
		return 0
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
