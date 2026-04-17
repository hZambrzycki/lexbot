package documentapp

import "testing"

func TestDeduplicateUpcomingEvents_MergesDuplicateEvents(t *testing.T) {
	input := []UpcomingEvent{
		{
			EventID:        "1",
			DocumentID:     "doc-1",
			OriginalName:   "eventos.txt",
			EventType:      "hearing",
			EventDate:      "2026-05-10",
			SourceText:     "Se señala juicio para el día 10/05/2026.",
			DaysRemaining:  30,
			Status:         "upcoming",
			Priority:       "high",
			DuplicateCount: 1,
			DocumentNames:  []string{"eventos.txt"},
			DocumentIDs:    []string{"doc-1"},
			DateKind:       "absolute",
			Computation:    "absolute date",
		},
		{
			EventID:        "2",
			DocumentID:     "doc-2",
			OriginalName:   "eventos_test.txt",
			EventType:      "hearing",
			EventDate:      "2026-05-10",
			SourceText:     "Se señala juicio para el día 10/05/2026.",
			DaysRemaining:  30,
			Status:         "upcoming",
			Priority:       "high",
			DuplicateCount: 1,
			DocumentNames:  []string{"eventos_test.txt"},
			DocumentIDs:    []string{"doc-2"},
			DateKind:       "absolute",
			Computation:    "absolute date",
		},
	}

	got := deduplicateUpcomingEvents(input)

	if len(got) != 1 {
		t.Fatalf("expected 1 deduplicated event, got %d", len(got))
	}

	if got[0].DuplicateCount != 2 {
		t.Fatalf("expected duplicate count 2, got %d", got[0].DuplicateCount)
	}

	if len(got[0].DocumentNames) != 2 {
		t.Fatalf("expected 2 document names, got %d", len(got[0].DocumentNames))
	}
}

func TestDeduplicateUpcomingEvents_KeepsDifferentTypesSeparate(t *testing.T) {
	input := []UpcomingEvent{
		{
			EventID:        "1",
			DocumentID:     "doc-1",
			OriginalName:   "a.txt",
			EventType:      "hearing",
			EventDate:      "2026-05-10",
			SourceText:     "Se señala juicio para el día 10/05/2026.",
			DaysRemaining:  30,
			Status:         "upcoming",
			Priority:       "high",
			DuplicateCount: 1,
			DocumentNames:  []string{"a.txt"},
			DocumentIDs:    []string{"doc-1"},
			DateKind:       "absolute",
			Computation:    "absolute date",
		},
		{
			EventID:        "2",
			DocumentID:     "doc-2",
			OriginalName:   "b.txt",
			EventType:      "requirement",
			EventDate:      "2026-05-10",
			SourceText:     "Requiérase a la parte demandada antes del 10/05/2026.",
			DaysRemaining:  30,
			Status:         "upcoming",
			Priority:       "medium",
			DuplicateCount: 1,
			DocumentNames:  []string{"b.txt"},
			DocumentIDs:    []string{"doc-2"},
			DateKind:       "absolute",
			Computation:    "absolute date",
		},
	}

	got := deduplicateUpcomingEvents(input)

	if len(got) != 2 {
		t.Fatalf("expected 2 separate events, got %d", len(got))
	}
}

func TestDeduplicateUpcomingEvents_KeepsDifferentRelativeSemanticsSeparate(t *testing.T) {
	input := []UpcomingEvent{
		{
			EventID:        "1",
			DocumentID:     "doc-1",
			OriginalName:   "a.txt",
			EventType:      "deadline",
			EventDate:      "2026-04-16",
			SourceText:     "En el plazo de cinco días deberá presentar escrito.",
			DaysRemaining:  0,
			Status:         "today",
			Priority:       "critical",
			DuplicateCount: 1,
			DocumentNames:  []string{"a.txt"},
			DocumentIDs:    []string{"doc-1"},
			AnchorDate:     "2026-04-11",
			DateKind:       "relative",
			AnchorSource:   "notification_line",
			RelativeDays:   5,
			IsBusinessDays: false,
			TriggerText:    "en el plazo de cinco dias",
			Computation:    "anchor + 5 natural days",
		},
		{
			EventID:        "2",
			DocumentID:     "doc-2",
			OriginalName:   "b.txt",
			EventType:      "deadline",
			EventDate:      "2026-04-16",
			SourceText:     "Se concede plazo de 3 días hábiles desde el siguiente día hábil para aportar la documentación.",
			DaysRemaining:  0,
			Status:         "today",
			Priority:       "critical",
			DuplicateCount: 1,
			DocumentNames:  []string{"b.txt"},
			DocumentIDs:    []string{"doc-2"},
			AnchorDate:     "2026-04-11",
			DateKind:       "relative",
			AnchorSource:   "notification_line",
			RelativeDays:   3,
			IsBusinessDays: true,
			TriggerText:    "plazo de 3 dias habiles desde el siguiente dia habil",
			Computation:    "anchor + next_day + 3 business days",
		},
	}

	got := deduplicateUpcomingEvents(input)

	if len(got) != 2 {
		t.Fatalf("expected 2 separate relative events, got %d", len(got))
	}
}

