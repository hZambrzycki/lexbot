package documentapp

import "time"

func addBusinessDays(base time.Time, days int) time.Time {
	if days <= 0 {
		return base
	}

	current := base
	added := 0

	for added < days {
		current = current.AddDate(0, 0, 1)
		if isBusinessDay(current) {
			added++
		}
	}

	return current
}

func isBusinessDay(t time.Time) bool {
	switch t.Weekday() {
	case time.Saturday, time.Sunday:
		return false
	default:
		return true
	}
}
