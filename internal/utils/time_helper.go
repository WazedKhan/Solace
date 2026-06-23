package utils

import "time"

func IsSameDay(t1, t2 time.Time) bool {
	return t1.Year() == t2.Year() &&
		t1.YearDay() == t2.YearDay()
}
