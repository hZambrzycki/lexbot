package documentapp

import "testing"

func TestDeduplicateByDateAndType_MergesDuplicateEvents(t *testing.T) {
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
		},
	}

	got := deduplicateByDateAndType(input)

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

func TestDeduplicateByDateAndType_KeepsDifferentTypesSeparate(t *testing.T) {
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
		},
	}

	got := deduplicateByDateAndType(input)

	if len(got) != 2 {
		t.Fatalf("expected 2 separate events, got %d", len(got))
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
