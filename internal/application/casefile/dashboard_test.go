package casefileapp

import (
	documentapp "lexbox/internal/application/document"
	"lexbox/internal/domain/document"
	"testing"
)

func TestBuildDashboardAlert_IgnoresResolvedEventsForTopAlert(t *testing.T) {
	result := CaseFileDashboardResult{
		ResolvedCount: 1,
	}

	activeEvents := []documentapp.UpcomingEvent{
		{
			EventID:       "evt-pending-1",
			EventType:     "deadline",
			EventDate:     "2026-04-16",
			Status:        "overdue",
			DaysRemaining: -7,
			Priority:      "critical",
			ReviewStatus:  "pending",
			DocumentNames: []string{"doc-pending.txt"},
		},
	}

	needsAttention, topAlert := buildDashboardAlert(result, activeEvents)

	if !needsAttention {
		t.Fatal("expected needsAttention=true")
	}

	want := "critical deadline on 2026-04-16 (7 days ago)"
	if topAlert != want {
		t.Fatalf("expected top alert %q, got %q", want, topAlert)
	}
}

func TestBuildDashboardAlert_WhenOnlyResolvedEventsRemain_ReturnsNoImmediateAlert(t *testing.T) {
	result := CaseFileDashboardResult{
		ResolvedCount: 1,
	}

	activeEvents := []documentapp.UpcomingEvent{}

	needsAttention, topAlert := buildDashboardAlert(result, activeEvents)

	if needsAttention {
		t.Fatal("expected needsAttention=false")
	}

	want := "all detected events are resolved"
	if topAlert != want {
		t.Fatalf("expected top alert %q, got %q", want, topAlert)
	}
}

func TestSelectTopEvent_PrefersPendingOverReviewed(t *testing.T) {
	events := []documentapp.UpcomingEvent{
		{
			EventID:       "evt-reviewed",
			EventType:     "deadline",
			EventDate:     "2026-04-16",
			Status:        "overdue",
			DaysRemaining: -7,
			Priority:      "critical",
			ReviewStatus:  "reviewed",
			DocumentNames: []string{"reviewed.txt"},
		},
		{
			EventID:       "evt-pending",
			EventType:     "deadline",
			EventDate:     "2026-04-16",
			Status:        "overdue",
			DaysRemaining: -7,
			Priority:      "critical",
			ReviewStatus:  "pending",
			DocumentNames: []string{"pending.txt"},
		},
	}

	best := selectTopEvent(events)
	if best == nil {
		t.Fatal("expected best event")
	}

	if best.EventID != "evt-pending" {
		t.Fatalf("expected pending event to win, got %q", best.EventID)
	}
}

func TestBuildDashboardRecommendedAction_PrefixesReviewedEventsWithRecheck(t *testing.T) {
	events := []documentapp.UpcomingEvent{
		{
			EventID:       "evt-reviewed",
			EventType:     "deadline",
			EventDate:     "2026-04-20",
			Status:        "upcoming",
			DaysRemaining: 2,
			Priority:      "high",
			ReviewStatus:  "reviewed",
			DocumentNames: []string{"demanda.txt"},
		},
	}

	got := buildDashboardRecommendedAction(events)
	want := "re-check review upcoming deadline (demanda.txt) for 2026-04-20"

	if got != want {
		t.Fatalf("expected %q, got %q", want, got)
	}
}

func TestBuildDashboardProceduralHint_UsesMostRelevantActiveEvent(t *testing.T) {
	events := []documentapp.UpcomingEvent{
		{
			EventID:       "evt-notification",
			EventType:     "notification",
			EventDate:     "2026-04-11",
			Status:        "overdue",
			DaysRemaining: -12,
			Priority:      "low",
			ReviewStatus:  "pending",
		},
		{
			EventID:       "evt-deadline",
			EventType:     "deadline",
			EventDate:     "2026-04-16",
			Status:        "overdue",
			DaysRemaining: -7,
			Priority:      "critical",
			ReviewStatus:  "pending",
		},
	}

	got := buildDashboardProceduralHint(events)
	want := "possible deadline breach"

	if got != want {
		t.Fatalf("expected %q, got %q", want, got)
	}
}

func TestFilterActiveDocumentEvents_RemovesResolvedOnly(t *testing.T) {
	events := []document.Event{
		{ID: "1", ReviewStatus: "pending"},
		{ID: "2", ReviewStatus: "reviewed"},
		{ID: "3", ReviewStatus: "resolved"},
		{ID: "4", ReviewStatus: ""},
	}

	active := filterActiveDocumentEvents(events)

	if len(active) != 3 {
		t.Fatalf("expected 3 active events, got %d", len(active))
	}

	if active[0].ID != "1" {
		t.Fatalf("expected first active id=%q, got %q", "1", active[0].ID)
	}
	if active[1].ID != "2" {
		t.Fatalf("expected second active id=%q, got %q", "2", active[1].ID)
	}
	if active[2].ID != "4" {
		t.Fatalf("expected third active id=%q, got %q", "4", active[2].ID)
	}
}

func TestNormalizeReviewStatus_DefaultsToPending(t *testing.T) {
	got := normalizeReviewStatus("   ")
	want := "pending"

	if got != want {
		t.Fatalf("expected %q, got %q", want, got)
	}
}

func TestReviewStatusRank_OrderIsPendingReviewedResolved(t *testing.T) {
	if reviewStatusRank("pending") >= reviewStatusRank("reviewed") {
		t.Fatal("expected pending to rank before reviewed")
	}

	if reviewStatusRank("reviewed") >= reviewStatusRank("resolved") {
		t.Fatal("expected reviewed to rank before resolved")
	}
}
