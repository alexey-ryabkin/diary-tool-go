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
	names := map[int]string{
		1:  "Январь",
		2:  "Февраль",
		3:  "Март",
		4:  "Апрель",
		5:  "Май",
		6:  "Июнь",
		7:  "Июль",
		8:  "Август",
		9:  "Сентябрь",
		10: "Октябрь",
		11: "Ноябрь",
		12: "Декабрь",
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

// Days возвращает slice дней, относящихся к месяцу m по календарю.
func (m Month) Days() []Day {
	days := make([]Day, 0, 31)
	start := time.Date(m.Year, time.Month(m.Month), 1, 0, 0, 0, 0, time.UTC)
	end := start.AddDate(0, 1, 0)

	for day := start; day.Before(end); day = day.AddDate(0, 0, 1) {
		days = append(days, Day{m.Year, m.Month, day.Day()})
	}

	return days
}

func (d Month) Before(e Month) bool {
	if d.Year != e.Year {
		return d.Year < e.Year
	}
	return d.Month < e.Month
}

func SortMonths(slice []Month) []Month {
	slices.SortFunc(slice, func(a, b Month) int {
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
