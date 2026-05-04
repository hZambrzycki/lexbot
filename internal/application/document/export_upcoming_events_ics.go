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
	sb.WriteString("PRODID:-//LEXBOX//Agenda Procesal//ES\r\n")
	sb.WriteString("CALSCALE:GREGORIAN\r\n")
	sb.WriteString("METHOD:PUBLISH\r\n")
	sb.WriteString("X-WR-CALNAME:LEXBOX - Agenda procesal\r\n")
	sb.WriteString("X-WR-CALDESC:Hitos procesales detectados automáticamente por LEXBOX\r\n")

	for i, e := range events {
		dateValue, err := time.Parse("2006-01-02", e.EventDate)
		if err != nil {
			continue
		}

		startDate := dateValue.Format("20060102")
		endDate := dateValue.AddDate(0, 0, 1).Format("20060102")

		summary := buildICSSummary(e)
		description := buildICSDescription(e)

		uidBase := e.EventID
		if strings.TrimSpace(uidBase) == "" {
			uidBase = fmt.Sprintf("%s-%s-%d", e.EventType, e.EventDate, i+1)
		}

		sb.WriteString("BEGIN:VEVENT\r\n")
		writeICSLine(&sb, "UID", fmt.Sprintf("%s@lexbox", uidBase))
		writeICSLine(&sb, "DTSTAMP", nowUTC)
		sb.WriteString(fmt.Sprintf("DTSTART;VALUE=DATE:%s\r\n", startDate))
		sb.WriteString(fmt.Sprintf("DTEND;VALUE=DATE:%s\r\n", endDate))
		writeICSLine(&sb, "SUMMARY", summary)
		writeICSLine(&sb, "DESCRIPTION", description)
		sb.WriteString("END:VEVENT\r\n")
	}

	sb.WriteString("END:VCALENDAR\r\n")

	return sb.String()
}

func buildICSSummary(e UpcomingEvent) string {
	priority := displayICSPriority(e.Priority)
	eventType := displayICSEventType(e.EventType)

	summary := fmt.Sprintf("[%s] %s", priority, eventType)

	if len(e.DocumentNames) > 0 {
		summary = fmt.Sprintf("%s - %s", summary, strings.Join(e.DocumentNames, ", "))
	}

	return summary
}

func buildICSDescription(e UpcomingEvent) string {
	descriptionParts := []string{
		fmt.Sprintf("Estado: %s", displayICSStatus(e.Status)),
		fmt.Sprintf("Prioridad: %s", displayICSPriority(e.Priority)),
		fmt.Sprintf("Tipo de evento: %s", displayICSEventType(e.EventType)),
		fmt.Sprintf("Fecha del evento: %s", e.EventDate),
	}

	if strings.TrimSpace(e.DateKind) != "" {
		descriptionParts = append(descriptionParts, fmt.Sprintf("Tipo de fecha: %s", displayICSDateKind(e.DateKind)))
	}

	if strings.TrimSpace(e.Computation) != "" {
		descriptionParts = append(descriptionParts, fmt.Sprintf("Cómputo: %s", displayICSComputation(e.Computation)))
	}

	if strings.TrimSpace(e.CalendarScope) != "" {
		descriptionParts = append(descriptionParts, fmt.Sprintf("Calendario aplicado: %s", e.CalendarScope))
	}

	if len(e.DocumentNames) > 0 {
		descriptionParts = append(descriptionParts, fmt.Sprintf("Documentos: %s", strings.Join(e.DocumentNames, ", ")))
	}

	if len(e.DocumentIDs) > 0 {
		descriptionParts = append(descriptionParts, fmt.Sprintf("IDs de documentos: %s", strings.Join(e.DocumentIDs, ", ")))
	}

	if strings.TrimSpace(e.SourceText) != "" {
		descriptionParts = append(descriptionParts, fmt.Sprintf("Texto origen: %s", e.SourceText))
	}

	if e.DateKind == "relative" {
		if strings.TrimSpace(e.AnchorDate) != "" {
			descriptionParts = append(descriptionParts, fmt.Sprintf("Fecha base: %s", e.AnchorDate))
		}

		if strings.TrimSpace(e.AnchorSource) != "" {
			descriptionParts = append(descriptionParts, fmt.Sprintf("Origen de fecha base: %s", displayICSAnchorSource(e.AnchorSource)))
		}

		if e.RelativeDays > 0 {
			descriptionParts = append(descriptionParts, fmt.Sprintf("Días relativos: %d", e.RelativeDays))
		}

		descriptionParts = append(descriptionParts, fmt.Sprintf("Días hábiles: %s", displayICSBool(e.IsBusinessDays)))
		descriptionParts = append(descriptionParts, fmt.Sprintf("Añade día siguiente: %s", displayICSBool(e.AddExtraDay)))

		if strings.TrimSpace(e.TriggerText) != "" {
			descriptionParts = append(descriptionParts, fmt.Sprintf("Disparador textual: %s", e.TriggerText))
		}
	}

	if e.DuplicateCount > 1 {
		descriptionParts = append(descriptionParts, fmt.Sprintf("Detecciones agrupadas: %d", e.DuplicateCount))
	}

	if strings.TrimSpace(e.ReviewStatus) != "" {
		descriptionParts = append(descriptionParts, fmt.Sprintf("Estado de revisión: %s", displayICSReviewStatus(e.ReviewStatus)))
	}

	if strings.TrimSpace(e.ResolutionNote) != "" {
		descriptionParts = append(descriptionParts, fmt.Sprintf("Nota de resolución: %s", e.ResolutionNote))
	}

	return strings.Join(descriptionParts, "\n")
}

