package documentapp

import "testing"

func TestExtractDocumentEvents_DetectsRelativeDeadlineFromPreviousAbsoluteDate(t *testing.T) {
	content := `
	Notifíquese la resolución el 11 de abril de 2026.
	Se concede plazo de 5 días para formular alegaciones.
	`

	got := extractDocumentEvents(content)

	if len(got) != 2 {
		t.Fatalf("expected 2 events, got %d", len(got))
	}

	assertHasExtractedEvent(t, got, "notification", "2026-04-11")
	assertHasExtractedEvent(t, got, "deadline", "2026-04-16")

	deadline := mustFindExtractedEvent(t, got, "deadline", "2026-04-16")
	if deadline.DateKind != extractedDateKindRelative {
		t.Fatalf("expected relative date kind, got %q", deadline.DateKind)
	}
	if deadline.AnchorDate != "2026-04-11" {
		t.Fatalf("expected anchor date 2026-04-11, got %q", deadline.AnchorDate)
	}
	if deadline.AnchorSource != anchorSourceNotificationLine {
		t.Fatalf("expected anchor source %q, got %q", anchorSourceNotificationLine, deadline.AnchorSource)
	}
	if deadline.RelativeDays != 5 {
		t.Fatalf("expected relative days 5, got %d", deadline.RelativeDays)
	}
	if deadline.IsBusinessDays {
		t.Fatalf("expected natural days, got business days")
	}
	if deadline.TriggerText == "" {
		t.Fatal("expected non-empty trigger text")
	}
}

func TestExtractDocumentEvents_DetectsRelativeDeadlineWithTextualNumber(t *testing.T) {
	content := `
	Notifíquese la resolución el 11 de abril de 2026.
	En el plazo de cinco días deberá presentar escrito.
	`

	got := extractDocumentEvents(content)

	assertHasExtractedEvent(t, got, "deadline", "2026-04-16")

	deadline := mustFindExtractedEvent(t, got, "deadline", "2026-04-16")
	if deadline.RelativeDays != 5 {
		t.Fatalf("expected relative days 5, got %d", deadline.RelativeDays)
	}
	if deadline.IsBusinessDays {
		t.Fatalf("expected natural days, got business days")
	}
	if deadline.AnchorSource != anchorSourceNotificationLine {
		t.Fatalf("expected anchor source %q, got %q", anchorSourceNotificationLine, deadline.AnchorSource)
	}
}

func TestExtractDocumentEvents_IgnoresRelativeDeadlineWithoutAnchorDate(t *testing.T) {
	content := `
	Se concede plazo de 5 días para formular alegaciones.
	`

	got := extractDocumentEvents(content)

	if len(got) != 0 {
		t.Fatalf("expected 0 events without anchor date, got %d", len(got))
	}
}

func TestExtractDocumentEvents_PrioritizesNotificationDateOverGenericPreviousAbsoluteDate(t *testing.T) {
	content := `
	DILIGENCIA DE ORDENACIÓN de 09/04/2026.
	Notifíquese la presente resolución a las partes el 11/04/2026.
	Se concede plazo de 5 días para formular alegaciones.
	`

	got := extractDocumentEvents(content)

	assertHasExtractedEvent(t, got, "notification", "2026-04-11")
	assertHasExtractedEvent(t, got, "deadline", "2026-04-16")

	deadline := mustFindExtractedEvent(t, got, "deadline", "2026-04-16")
	if deadline.AnchorDate != "2026-04-11" {
		t.Fatalf("expected anchor date 2026-04-11, got %q", deadline.AnchorDate)
	}
	if deadline.AnchorSource != anchorSourceNotificationLine {
		t.Fatalf("expected anchor source %q, got %q", anchorSourceNotificationLine, deadline.AnchorSource)
	}
}

func TestExtractDocumentEvents_DetectsNextDayAfterNotification(t *testing.T) {
	content := `
	Notifíquese la resolución el 11/04/2026.
	Al día siguiente de la notificación deberá aportar la documentación.
	`

	got := extractDocumentEvents(content)

	assertHasExtractedEvent(t, got, "notification", "2026-04-11")
	assertHasExtractedEvent(t, got, "deadline", "2026-04-12")

	event := mustFindExtractedEvent(t, got, "deadline", "2026-04-12")
	if event.RelativeDays != 1 {
		t.Fatalf("expected relative days 1, got %d", event.RelativeDays)
	}
	if event.AnchorDate != "2026-04-11" {
		t.Fatalf("expected anchor date 2026-04-11, got %q", event.AnchorDate)
	}
	if event.AnchorSource != anchorSourceNotificationLine {
		t.Fatalf("expected anchor source %q, got %q", anchorSourceNotificationLine, event.AnchorSource)
	}
}
