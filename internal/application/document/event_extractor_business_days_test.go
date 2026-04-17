package documentapp

import "testing"

func TestExtractDocumentEvents_DetectsRelativeDeadlineInBusinessDays(t *testing.T) {
	content := `
	Notifíquese la resolución el 13 de abril de 2026.
	Dentro de 3 días hábiles deberá aportar la documentación.
	`

	got := extractDocumentEvents(content)

	assertHasExtractedEvent(t, got, "notification", "2026-04-13")
	assertHasExtractedEvent(t, got, "deadline", "2026-04-16")

	deadline := mustFindExtractedEvent(t, got, "deadline", "2026-04-16")
	assertRelativeDeadline(t, deadline, relativeDeadlineExpectation{
		date:           "2026-04-16",
		anchorDate:     "2026-04-13",
		anchorSource:   anchorSourceNotificationLine,
		relativeDays:   3,
		businessDays:   true,
		addExtraDay:    false,
		requireTrigger: true,
	})
}

func TestExtractDocumentEvents_DetectsRelativeDeadlineInNaturalDaysWhenExplicit(t *testing.T) {
	content := `
	Notifíquese la resolución el 11 de abril de 2026.
	Dentro de 3 días naturales deberá aportar la documentación.
	`

	got := extractDocumentEvents(content)

	assertHasExtractedEvent(t, got, "deadline", "2026-04-14")

	deadline := mustFindExtractedEvent(t, got, "deadline", "2026-04-14")
	assertRelativeDeadline(t, deadline, relativeDeadlineExpectation{
		date:           "2026-04-14",
		anchorDate:     "2026-04-11",
		anchorSource:   anchorSourceNotificationLine,
		relativeDays:   3,
		businessDays:   false,
		addExtraDay:    false,
		requireTrigger: true,
	})
}

func TestExtractDocumentEvents_DetectsBusinessDaysCountingFromNextDay(t *testing.T) {
	content := `
	Notifíquese la resolución el 13/04/2026.
	Dentro de 3 días hábiles a contar desde el día siguiente deberá aportar la documentación.
	`

	got := extractDocumentEvents(content)

	assertHasExtractedEvent(t, got, "notification", "2026-04-13")
	assertHasExtractedEvent(t, got, "deadline", "2026-04-17")

	deadline := mustFindExtractedEvent(t, got, "deadline", "2026-04-17")
	assertRelativeDeadline(t, deadline, relativeDeadlineExpectation{
		date:           "2026-04-17",
		anchorDate:     "2026-04-13",
		anchorSource:   anchorSourceNotificationLine,
		relativeDays:   3,
		businessDays:   true,
		addExtraDay:    true,
		requireTrigger: true,
	})
}

func TestExtractDocumentEvents_DetectsNaturalDaysCountingFromNextDay(t *testing.T) {
	content := `
	Notifíquese la resolución el 11/04/2026.
	Dentro de 3 días naturales a contar desde el día siguiente deberá aportar la documentación.
	`

	got := extractDocumentEvents(content)

	assertHasExtractedEvent(t, got, "notification", "2026-04-11")
	assertHasExtractedEvent(t, got, "deadline", "2026-04-15")

	deadline := mustFindExtractedEvent(t, got, "deadline", "2026-04-15")
	assertRelativeDeadline(t, deadline, relativeDeadlineExpectation{
		date:           "2026-04-15",
		anchorDate:     "2026-04-11",
		anchorSource:   anchorSourceNotificationLine,
		relativeDays:   3,
		businessDays:   false,
		addExtraDay:    true,
		requireTrigger: true,
	})
}

func TestExtractDocumentEvents_DetectsBusinessDaysFromNextBusinessDay(t *testing.T) {
	content := `
	Notifíquese la resolución el 13/04/2026.
	Se concede plazo de 3 días hábiles desde el siguiente día hábil para aportar la documentación.
	`

	got := extractDocumentEvents(content)

	assertHasExtractedEvent(t, got, "notification", "2026-04-13")
	assertHasExtractedEvent(t, got, "deadline", "2026-04-17")

	deadline := mustFindExtractedEvent(t, got, "deadline", "2026-04-17")
	assertRelativeDeadline(t, deadline, relativeDeadlineExpectation{
		date:           "2026-04-17",
		anchorDate:     "2026-04-13",
		anchorSource:   anchorSourceNotificationLine,
		relativeDays:   3,
		businessDays:   true,
		addExtraDay:    true,
		requireTrigger: true,
	})
}

