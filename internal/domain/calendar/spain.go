package calendar

import (
	"fmt"
	"time"
)

const (
	ScopeState  = "state"
	ScopeMadrid = "madrid"
)

func Holidays(scope string, year int) []string {
	switch scope {
	case ScopeMadrid:
		return MadridHolidays(year)
	case ScopeState:
		fallthrough
	default:
		return SpainNationalHolidays(year)
	}
}

func SpainNationalHolidays(year int) []string {
	goodFriday := easterSunday(year).AddDate(0, 0, -2)

	return []string{
		fmt.Sprintf("%d-01-01", year),
		fmt.Sprintf("%d-01-06", year),
		goodFriday.Format("2006-01-02"),
		fmt.Sprintf("%d-05-01", year),
		fmt.Sprintf("%d-08-15", year),
		fmt.Sprintf("%d-10-12", year),
		fmt.Sprintf("%d-12-08", year),
		fmt.Sprintf("%d-12-25", year),
	}
}

func MadridHolidays(year int) []string {
	return mergeHolidayLists(
		SpainNationalHolidays(year),
		madridRegionalHolidays(year),
	)
}

func madridRegionalHolidays(year int) []string {
	return []string{
		fmt.Sprintf("%d-05-02", year),
	}
}

func mergeHolidayLists(lists ...[]string) []string {
	seen := make(map[string]struct{})
	result := make([]string, 0)

	for _, list := range lists {
		for _, item := range list {
			if _, exists := seen[item]; exists {
				continue
			}
			seen[item] = struct{}{}
			result = append(result, item)
		}
	}

	return result
}

// Anonymous Gregorian algorithm
func easterSunday(year int) time.Time {
	a := year % 19
	b := year / 100
	c := year % 100
	d := b / 4
	e := b % 4
	f := (b + 8) / 25
	g := (b - f + 1) / 3
	h := (19*a + b - d - g + 15) % 30
	i := c / 4
	k := c % 4
	l := (32 + 2*e + 2*i - h - k) % 7
	m := (a + 11*h + 22*l) / 451
	month := (h + l - 7*m + 114) / 31
	day := ((h + l - 7*m + 114) % 31) + 1

	return time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)
}
