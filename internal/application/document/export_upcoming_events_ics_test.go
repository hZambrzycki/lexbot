package documentapp

import (
	"strings"
	"testing"
)

func unfoldICS(value string) string {
	value = strings.ReplaceAll(value, "\r\n ", "")
	value = strings.ReplaceAll(value, "\n ", "")
	return value
}

func TestBuildICSCalendar_IncludesComputationAndTrigger(t *testing.T) {
	events := []UpcomingEvent{
		{
			EventID:        "1",
			DocumentID:     "doc-1",
			OriginalName:   "test.txt",
			EventType:      "deadline",
			EventDate:      "2026-04-16",
			SourceText:     "Se concede plazo de 3 días hábiles desde el siguiente día hábil.",
			DaysRemaining:  0,
			Status:         "today",
			Priority:       "critical",
			DuplicateCount: 1,
			DocumentNames:  []string{"test.txt"},
			DocumentIDs:    []string{"doc-1"},

			DateKind:       "relative",
			AnchorDate:     "2026-04-11",
			AnchorSource:   "notification_line",
			RelativeDays:   3,
			IsBusinessDays: true,
			AddExtraDay:    true,
			TriggerText:    "plazo de 3 dias habiles desde el siguiente dia habil",
			Computation:    "anchor + next_day + 3 business days",
		},
	}

	ics := unfoldICS(buildICSCalendar(events))

	if !strings.Contains(ics, "BEGIN:VEVENT") {
		t.Fatal("expected VEVENT block")
	}

	if !strings.Contains(ics, "Cómputo: Fecha base + día siguiente + 3 días hábiles") {
		t.Fatal("missing translated computation in DESCRIPTION")
	}

	if !strings.Contains(ics, "Disparador textual: plazo de 3 dias habiles desde el siguiente dia habil") {
		t.Fatal("missing trigger in DESCRIPTION")
	}

	if !strings.Contains(ics, "Fecha base: 2026-04-11") {
		t.Fatal("missing anchor date in DESCRIPTION")
	}

	if !strings.Contains(ics, "Días hábiles: sí") {
		t.Fatal("missing business days flag")
	}

	if !strings.Contains(ics, "Añade día siguiente: sí") {
		t.Fatal("missing add extra day flag")
	}
}

func TestBuildICSCalendar_AbsoluteEventIncludesComputation(t *testing.T) {
	events := []UpcomingEvent{
		{
			EventID:       "2",
			EventType:     "notification",
			EventDate:     "2026-04-11",
			Status:        "overdue",
			Priority:      "low",
			DateKind:      "absolute",
			Computation:   "absolute date",
			SourceText:    "Notifíquese la resolución el 11/04/2026.",
			DocumentNames: []string{"notif.txt"},
			DocumentIDs:   []string{"doc-2"},
		},
	}

	ics := unfoldICS(buildICSCalendar(events))

	if !strings.Contains(ics, "Cómputo: Fecha absoluta") {
		t.Fatal("expected absolute computation in DESCRIPTION")
	}

	if !strings.Contains(ics, "Tipo de fecha: Fecha absoluta") {
		t.Fatal("expected absolute date kind")
	}
}

func TestBuildICSCalendar_MultipleEventsRemainSeparated(t *testing.T) {
	events := []UpcomingEvent{
		{
			EventID:      "1",
			EventType:    "deadline",
			EventDate:    "2026-04-16",
			Priority:     "critical",
			Status:       "today",
			DateKind:     "relative",
			RelativeDays: 3,
			Computation:  "anchor + 3 business days",
		},
		{
			EventID:      "2",
			EventType:    "deadline",
			EventDate:    "2026-04-16",
			Priority:     "critical",
			Status:       "today",
			DateKind:     "relative",
			RelativeDays: 5,
			Computation:  "anchor + 5 natural days",
		},
	}

	ics := buildICSCalendar(events)

	count := strings.Count(ics, "BEGIN:VEVENT")
	if count != 2 {
		t.Fatalf("expected 2 events, got %d", count)
	}
}

func TestBuildICSCalendar_UsesSpanishCalendarMetadata(t *testing.T) {
	events := []UpcomingEvent{
		{
			EventID:     "1",
			EventType:   "deadline",
			EventDate:   "2026-04-16",
			Priority:    "critical",
			Status:      "today",
			DateKind:    "relative",
			Computation: "anchor + 3 business days",
		},
	}

	ics := buildICSCalendar(events)

	if !strings.Contains(ics, "PRODID:-//LEXBOX//Agenda Procesal//ES") {
		t.Fatal("expected Spanish PRODID")
	}

	if !strings.Contains(ics, "X-WR-CALNAME:LEXBOX - Agenda procesal") {
		t.Fatal("expected calendar name")
	}
}

func TestEscapeICSText(t *testing.T) {
	input := "texto, con; caracteres \\ raros\nnueva linea"
	got := escapeICSText(input)

	if strings.Contains(got, "\n") {
		t.Fatal("newline should be escaped")
	}

	if !strings.Contains(got, `\,`) {
		t.Fatal("comma not escaped")
	}

	if !strings.Contains(got, `\;`) {
		t.Fatal("semicolon not escaped")
	}

	if !strings.Contains(got, `\\`) {
		t.Fatal("backslash not escaped")
	}

	if !strings.Contains(got, `\n`) {
		t.Fatal("newline escape sequence missing")
	}
}

func TestWriteICSLine_FoldsLongLines(t *testing.T) {
	var sb strings.Builder

	writeICSLine(&sb, "DESCRIPTION", strings.Repeat("a", 120))

	got := sb.String()

	if !strings.Contains(got, "\r\n ") {
		t.Fatal("expected folded ICS line continuation")
	}
}
