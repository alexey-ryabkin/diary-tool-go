// Пакет models реализует две основые структуры данных - месяц и день,
// включая способы определения их связи и трансформации
package models

import (
	"fmt"
	"slices"
	"time"
)

// Month описывает один месяц дневника - одну папку формата 06.01, January
type Month struct {
	Year  int
	Month int
}

func (m Month) monthName() string {
	names := []string{
		"",
		"Январь", "Февраль", "Март", "Апрель",
		"Май", "Июнь", "Июль", "Август",
		"Сентябрь", "Октябрь", "Ноябрь", "Декабрь",
	}
	return names[m.Month]
}

func (m Month) String() string {
	return fmt.Sprintf("%02d.%02d, %v", m.Year%100, m.Month, m.monthName())
}

func (m Month) Next() Month {
	y := m.Year
	mo := m.Month + 1

	if mo > 12 {
		mo = 1
		y++
	}

	return Month{Year: y, Month: mo}
}

func (m Month) Days() []Day {
	days := make([]Day, 0, 31)
	start := time.Date(m.Year, time.Month(m.Month), 1, 0, 0, 0, 0, time.UTC)
	end := start.AddDate(0, 1, 0)

	for day := start; day.Before(end); day = day.AddDate(0, 0, 1) {
		days = append(days, Day{m.Year, m.Month, day.Day()})
	}

	return days
}

func SortMonths(slice []Month) []Month {
	slices.SortFunc(slice, func(a, b Month) int {
		if a.Year < b.Year {
			return -1
		}
		return a.Month - b.Month
	})
	return slice
}