package calendar

import "testing"

func TestSpainNationalHolidays_UsesDynamicGoodFriday_2026(t *testing.T) {
	holidays := SpainNationalHolidays(2026)

	want := "2026-04-03"
	if !containsHoliday(holidays, want) {
		t.Fatalf("expected Good Friday %q in national holidays", want)
	}
}

func TestSpainNationalHolidays_UsesDynamicGoodFriday_2027(t *testing.T) {
	holidays := SpainNationalHolidays(2027)

	want := "2027-03-26"
	if !containsHoliday(holidays, want) {
		t.Fatalf("expected Good Friday %q in national holidays", want)
	}
}

func TestMadridHolidays_IncludeRegionalHoliday(t *testing.T) {
	holidays := MadridHolidays(2026)

	want := "2026-05-02"
	if !containsHoliday(holidays, want) {
		t.Fatalf("expected Madrid regional holiday %q", want)
	}
}

func TestMadridHolidays_KeepNationalHoliday(t *testing.T) {
	holidays := MadridHolidays(2026)

	want := "2026-12-25"
	if !containsHoliday(holidays, want) {
		t.Fatalf("expected national holiday %q", want)
	}
}

func TestMergeHolidayLists_Deduplicates(t *testing.T) {
	got := mergeHolidayLists(
		[]string{"2026-05-01", "2026-05-02"},
		[]string{"2026-05-02", "2026-12-25"},
	)

	if len(got) != 3 {
		t.Fatalf("expected 3 merged holidays, got %d", len(got))
	}
}

func containsHoliday(items []string, want string) bool {
	for _, item := range items {
		if item == want {
			return true
		}
	}
	return false
}
