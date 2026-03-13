package fs

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/alexey-ryabkin/diary-tool-go/internal/models"
	"github.com/alexey-ryabkin/diary-tool-go/internal/parse"
	"github.com/alexey-ryabkin/diary-tool-go/internal/utils"
)

// GetLatestFullMonth сканирует поданный каталог 
// и возвращает самый поздний месяц с файлами дневника в полном составе.
func GetLatestFullMonth(path string) (month models.Month, err error) {
	var (
		validFolders []models.Month
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
		correct, _ := CheckDirForCorrectness(filepath.Join(path, folder.Name()))
		if correct {
			rollingMonth, _ := parse.FolderName(folder.Name())
			validFolders = append(validFolders, rollingMonth)
		}
	}

	if len(validFolders) == 0 {
		return month, errors.New("в рабочем каталоге нет ни одной папки дневника")
	}

	validFolders = models.SortMonths(validFolders)

	return validFolders[len(validFolders)-1], nil
}

// CheckDirForCorrectness проверяет, соответствует ли папка по заданному пути идеалу папки дневника.
func CheckDirForCorrectness(path string) (bool, error) {
	var (
		month                   models.Month
		contentsStr             []string
		referenceContents       []models.Day
		referenceContentsString []string
	)

	fileInfo, err := os.Stat(path)
	if err != nil {
		return false, errors.New("папки не существует")
	}
	if !fileInfo.IsDir() {
		return false, errors.New("по указанному пути находится файл")
	}
	if month, err = parse.FolderName(filepath.Base(path)); err != nil {
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
	return utils.AreSetsEqual(contentsStr, referenceContentsString), nil
}

type ScanResult struct {
	// Уже существующие дни
	Existed []models.Day
	// Недостающие дни
	Missing []models.Day
	// Лишние файлы
	Garbage []string
}

// FolderContents сканирует содержимое папки дневника за месяц.
// Она:
//   - определяет месяц по имени папки
//   - проверяет какие файлы дней уже существуют
//   - определяет недостающие дни
//   - собирает список посторонних файлов
func FolderContents(path string) (ScanResult, error) {
	var (
		res ScanResult
		month models.Month
		days []models.Day
	)

	// Получаем имя папки (например "2024-03")
	folderName := filepath.Base(path)

	// Парсим месяц
	month, err := parse.FolderName(folderName)
	if err != nil {
		return res, err
	}

	// Если папка не существует, то в ней не хватает всех файлов
	res.Missing = month.Days()

	// Проверяем что путь указывает на директорию
	fileInfo, err := os.Stat(path)
	if err != nil {
		return res, err
	}
	if !fileInfo.IsDir() {
		return res, errors.New("по указанному пути находится файл")
	}

	// Читаем содержимое директории
	contents, err := os.ReadDir(path)
	if err != nil {
		return res, err
	}

	// Получаем ожидаемые дни месяца
	days = month.Days()

	// Карта ожидаемых файлов
	targetMap := make(map[string]models.Day, len(days))
	for _, d := range days {
		targetMap[d.String()] = d
	}

	res.Existed = make([]models.Day, 0, len(contents))
	res.Garbage = make([]string, 0, len(contents))

	// Фильтрация существующих файлов
	for _, entry := range contents {
		name := entry.Name()

		if day, ok := targetMap[name]; ok {
			res.Existed = append(res.Existed, day)
			delete(targetMap, name)
			continue
		}

		res.Garbage = append(res.Garbage, name)
	}

	// Оставшиеся элементы карты — отсутствующие дни
	res.Missing = make([]models.Day, 0, len(targetMap))
	for _, day := range targetMap {
		res.Missing = append(res.Missing, day)
	}

	return res, nil
}
