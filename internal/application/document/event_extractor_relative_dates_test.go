package documentapp

import "testing"

func TestExtractDocumentEvents_DetectsRelativeDeadlineFromNotificationDate(t *testing.T) {
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
	assertRelativeDeadline(t, deadline, relativeDeadlineExpectation{
		date:           "2026-04-16",
		anchorDate:     "2026-04-11",
		anchorSource:   anchorSourceNotificationLine,
		relativeDays:   5,
		businessDays:   false,
		addExtraDay:    false,
		requireTrigger: true,
	})
}

func TestExtractDocumentEvents_DetectsRelativeDeadlineWithTextualNumber(t *testing.T) {
	content := `
	Notifíquese la resolución el 11 de abril de 2026.
	En el plazo de cinco días deberá presentar escrito.
	`

	got := extractDocumentEvents(content)

	assertHasExtractedEvent(t, got, "deadline", "2026-04-16")

	deadline := mustFindExtractedEvent(t, got, "deadline", "2026-04-16")
	assertRelativeDeadline(t, deadline, relativeDeadlineExpectation{
		date:           "2026-04-16",
		anchorDate:     "2026-04-11",
		anchorSource:   anchorSourceNotificationLine,
		relativeDays:   5,
		businessDays:   false,
		addExtraDay:    false,
		requireTrigger: true,
	})
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

func TestExtractDocumentEvents_PrioritizesNotificationDateOverPreviousAbsoluteDate(t *testing.T) {
	content := `
	DILIGENCIA DE ORDENACIÓN de 09/04/2026.
	Notifíquese la presente resolución a las partes el 11/04/2026.
	Se concede plazo de 5 días para formular alegaciones.
	`

	got := extractDocumentEvents(content)

	assertHasExtractedEvent(t, got, "notification", "2026-04-11")
	assertHasExtractedEvent(t, got, "deadline", "2026-04-16")

	deadline := mustFindExtractedEvent(t, got, "deadline", "2026-04-16")
	assertRelativeDeadline(t, deadline, relativeDeadlineExpectation{
		date:           "2026-04-16",
		anchorDate:     "2026-04-11",
		anchorSource:   anchorSourceNotificationLine,
		relativeDays:   5,
		businessDays:   false,
		addExtraDay:    false,
		requireTrigger: true,
	})
}

func TestExtractDocumentEvents_DetectsNextDayAfterNotification(t *testing.T) {
	content := `
	Notifíquese la resolución el 11/04/2026.
	Al día siguiente de la notificación deberá aportar la documentación.
	`

	got := extractDocumentEvents(content)

	assertHasExtractedEvent(t, got, "notification", "2026-04-11")
	assertHasExtractedEvent(t, got, "deadline", "2026-04-12")

	deadline := mustFindExtractedEvent(t, got, "deadline", "2026-04-12")
	assertRelativeDeadline(t, deadline, relativeDeadlineExpectation{
		date:           "2026-04-12",
		anchorDate:     "2026-04-11",
		anchorSource:   anchorSourceNotificationLine,
		relativeDays:   1,
		businessDays:   false,
		addExtraDay:    false,
		requireTrigger: true,
	})
}

func TestExtractDocumentEvents_DetectsRelativeDeadlineCountingFromNextDay(t *testing.T) {
	content := `
	Notifíquese la resolución el 11/04/2026.
	Se concede plazo de 5 días a contar desde el día siguiente para formular alegaciones.
	`

	got := extractDocumentEvents(content)

	assertHasExtractedEvent(t, got, "notification", "2026-04-11")
	assertHasExtractedEvent(t, got, "deadline", "2026-04-17")

	deadline := mustFindExtractedEvent(t, got, "deadline", "2026-04-17")
	assertRelativeDeadline(t, deadline, relativeDeadlineExpectation{
		date:           "2026-04-17",
		anchorDate:     "2026-04-11",
		anchorSource:   anchorSourceNotificationLine,
		relativeDays:   5,
		businessDays:   false,
		addExtraDay:    true,
		requireTrigger: true,
	})
}

func TestExtractDocumentEvents_DetectsRelativeDeadlineFromNextDayVariant(t *testing.T) {
	content := `
	Notifíquese la resolución el 11/04/2026.
	Desde el día siguiente se concede plazo de 5 días para recurrir.
	`

	got := extractDocumentEvents(content)

	assertHasExtractedEvent(t, got, "notification", "2026-04-11")
	assertHasExtractedEvent(t, got, "deadline", "2026-04-17")

	deadline := mustFindExtractedEvent(t, got, "deadline", "2026-04-17")
	assertRelativeDeadline(t, deadline, relativeDeadlineExpectation{
		date:           "2026-04-17",
		anchorDate:     "2026-04-11",
		anchorSource:   anchorSourceNotificationLine,
		relativeDays:   5,
		businessDays:   false,
		addExtraDay:    true,
		requireTrigger: true,
	})
}

func TestExtractDocumentEvents_FromNotificationWithoutNextDay(t *testing.T) {
	content := `
	Notifíquese la resolución el 11/04/2026.
	Se concede plazo de 5 días desde la notificación para formular alegaciones.
	`

	got := extractDocumentEvents(content)

	assertHasExtractedEvent(t, got, "notification", "2026-04-11")
	assertHasExtractedEvent(t, got, "deadline", "2026-04-16")

	deadline := mustFindExtractedEvent(t, got, "deadline", "2026-04-16")
	assertRelativeDeadline(t, deadline, relativeDeadlineExpectation{
		date:           "2026-04-16",
		anchorDate:     "2026-04-11",
		anchorSource:   anchorSourceNotificationLine,
		relativeDays:   5,
		businessDays:   false,
		addExtraDay:    false,
		requireTrigger: true,
	})
}

func TestExtractDocumentEvents_FromNextDayVariant_ApartirDelDiaSiguiente(t *testing.T) {
	content := `
	Notifíquese la resolución el 11/04/2026.
	Se concede plazo de 5 días a partir del día siguiente para formular alegaciones.
	`

	got := extractDocumentEvents(content)

	assertHasExtractedEvent(t, got, "notification", "2026-04-11")
	assertHasExtractedEvent(t, got, "deadline", "2026-04-17")

	deadline := mustFindExtractedEvent(t, got, "deadline", "2026-04-17")
	assertRelativeDeadline(t, deadline, relativeDeadlineExpectation{
		date:           "2026-04-17",
		anchorDate:     "2026-04-11",
		anchorSource:   anchorSourceNotificationLine,
		relativeDays:   5,
		businessDays:   false,
		addExtraDay:    true,
		requireTrigger: true,
	})
}

func TestExtractDocumentEvents_FromNotificationPhrase_ViaSuNotificacion(t *testing.T) {
	content := `
	Notifíquese la resolución el 11/04/2026.
	Se concede plazo de 5 días desde su notificación para formular alegaciones.
	`

	got := extractDocumentEvents(content)

	assertHasExtractedEvent(t, got, "notification", "2026-04-11")
	assertHasExtractedEvent(t, got, "deadline", "2026-04-16")

	deadline := mustFindExtractedEvent(t, got, "deadline", "2026-04-16")
	assertRelativeDeadline(t, deadline, relativeDeadlineExpectation{
		date:           "2026-04-16",
		anchorDate:     "2026-04-11",
		anchorSource:   anchorSourceNotificationLine,
		relativeDays:   5,
		businessDays:   false,
		addExtraDay:    false,
		requireTrigger: true,
	})
}

type relativeDeadlineExpectation struct {
	date           string
	anchorDate     string
	anchorSource   string
	relativeDays   int
	businessDays   bool
	addExtraDay    bool
	requireTrigger bool
}

func assertRelativeDeadline(t *testing.T, got extractedEventCandidate, want relativeDeadlineExpectation) {
	t.Helper()

	if got.DateKind != extractedDateKindRelative {
		t.Fatalf("expected relative date kind, got %q", got.DateKind)
	}
	if got.EventDate != want.date {
		t.Fatalf("expected event date %q, got %q", want.date, got.EventDate)
	}
	if got.AnchorDate != want.anchorDate {
		t.Fatalf("expected anchor date %q, got %q", want.anchorDate, got.AnchorDate)
	}
	if got.AnchorSource != want.anchorSource {
		t.Fatalf("expected anchor source %q, got %q", want.anchorSource, got.AnchorSource)
	}
	if got.RelativeDays != want.relativeDays {
		t.Fatalf("expected relative days %d, got %d", want.relativeDays, got.RelativeDays)
	}
	if got.IsBusinessDays != want.businessDays {
		t.Fatalf("expected business days %v, got %v", want.businessDays, got.IsBusinessDays)
	}
	if got.AddExtraDay != want.addExtraDay {
		t.Fatalf("expected add extra day %v, got %v", want.addExtraDay, got.AddExtraDay)
	}
	if want.requireTrigger && got.TriggerText == "" {
		t.Fatal("expected non-empty trigger text")
	}
}
