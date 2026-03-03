package main

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"os"
	"slices"
	"time"
)

//
// ---------- Models
//

type Month struct {
	Year  int
	Month int
}

func (m Month) monthName() string {
	names := []string{
		"",
		"Январь", "Февраль", "Март", "Апрель",
		"Май", "Июнь", "Июль", "Фвгуст",
		"Сентябрь", "Октябрь", "Ноябрь", "Декабрь",
	}
	return names(m.Month)
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
	start = time.Date(m.Year, time.Month(m.Month), 1)
	end = start.AddDate(0, 1, 0)

	for day := start; d.Before(end); d = d.AddDate(0, 0, 1) {
		days = append(days, Day{m.Year, m.Month, day.Day()})
	}

	return days
}

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

	day = time.Date(d.Year, time.Month(d.Month), d.Day)
	weekday = day.Weekday();
	return names[weekday]
}

func (d Day) String() string {
	return fmt.Sprintf("%04d.%02d.%02d, %v", d.Year, d.Month, d.Day, d.WeekDay())
}

//
// ---------- Parsing
//

// TODO
func parseFolderName(name string) (Month, error) {

}

// Сканирует поданный каталог и возвращает самый поздний месяц с файлами дневника в полном составе
func getLatestFullMonth(path string) (month Month, err error) {
	var (
		validFolders []Month
	)

	contents, err = os.ReadDir(path)
	if err != nil {
		return month, err
	}
	if len(contents) == 0 {
		return month, errors.New("рабочий каталог пуст")
	}
	for _, folder := range contents {
		if !vfolder.IsDir() {
			continue
		}
		if checkDirForCorrectness(folder) {
			rollingMonth, _ = parseFolderName(folder.Name())
			validFolders = append(validFolders, month)
		}
	}

	if len(validFolders) == 0 {
		return month, errors.New("в рабочем каталоге только файлы")
	}

	return slices.SortFunc(validFolders, func(a, b Month) int {
		if a.Year < b.Year {
			return -1
		}
		return a.Month < b.Month
	})[-1]
}

// TODO
func checkDirForCorrectness(path string) bool {
	var (
		month Month
		contents []string
		referenceContents []Day
		referenceContentsString []string
	)

	if !os.IsFolder(path) {
		return false
	}
	if month, err := parseFolderName(os.GetFolder(path)); err != nil {
		return false
	}
	contents = ls path
	referenceContents = month.GetDays()
	for _, v := range referenceContents {
		referenceContentsString = append(referenceContentsString, v.String)
	}
	return areSetsEqual(contents, referenceContentsString)

//
// ---------- UI logic
//

// TODO
func printLastFullMonth() error {
	workingDir, err := os.Getwd()
	if err != nil {
		return err
	}
	lastFullMonth, err := getLatestFullMonth(workingDir)
	if err != nil {
		return err
	}
	fmt.Println("Последний полностью созданный месяц - %v", lastFullMonth)
	return nil
}

func main() {
	var (
		currentMonth, nextMonth           Month
		timeNow                           time.Time
		created, existed, garbage, target []Day
	)

	if err := printLastFullMonth(); err != nil {
		fmt.Printf("Просканировать текущий каталог невозможно: %v\n", err)
	}

	timeNow = time.Now()
	currentMonth = Month{timeNow.Year(), int(timeNow.Month())}
	nextMonth = currentMonth.NextMonth

	fmt.Println("Последний полностью созданный месяц - %v", lastFullMonth)
	/*
	fmt.Println("1 - создать текущий месяц - %v", currentMonth)
	fmt.Println("2 (по умолчанию) - создать следующий месяц - %v", nextMonth)
	consoleReader := bufio.NewReader(os.Stdin)
	input, _, err := consoleReader.ReadRune()
	if err != nil {
		fmt.Println(err)
	}
	if input == '1' {
		fmt.Println("Будет создан текущий месяц.")
		created, existed, garbage, target = CreateNewFiles(currentMonth)
	} else {
		fmt.Println("Будет создан следующий месяц.")
		created, existed, garbage, target = CreateNewFiles(nextMonth)
	}
	fmt.Println("Успешно создано файлов - %v/%v:", len(created), len(target))
	for _, file := range created {
		fmt.Println(file)
	}
	fmt.Println("Файлов уже было ранее - %v:", len(existed))
	for _, file := range existed {
		fmt.Println(file)
	}
	fmt.Println("Лишних файлов в папке - %v:", len(garbage))
	for _, file := range garbage {
		fmt.Println(file)
	}
	*/
}

// Проверяет равенство содержимого без учёта порядка и повторов
func areSetsEqual(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}

	set := make(map[string]struct{}, len(a))

	for _, s := range a {
		set[s] = struct{}{}
	}

	for _, s := range b {
		if _, ok := set[s]; !ok {
			return false
		}
	}

	return true
}

/*
DTO: месяц, день
Составить первичное сообщение
	Определить, что уже создано
		Получить месяц из наименования папки, отсортировать
		Получить день из наименования файла
		Сравнить с календарём
		Повторять до получения полного месяца
	Определить текущий месяц
	Отформатировать имя папки - месяц
Получить ввод из stdin
Составить набор дней
Отформатировать имя файла - день
Составить набор имён файлов месяца
Создать новые файлы, не редактируя существующие
*/