func TestExtractDocumentEvents_DetectsBusinessDaysFromNextBusinessDay_Apartir(t *testing.T) {
	content := `
	Notifíquese la resolución el 13/04/2026.
	Se concede plazo de 3 días hábiles a partir del siguiente día hábil para aportar la documentación.
	`

	got := extractDocumentEvents(content)

	assertHasExtractedEvent(t, got, "notification", "2026-04-13")
	assertHasExtractedEvent(t, got, "deadline", "2026-04-17")

	deadline := mustFindExtractedEvent(t, got, "deadline", "2026-04-17")
	assertRelativeDeadline(t, deadline, relativeDeadlineExpectation{
		date:           "2026-04-17",
		anchorDate:     "2026-04-13",
		anchorSource:   anchorSourceNotificationLine,
		relativeDays:   3,
		businessDays:   true,
		addExtraDay:    true,
		requireTrigger: true,
	})
}

func TestExtractDocumentEvents_TriggerPreservesFullNextBusinessDayPhrase(t *testing.T) {
	content := `
	Notifíquese la resolución el 13/04/2026.
	Se concede plazo de 3 días hábiles desde el siguiente día hábil para aportar la documentación.
	`

	got := extractDocumentEvents(content)

	deadline := mustFindExtractedEvent(t, got, "deadline", "2026-04-17")

	want := "plazo de 3 dias habiles desde el siguiente dia habil"
	if deadline.TriggerText != want {
		t.Fatalf("expected trigger %q, got %q", want, deadline.TriggerText)
	}
}

func TestExtractDocumentEvents_BusinessDaysSkipConfiguredHoliday(t *testing.T) {
	content := `
	Notifíquese la resolución el 02/04/2026.
	Dentro de 1 día hábil deberá aportar la documentación.
	`

	got := extractDocumentEvents(content)

	assertHasExtractedEvent(t, got, "notification", "2026-04-02")
	assertHasExtractedEvent(t, got, "deadline", "2026-04-06")

	deadline := mustFindExtractedEvent(t, got, "deadline", "2026-04-06")
	assertRelativeDeadline(t, deadline, relativeDeadlineExpectation{
		date:           "2026-04-06",
		anchorDate:     "2026-04-02",
		anchorSource:   anchorSourceNotificationLine,
		relativeDays:   1,
		businessDays:   true,
		addExtraDay:    false,
		requireTrigger: true,
	})
}

func TestExtractDocumentEvents_DetectsSingleBusinessDay(t *testing.T) {
	content := `
	Notifíquese la resolución el 02/04/2026.
	Dentro de 1 día hábil deberá aportar la documentación.
	`

	got := extractDocumentEvents(content)

	assertHasExtractedEvent(t, got, "notification", "2026-04-02")
	assertHasExtractedEvent(t, got, "deadline", "2026-04-06")
}

func TestExtractDocumentEvents_BusinessDaysUseMadridScope(t *testing.T) {
	content := `
	Notifíquese la resolución el 01/05/2026.
	Dentro de 1 día hábil deberá aportar la documentación.
	`

	got := extractDocumentEvents(content)

	assertHasExtractedEvent(t, got, "notification", "2026-05-01")
	assertHasExtractedEvent(t, got, "deadline", "2026-05-04")

	deadline := mustFindExtractedEvent(t, got, "deadline", "2026-05-04")
	assertRelativeDeadline(t, deadline, relativeDeadlineExpectation{
		date:           "2026-05-04",
		anchorDate:     "2026-05-01",
		anchorSource:   anchorSourceNotificationLine,
		relativeDays:   1,
		businessDays:   true,
		addExtraDay:    false,
		requireTrigger: true,
	})
}

func TestExtractDocumentEvents_BusinessDaysSkipAugustWhenProceduralRuleEnabled(t *testing.T) {
	content := `
	Notifíquese la resolución el 31/07/2026.
	Dentro de 1 día hábil deberá aportar la documentación.
	`

	got := extractDocumentEvents(content)

	assertHasExtractedEvent(t, got, "notification", "2026-07-31")
	assertHasExtractedEvent(t, got, "deadline", "2026-09-01")

	deadline := mustFindExtractedEvent(t, got, "deadline", "2026-09-01")
	assertRelativeDeadline(t, deadline, relativeDeadlineExpectation{
		date:           "2026-09-01",
		anchorDate:     "2026-07-31",
		anchorSource:   anchorSourceNotificationLine,
		relativeDays:   1,
		businessDays:   true,
		addExtraDay:    false,
		requireTrigger: true,
	})
}