func TestPriorityRank(t *testing.T) {
	if priorityRank("critical") >= priorityRank("high") {
		t.Fatalf("expected critical to rank before high")
	}

	if priorityRank("high") >= priorityRank("medium") {
		t.Fatalf("expected high to rank before medium")
	}

	if priorityRank("medium") >= priorityRank("low") {
		t.Fatalf("expected medium to rank before low")
	}
}

func TestStatusRank(t *testing.T) {
	if statusRank("overdue") >= statusRank("today") {
		t.Fatalf("expected overdue to rank before today")
	}

	if statusRank("today") >= statusRank("upcoming") {
		t.Fatalf("expected today to rank before upcoming")
	}
}

func TestClassifyUpcomingPriority(t *testing.T) {
	tests := []struct {
		name          string
		eventType     string
		daysRemaining int
		want          string
	}{
		{"overdue deadline", "deadline", -1, "critical"},
		{"today hearing", "hearing", 0, "critical"},
		{"future hearing", "hearing", 10, "high"},
		{"requirement", "requirement", 5, "medium"},
		{"notification", "notification", 0, "low"},
		{"filing overdue", "filing", -2, "high"},
	}

	for _, tt := range tests {
		got := classifyUpcomingPriority(tt.eventType, tt.daysRemaining)
		if got != tt.want {
			t.Fatalf("%s: expected %q, got %q", tt.name, tt.want, got)
		}
	}
}

func TestBuildUpcomingComputation_IncludesNextDayForRelativeChain(t *testing.T) {
	got := BuildUpcomingComputation(
		"relative",
		"2026-04-11",
		5,
		false,
		"plazo de 5 dias a contar desde el dia siguiente",
	)

	want := "anchor + next_day + 5 natural days"
	if got != want {
		t.Fatalf("expected %q, got %q", want, got)
	}
}

func TestBuildUpcomingComputation_UsesBusinessDays(t *testing.T) {
	got := BuildUpcomingComputation(
		"relative",
		"2026-04-11",
		3,
		true,
		"dentro de 3 dias habiles",
	)

	want := "anchor + 3 business days"
	if got != want {
		t.Fatalf("expected %q, got %q", want, got)
	}
}

func TestBuildUpcomingComputation_UsesBusinessDaysWithNextBusinessDayTrigger(t *testing.T) {
	got := BuildUpcomingComputation(
		"relative",
		"2026-04-11",
		3,
		true,
		"plazo de 3 dias habiles desde el siguiente dia habil",
	)

	want := "anchor + next_day + 3 business days"
	if got != want {
		t.Fatalf("expected %q, got %q", want, got)
	}
}

func TestBuildUpcomingComputation_ForAbsoluteDate(t *testing.T) {
	got := BuildUpcomingComputation(
		"absolute",
		"2026-04-11",
		0,
		false,
		"",
	)

	want := "absolute date"
	if got != want {
		t.Fatalf("expected %q, got %q", want, got)
	}
}

func TestBuildUpcomingComputation_ForRelativeWithoutAnchorOrDays(t *testing.T) {
	got := BuildUpcomingComputation(
		"relative",
		"",
		0,
		false,
		"",
	)

	want := "relative date"
	if got != want {
		t.Fatalf("expected %q, got %q", want, got)
	}
}

func TestTriggerAddsExtraDay(t *testing.T) {
	tests := []struct {
		name    string
		trigger string
		want    bool
	}{
		{
			name:    "a contar desde el dia siguiente",
			trigger: "plazo de 5 dias a contar desde el dia siguiente",
			want:    true,
		},
		{
			name:    "desde el siguiente",
			trigger: "plazo de 5 dias desde el siguiente",
			want:    true,
		},
		{
			name:    "desde el siguiente dia habil",
			trigger: "plazo de 3 dias habiles desde el siguiente dia habil",
			want:    true,
		},
		{
			name:    "plain relative deadline",
			trigger: "dentro de 3 dias habiles",
			want:    false,
		},
	}

	for _, tt := range tests {
		got := triggerAddsExtraDay(tt.trigger)
		if got != tt.want {
			t.Fatalf("%s: expected %v, got %v", tt.name, tt.want, got)
		}
	}
}
