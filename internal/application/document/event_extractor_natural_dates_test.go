package documentapp

import "testing"

func TestExtractDocumentEvents_DetectsSpanishTextualDates(t *testing.T) {
	content := `
	Se señala juicio para el día 10 de mayo de 2026.
	La comparecencia tendrá lugar el 15 de mayo de 2026.
	Notifíquese la resolución el 11 de abril de 2026.
	`

	got := extractDocumentEvents(content, DefaultEventComputationConfig())
	if len(got) != 3 {
		t.Fatalf("expected 3 events, got %d", len(got))
	}

	assertHasExtractedEvent(t, got, "hearing", "2026-05-10")
	assertHasExtractedEvent(t, got, "appearance", "2026-05-15")
	assertHasExtractedEvent(t, got, "notification", "2026-04-11")

	hearing := mustFindExtractedEvent(t, got, "hearing", "2026-05-10")
	if hearing.DateKind != extractedDateKindAbsolute {
		t.Fatalf("expected absolute date kind, got %q", hearing.DateKind)
	}
	if hearing.AnchorDate != "2026-05-10" {
		t.Fatalf("expected anchor date 2026-05-10, got %q", hearing.AnchorDate)
	}
	if hearing.AnchorSource != anchorSourceInline {
		t.Fatalf("expected anchor source %q, got %q", anchorSourceInline, hearing.AnchorSource)
	}
	if hearing.RelativeDays != 0 {
		t.Fatalf("expected relative days 0, got %d", hearing.RelativeDays)
	}
	if hearing.IsBusinessDays {
		t.Fatalf("expected non-business absolute date")
	}
	if hearing.TriggerText != "" {
		t.Fatalf("expected empty trigger text, got %q", hearing.TriggerText)
	}
}
