package models

import (
	"fmt"
	"slices"
	"time"
)

// Day описывает один день дневника - один файл формата 2006.01.02, Monday.txt
type Day struct {
	Year  int
	Month int
	Day   int
}

func (d Day) WeekDay() string {
	names := []string{
		"воскресенье",
		"понедельник",
		"вторник",
		"среда",
		"четверг",
		"пятница",
		"суббота",
	}

	day := time.Date(d.Year, time.Month(d.Month), d.Day, 0, 0, 0, 0, time.UTC)
	weekday := day.Weekday()
	return names[weekday]
}

func (d Day) String() string {
	return fmt.Sprintf("%04d.%02d.%02d, %v.txt", d.Year, d.Month, d.Day, d.WeekDay())
}

func (d Day) Before(e Day) bool {
	if d.Year != e.Year {
		return d.Year < e.Year
	}
	if d.Month != e.Month {
		return d.Month < e.Month
	}
	return d.Day < e.Day
}

func SortDays(slice []Day) []Day {
	slices.SortFunc(slice, func(a, b Day) int {
		if a.Before(b) {
			return -1
		}
		if b.Before(a) {
			return 1
		}
		return 0
	})
	return slice
}
