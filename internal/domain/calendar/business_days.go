package calendar

import "time"

type Calendar struct {
	Holidays map[string]struct{}
}

func NewCalendar(holidays []string) Calendar {
	items := make(map[string]struct{}, len(holidays))
	for _, holiday := range holidays {
		items[holiday] = struct{}{}
	}

	return Calendar{
		Holidays: items,
	}
}

func (c Calendar) IsBusinessDay(t time.Time) bool {
	return c.IsBusinessDayWithRules(t, DefaultProceduralRules())
}

func (c Calendar) IsBusinessDayWithRules(t time.Time, rules ProceduralRules) bool {
	switch t.Weekday() {
	case time.Saturday, time.Sunday:
		return false
	}

	if rules.AugustNonBusiness && t.Month() == time.August {
		return false
	}

	_, isHoliday := c.Holidays[t.Format("2006-01-02")]
	return !isHoliday
}

func (c Calendar) AddBusinessDays(start time.Time, days int) time.Time {
	return c.AddBusinessDaysWithRules(start, days, DefaultProceduralRules())
}

func (c Calendar) AddBusinessDaysWithRules(start time.Time, days int, rules ProceduralRules) time.Time {
	if days <= 0 {
		return start
	}

	current := start
	added := 0

	for added < days {
		current = current.AddDate(0, 0, 1)
		if c.IsBusinessDayWithRules(current, rules) {
			added++
		}
	}

	return current
}
