package shared

import "time"

type Timestamp time.Time

func Now() Timestamp {
	return Timestamp(time.Now())
}

func (t Timestamp) Time() time.Time {
	return time.Time(t)
}