func writeICSLine(sb *strings.Builder, key string, value string) {
	escaped := escapeICSText(value)
	line := fmt.Sprintf("%s:%s", key, escaped)

	for _, folded := range foldICSLine(line) {
		sb.WriteString(folded)
		sb.WriteString("\r\n")
	}
}

func foldICSLine(line string) []string {
	const limit = 75

	if len(line) <= limit {
		return []string{line}
	}

	lines := make([]string, 0)

	for len(line) > limit {
		cut := limit

		for cut > 0 && !isSafeUTF8Boundary(line, cut) {
			cut--
		}

		if cut <= 0 {
			cut = limit
		}

		lines = append(lines, line[:cut])
		line = " " + line[cut:]
	}

	lines = append(lines, line)

	return lines
}

func isSafeUTF8Boundary(value string, index int) bool {
	if index <= 0 || index >= len(value) {
		return true
	}

	return value[index]&0xC0 != 0x80
}

func escapeICSText(value string) string {
	value = strings.ReplaceAll(value, `\`, `\\`)
	value = strings.ReplaceAll(value, ";", `\;`)
	value = strings.ReplaceAll(value, ",", `\,`)
	value = strings.ReplaceAll(value, "\r\n", "\n")
	value = strings.ReplaceAll(value, "\r", "\n")
	value = strings.ReplaceAll(value, "\n", `\n`)
	return value
}

func displayICSEventType(value string) string {
	switch value {
	case "deadline":
		return "Plazo"
	case "notification":
		return "Notificación"
	case "requirement":
		return "Requerimiento"
	case "hearing":
		return "Vista"
	case "appearance":
		return "Comparecencia"
	case "filing":
		return "Presentación"
	default:
		if strings.TrimSpace(value) == "" {
			return "Evento"
		}
		return value
	}
}

func displayICSStatus(value string) string {
	switch value {
	case "overdue":
		return "Vencido"
	case "today":
		return "Hoy"
	case "upcoming":
		return "Próximo"
	default:
		if strings.TrimSpace(value) == "" {
			return "Sin estado"
		}
		return value
	}
}

func displayICSPriority(value string) string {
	switch value {
	case "critical":
		return "Crítico"
	case "high":
		return "Alta"
	case "medium":
		return "Media"
	case "low":
		return "Baja"
	default:
		if strings.TrimSpace(value) == "" {
			return "Sin prioridad"
		}
		return value
	}
}

func displayICSDateKind(value string) string {
	switch value {
	case "absolute":
		return "Fecha absoluta"
	case "relative":
		return "Fecha relativa"
	default:
		if strings.TrimSpace(value) == "" {
			return "Sin tipo"
		}
		return value
	}
}

func displayICSAnchorSource(value string) string {
	switch value {
	case "inline":
		return "Fecha en la misma línea"
	case "previous_line":
		return "Fecha en línea anterior"
	case "notification_line":
		return "Línea de notificación"
	case "procedural_anchor_line":
		return "Línea procesal de referencia"
	default:
		if strings.TrimSpace(value) == "" {
			return "Sin origen"
		}
		return value
	}
}

func displayICSReviewStatus(value string) string {
	switch value {
	case "pending":
		return "Pendiente"
	case "reviewed":
		return "Revisado"
	case "resolved":
		return "Resuelto"
	default:
		if strings.TrimSpace(value) == "" {
			return "Pendiente"
		}
		return value
	}
}

func displayICSBool(value bool) string {
	if value {
		return "sí"
	}
	return "no"
}

func displayICSComputation(value string) string {
	switch value {
	case "absolute date":
		return "Fecha absoluta"
	case "relative date":
		return "Fecha relativa"
	case "anchor + 1 business days":
		return "Fecha base + 1 día hábil"
	case "anchor + 1 natural days":
		return "Fecha base + 1 día natural"
	case "anchor + next_day + 1 business days":
		return "Fecha base + día siguiente + 1 día hábil"
	case "anchor + next_day + 1 natural days":
		return "Fecha base + día siguiente + 1 día natural"
	default:
		if strings.TrimSpace(value) == "" {
			return "-"
		}

		return strings.NewReplacer(
			"anchor", "Fecha base",
			"next_day", "día siguiente",
			"business days", "días hábiles",
			"natural days", "días naturales",
		).Replace(value)
	}
}
