package calendar

import (
	"testing"
	"time"
)

func TestAddBusinessDaysWithRules_SkipsEntireAugust(t *testing.T) {
	cal := NewCalendar(nil)

	start := time.Date(2026, 7, 31, 0, 0, 0, 0, time.UTC)

	got := cal.AddBusinessDaysWithRules(start, 1, ProceduralRules{
		AugustNonBusiness: true,
	})

	want := time.Date(2026, 9, 1, 0, 0, 0, 0, time.UTC)
	if !got.Equal(want) {
		t.Fatalf("expected %v, got %v", want, got)
	}
}
