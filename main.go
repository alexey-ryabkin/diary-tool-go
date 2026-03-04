package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strconv"
	"strings"
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

//
// ---------- Parsing
//

// 24.04, Апрель
func parseFolderName(name string) (Month, error) {
	parts := strings.Split(name, ",")
	if len(parts) < 2 {
		return Month{}, errors.New("наименование папки не содержит запятой")
	}

	parts = strings.Split(parts[0], ".")
	if len(parts) < 2 {
		return Month{}, errors.New("год и месяц не разделены точкой")
	}

	year, err := strconv.Atoi(parts[0])
	if err != nil {
		return Month{}, err
	}
	if year < 0 || year > 99 {
		return Month{}, errors.New("год указан не двумя цифрами")
	}
	year = 2000 + year

	month, err := strconv.Atoi(parts[1])
	if err != nil {
		return Month{}, err
	}
	if month < 1 || month > 12 {
		return Month{}, errors.New("месяц не может быть меньше одного или больше двенадцати")
	}

	validMonth := Month{year, month}
	canon := validMonth.String()
	if name != canon {
		errorMsg := fmt.Sprint("формат ", name, " не соответствует каноническому ", canon)
		return Month{}, errors.New(errorMsg)
	}

	return validMonth, nil
}

// Сканирует поданный каталог и возвращает самый поздний месяц с файлами дневника в полном составе
func getLatestFullMonth(path string) (month Month, err error) {
	var (
		validFolders []Month
	)

	contents, err := os.ReadDir(path)
	if err != nil {
		return month, err
	}
	if len(contents) == 0 {
		return month, errors.New("рабочий каталог пуст")
	}
	for _, folder := range contents {
		if !folder.IsDir() {
			continue
		}
		correct, _ := checkDirForCorrectness(filepath.Join(path, folder.Name()))
		if correct {
			rollingMonth, _ := parseFolderName(folder.Name())
			validFolders = append(validFolders, rollingMonth)
		}
	}

	if len(validFolders) == 0 {
		return month, errors.New("в рабочем каталоге нет ни одной папки дневника")
	}

	slices.SortFunc(validFolders, func(a, b Month) int {
		if a.Year < b.Year {
			return -1
		}
		return a.Month - b.Month
	})

	return validFolders[len(validFolders)-1], nil
}

func checkDirForCorrectness(path string) (bool, error) {
	var (
		month                   Month
		contentsStr             []string
		referenceContents       []Day
		referenceContentsString []string
	)

	fileInfo, err := os.Stat(path)
	if err != nil {
		return false, errors.New("папки не существует")
	}
	if !fileInfo.IsDir() {
		return false, errors.New("по указанному пути находится файл")
	}
	if month, err = parseFolderName(filepath.Base(path)); err != nil {
		return false, err
	}

	contents, err := os.ReadDir(path)
	if err != nil {
		return false, err
	}

	referenceContents = month.Days()
	for _, v := range referenceContents {
		referenceContentsString = append(referenceContentsString, v.String())
	}

	contentsStr = make([]string, 0, len(contents))
	for _, e := range contents {
		contentsStr = append(contentsStr, e.Name())
	}
	return areSetsEqual(contentsStr, referenceContentsString), nil
}

//
// ---------- UI logic
//

func printLastFullMonth() error {
	workingDir, err := os.Getwd()
	if err != nil {
		return err
	}
	lastFullMonth, err := getLatestFullMonth(workingDir)
	if err != nil {
		return err
	}
	fmt.Printf("Последний полностью созданный месяц - %v\n", lastFullMonth)
	return nil
}

func main() {
	var (
		currentMonth, nextMonth           Month
		timeNow                           time.Time
		created, existed, garbage, target []Day
	)

	if err := printLastFullMonth(); err != nil {
		fmt.Printf("Определить последний созданный месяц невозможно: %v\n", err)
	}

	timeNow = time.Now()
	currentMonth = Month{timeNow.Year(), int(timeNow.Month())}
	nextMonth = currentMonth.Next()
	fmt.Printf("1 - создать текущий месяц - %v\n", currentMonth)
	fmt.Printf("2 (по умолчанию) - создать следующий месяц - %v\n", nextMonth)
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
	fmt.Printf("Успешно создано файлов - %v/%v:\n", len(created), len(target))
	for _, file := range created {
		fmt.Println(file)
	}
	fmt.Printf("Файлов уже было ранее - %v:\n", len(existed))
	for _, file := range existed {
		fmt.Println(file)
	}
	fmt.Printf("Лишних файлов в папке - %v:\n", len(garbage))
	for _, file := range garbage {
		fmt.Println(file)
	}
}

func CreateNewFiles(currentMonth Month) ([]Day, []Day, []Day, []Day) {
	fmt.Print("unimplemented")
	a := make([]Day, 0)
	return a, a, a, a
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
