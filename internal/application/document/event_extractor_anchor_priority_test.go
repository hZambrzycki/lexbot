package documentapp

import "testing"

func TestExtractDocumentEvents_UsesStrongProceduralAnchorWhenNotificationDoesNotExist(t *testing.T) {
	content := `
	DECRETO de fecha 09/04/2026.
	Se concede plazo de 5 días para formular alegaciones.
	`

	got := extractDocumentEvents(content, DefaultEventComputationConfig())
	assertHasExtractedEvent(t, got, "deadline", "2026-04-14")

	deadline := mustFindExtractedEvent(t, got, "deadline", "2026-04-14")
	if deadline.AnchorDate != "2026-04-09" {
		t.Fatalf("expected anchor date 2026-04-09, got %q", deadline.AnchorDate)
	}
	if deadline.AnchorSource != anchorSourceProceduralAnchorLine {
		t.Fatalf("expected anchor source %q, got %q", anchorSourceProceduralAnchorLine, deadline.AnchorSource)
	}
	if deadline.DateKind != extractedDateKindRelative {
		t.Fatalf("expected date kind %q, got %q", extractedDateKindRelative, deadline.DateKind)
	}
}

func TestExtractDocumentEvents_PrioritizesNotificationOverStrongProceduralAnchor(t *testing.T) {
	content := `
	AUTO de fecha 09/04/2026.
	Notifíquese la presente resolución a las partes el 11/04/2026.
	Se concede plazo de 5 días a contar desde la notificación.
	`

	got := extractDocumentEvents(content, DefaultEventComputationConfig())

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
