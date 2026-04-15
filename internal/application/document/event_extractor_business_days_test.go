package documentapp

import "testing"

func TestExtractDocumentEvents_DetectsRelativeDeadlineInBusinessDays(t *testing.T) {
	content := `
	Notifíquese la resolución el 11 de abril de 2026.
	Dentro de 3 días hábiles deberá aportar la documentación.
	`

	got := extractDocumentEvents(content)

	assertHasExtractedEvent(t, got, "notification", "2026-04-11")
	assertHasExtractedEvent(t, got, "deadline", "2026-04-15")

	deadline := mustFindExtractedEvent(t, got, "deadline", "2026-04-15")
	if deadline.DateKind != extractedDateKindRelative {
		t.Fatalf("expected relative date kind, got %q", deadline.DateKind)
	}
	if deadline.AnchorDate != "2026-04-11" {
		t.Fatalf("expected anchor date 2026-04-11, got %q", deadline.AnchorDate)
	}
	if deadline.AnchorSource != anchorSourceNotificationLine {
		t.Fatalf("expected anchor source %q, got %q", anchorSourceNotificationLine, deadline.AnchorSource)
	}
	if deadline.RelativeDays != 3 {
		t.Fatalf("expected relative days 3, got %d", deadline.RelativeDays)
	}
	if !deadline.IsBusinessDays {
		t.Fatalf("expected business days to be true")
	}
	if deadline.TriggerText == "" {
		t.Fatal("expected non-empty trigger text")
	}
}

func TestExtractDocumentEvents_DetectsRelativeDeadlineInNaturalDaysWhenExplicit(t *testing.T) {
	content := `
	Notifíquese la resolución el 11 de abril de 2026.
	Dentro de 3 días naturales deberá aportar la documentación.
	`

	got := extractDocumentEvents(content)

	assertHasExtractedEvent(t, got, "deadline", "2026-04-14")

	deadline := mustFindExtractedEvent(t, got, "deadline", "2026-04-14")
	if deadline.IsBusinessDays {
		t.Fatalf("expected natural days, got business days")
	}
	if deadline.RelativeDays != 3 {
		t.Fatalf("expected relative days 3, got %d", deadline.RelativeDays)
	}
	if deadline.AnchorSource != anchorSourceNotificationLine {
		t.Fatalf("expected anchor source %q, got %q", anchorSourceNotificationLine, deadline.AnchorSource)
	}
}
