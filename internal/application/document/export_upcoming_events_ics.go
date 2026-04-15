package documentapp

import (
	"context"
	"fmt"
	"strings"
	"time"
)

type ExportUpcomingEventsICSInput struct {
	CaseFileID string
	EventType  string
}

type ExportUpcomingEventsICSResult struct {
	Content    string
	EventCount int
}

type ExportUpcomingEventsICS struct {
	ListUpcomingEvents ListUpcomingEvents
}

func (uc ExportUpcomingEventsICS) Execute(ctx context.Context, in ExportUpcomingEventsICSInput) (ExportUpcomingEventsICSResult, error) {
	events, err := uc.ListUpcomingEvents.Execute(ctx, ListUpcomingEventsInput{
		CaseFileID: in.CaseFileID,
		EventType:  in.EventType,
	})
	if err != nil {
		return ExportUpcomingEventsICSResult{}, err
	}

	content := buildICSCalendar(events)

	return ExportUpcomingEventsICSResult{
		Content:    content,
		EventCount: len(events),
	}, nil
}

func buildICSCalendar(events []UpcomingEvent) string {
	var sb strings.Builder

	nowUTC := time.Now().UTC().Format("20060102T150405Z")

	sb.WriteString("BEGIN:VCALENDAR\r\n")
	sb.WriteString("VERSION:2.0\r\n")
	sb.WriteString("PRODID:-//LEXBOX//Upcoming Events//EN\r\n")
	sb.WriteString("CALSCALE:GREGORIAN\r\n")
	sb.WriteString("METHOD:PUBLISH\r\n")

	for i, e := range events {
		dateValue, err := time.Parse("2006-01-02", e.EventDate)
		if err != nil {
			continue
		}

		startDate := dateValue.Format("20060102")
		endDate := dateValue.AddDate(0, 0, 1).Format("20060102")

		summary := fmt.Sprintf("[%s] %s", strings.ToUpper(e.Priority), e.EventType)
		if len(e.DocumentNames) > 0 {
			summary = fmt.Sprintf("%s - %s", summary, strings.Join(e.DocumentNames, ", "))
		}

		descriptionParts := []string{
			fmt.Sprintf("Status: %s", e.Status),
			fmt.Sprintf("Priority: %s", e.Priority),
			fmt.Sprintf("Event type: %s", e.EventType),
			fmt.Sprintf("Event date: %s", e.EventDate),
		}

		if len(e.DocumentNames) > 0 {
			descriptionParts = append(descriptionParts, fmt.Sprintf("Documents: %s", strings.Join(e.DocumentNames, ", ")))
		}

		if len(e.DocumentIDs) > 0 {
			descriptionParts = append(descriptionParts, fmt.Sprintf("Document IDs: %s", strings.Join(e.DocumentIDs, ", ")))
		}

		if strings.TrimSpace(e.SourceText) != "" {
			descriptionParts = append(descriptionParts, fmt.Sprintf("Source: %s", e.SourceText))
		}

		if e.DateKind == "relative" {
			descriptionParts = append(descriptionParts, "Date kind: relative")

			if strings.TrimSpace(e.AnchorDate) != "" {
				descriptionParts = append(descriptionParts, fmt.Sprintf("Anchor date: %s", e.AnchorDate))
			}

			if strings.TrimSpace(e.AnchorSource) != "" {
				descriptionParts = append(descriptionParts, fmt.Sprintf("Anchor source: %s", e.AnchorSource))
			}

			if e.RelativeDays > 0 {
				descriptionParts = append(descriptionParts, fmt.Sprintf("Relative days: %d", e.RelativeDays))
			}

			descriptionParts = append(descriptionParts, fmt.Sprintf("Business days: %t", e.IsBusinessDays))

			if strings.TrimSpace(e.TriggerText) != "" {
				descriptionParts = append(descriptionParts, fmt.Sprintf("Trigger: %s", e.TriggerText))
			}
		}

		if e.DuplicateCount > 1 {
			descriptionParts = append(descriptionParts, fmt.Sprintf("Duplicate count: %d", e.DuplicateCount))
		}

		description := escapeICSText(strings.Join(descriptionParts, "\\n"))

		uidBase := e.EventID
		if uidBase == "" {
			uidBase = fmt.Sprintf("%s-%s-%d", e.EventType, e.EventDate, i+1)
		}

		sb.WriteString("BEGIN:VEVENT\r\n")
		sb.WriteString(fmt.Sprintf("UID:%s@lexbox\r\n", escapeICSText(uidBase)))
		sb.WriteString(fmt.Sprintf("DTSTAMP:%s\r\n", nowUTC))
		sb.WriteString(fmt.Sprintf("DTSTART;VALUE=DATE:%s\r\n", startDate))
		sb.WriteString(fmt.Sprintf("DTEND;VALUE=DATE:%s\r\n", endDate))
		sb.WriteString(fmt.Sprintf("SUMMARY:%s\r\n", escapeICSText(summary)))
		sb.WriteString(fmt.Sprintf("DESCRIPTION:%s\r\n", description))
		sb.WriteString("END:VEVENT\r\n")
	}

	sb.WriteString("END:VCALENDAR\r\n")

	return sb.String()
}

func escapeICSText(value string) string {
	value = strings.ReplaceAll(value, `\`, `\\`)
	value = strings.ReplaceAll(value, ";", `\;`)
	value = strings.ReplaceAll(value, ",", `\,`)
	value = strings.ReplaceAll(value, "\n", `\n`)
	value = strings.ReplaceAll(value, "\r", "")
	return value
}
