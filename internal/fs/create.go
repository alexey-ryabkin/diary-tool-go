package fs

import (
	"fmt"
	"os"
	"path/filepath"
	"slices"

	"github.com/alexey-ryabkin/diary-tool-go/internal/models"
)

type Result struct {
	Created []models.Day
	Existed []models.Day
	Garbage []string
	Target  []models.Day
}

// CreateNewFiles создаёт каталог дневника заданного месяца по заданному пути.
func CreateNewFiles(month models.Month, path string) (Result, error) {
	var (
		targetMap  map[string]models.Day
		generalErr error
		res        Result = Result{}
		currentContents ScanResult
	)


	path = filepath.Join(path, month.String())

	currentContents, _ = FolderContents(path)
	res.Existed = currentContents.Existed
	res.Garbage = currentContents.Garbage

	// Получение списка дней
	res.Target = month.Days()
	targetMap = make(map[string]models.Day, 31)
	for _, day := range currentContents.Missing {
		targetMap[day.String()] = day
	}

	// Создание папки месяца
	err := os.MkdirAll(path, 0o777) // rwx(7)rwx(7)rwx(7)
	if err != nil {
		err = fmt.Errorf("Создать каталог %v не удалось: %w\n", path, err)
		return res, err
	}

	// Создание файлов дневника
	res.Created = make([]models.Day, 0, len(targetMap))
	for name, day := range targetMap {
		file, err := os.OpenFile(
			filepath.Join(path, name),
			os.O_CREATE|os.O_EXCL|os.O_WRONLY,
			0o666)
		if err != nil {
			if generalErr != nil {
				generalErr = fmt.Errorf("%v\nНе удалось создать файл %v: %w", generalErr, name, err)
			} else {
				generalErr = fmt.Errorf("Не удалось создать файл %v: %w", name, err)
			}
		} else {
			defer file.Close()
			res.Created = append(res.Created, day)
		}
	}

	res.Created = models.SortDays(res.Created)
	res.Existed = models.SortDays(res.Existed)
	slices.Sort(res.Garbage)

	return res, generalErr
}
