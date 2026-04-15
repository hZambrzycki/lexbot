package documentapp

import "testing"

func TestExtractDocumentEvents_DetectsMultipleEventTypes(t *testing.T) {
	content := `
	DILIGENCIA DE ORDENACIÓN

	Notifíquese la resolución el 10/04/2026.
	Se señala juicio para el día 10/05/2026.
	Requiérase a la parte demandada antes del 20/04/2026.
	La comparecencia tendrá lugar el 15/05/2026.
	Se concede plazo hasta el 01/04/2026 para presentar alegaciones.
	`

	got := extractDocumentEvents(content)

	if len(got) != 5 {
		t.Fatalf("expected 5 events, got %d", len(got))
	}

	assertHasExtractedEvent(t, got, "notification", "2026-04-10")
	assertHasExtractedEvent(t, got, "hearing", "2026-05-10")
	assertHasExtractedEvent(t, got, "requirement", "2026-04-20")
	assertHasExtractedEvent(t, got, "appearance", "2026-05-15")
	assertHasExtractedEvent(t, got, "deadline", "2026-04-01")

	hearing := mustFindExtractedEvent(t, got, "hearing", "2026-05-10")
	if hearing.DateKind != extractedDateKindAbsolute {
		t.Fatalf("expected hearing date kind %q, got %q", extractedDateKindAbsolute, hearing.DateKind)
	}
	if hearing.AnchorDate != "2026-05-10" {
		t.Fatalf("expected hearing anchor date %q, got %q", "2026-05-10", hearing.AnchorDate)
	}
	if hearing.AnchorSource != anchorSourceInline {
		t.Fatalf("expected hearing anchor source %q, got %q", anchorSourceInline, hearing.AnchorSource)
	}
}

func TestExtractDocumentEvents_DeduplicatesSameLineEvent(t *testing.T) {
	content := `
	Se señala juicio para el día 10/05/2026.
	Se señala juicio para el día 10/05/2026.
	`

	got := extractDocumentEvents(content)

	if len(got) != 1 {
		t.Fatalf("expected 1 deduplicated event, got %d", len(got))
	}

	if got[0].EventType != "hearing" {
		t.Fatalf("expected event type %q, got %q", "hearing", got[0].EventType)
	}

	if got[0].EventDate != "2026-05-10" {
		t.Fatalf("expected event date %q, got %q", "2026-05-10", got[0].EventDate)
	}

	if got[0].DateKind != extractedDateKindAbsolute {
		t.Fatalf("expected date kind %q, got %q", extractedDateKindAbsolute, got[0].DateKind)
	}
}

func TestExtractDocumentEvents_IgnoresLinesWithoutRecognizedEventType(t *testing.T) {
	content := `
	Este documento menciona la fecha 10/05/2026 pero no contiene vocabulario procesal relevante.
	`

	got := extractDocumentEvents(content)

	if len(got) != 0 {
		t.Fatalf("expected 0 events, got %d", len(got))
	}
}

func TestExtractDocumentEvents_DetectsFilingEvent(t *testing.T) {
	content := `
	La presentación de demanda deberá realizarse antes del 18/04/2026.
	`

	got := extractDocumentEvents(content)

	if len(got) != 1 {
		t.Fatalf("expected 1 event, got %d", len(got))
	}

	if got[0].EventType != "filing" {
		t.Fatalf("expected event type %q, got %q", "filing", got[0].EventType)
	}

	if got[0].EventDate != "2026-04-18" {
		t.Fatalf("expected event date %q, got %q", "2026-04-18", got[0].EventDate)
	}

	if got[0].DateKind != extractedDateKindAbsolute {
		t.Fatalf("expected date kind %q, got %q", extractedDateKindAbsolute, got[0].DateKind)
	}
	if got[0].AnchorDate != "2026-04-18" {
		t.Fatalf("expected anchor date %q, got %q", "2026-04-18", got[0].AnchorDate)
	}
	if got[0].AnchorSource != anchorSourceInline {
		t.Fatalf("expected anchor source %q, got %q", anchorSourceInline, got[0].AnchorSource)
	}
	if got[0].RelativeDays != 0 {
		t.Fatalf("expected relative days 0, got %d", got[0].RelativeDays)
	}
	if got[0].IsBusinessDays {
		t.Fatalf("expected business days false for absolute date")
	}
	if got[0].TriggerText != "" {
		t.Fatalf("expected empty trigger text, got %q", got[0].TriggerText)
	}
}

func assertHasExtractedEvent(t *testing.T, events []extractedEventCandidate, wantType, wantDate string) {
	t.Helper()

	for _, e := range events {
		if e.EventType == wantType && e.EventDate == wantDate {
			return
		}
	}

	t.Fatalf("expected event type=%q date=%q not found", wantType, wantDate)
}

func mustFindExtractedEvent(t *testing.T, events []extractedEventCandidate, wantType, wantDate string) extractedEventCandidate {
	t.Helper()

	for _, e := range events {
		if e.EventType == wantType && e.EventDate == wantDate {
			return e
		}
	}

	t.Fatalf("expected event type=%q date=%q not found", wantType, wantDate)
	return extractedEventCandidate{}
}
