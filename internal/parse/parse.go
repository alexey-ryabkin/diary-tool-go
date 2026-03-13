package parse

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/alexey-ryabkin/diary-tool-go/internal/models"
)

// FolderName переводит наименование папки в объект, описывающий месяц дневника.
// Если папка называется не по стандарту, выдаёт ошибку с описанием несоответствия.
func FolderName(name string) (models.Month, error) {
	var (
		validMonth models.Month
		canonName string
	)

	parts := strings.Split(name, ",")
	if len(parts) < 2 {
		return models.Month{}, errors.New("наименование папки должно содержать запятую")
	}

	parts = strings.Split(parts[0], ".")
	if len(parts) < 2 {
		return models.Month{}, errors.New("год и месяц должны быть разделены точкой")
	}

	year, err := strconv.Atoi(parts[0])
	if err != nil {
		return models.Month{}, err
	}
	if year < 0 || year > 99 {
		return models.Month{}, errors.New("год должен быть указан двумя цифрами")
	}
	year = 2000 + year

	month, err := strconv.Atoi(parts[1])
	if err != nil {
		return models.Month{}, err
	}
	if month < 1 || month > 12 {
		return models.Month{}, errors.New("месяц должен быть от 1 до 12 включительно")
	}

	validMonth = models.Month{Year: year, Month: month}
	canonName = validMonth.String()
	if name != canonName {
		return models.Month{}, fmt.Errorf("формат %v не соответствует каноническому %v", name, canonName)
	}

	return validMonth, nil
}
