package parse

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/alexey-ryabkin/diary-tool-go/internal/models"
)

func parseFolderName(name string) (models.Month, error) {
	parts := strings.Split(name, ",")
	if len(parts) < 2 {
		return models.Month{}, errors.New("наименование папки не содержит запятой")
	}

	parts = strings.Split(parts[0], ".")
	if len(parts) < 2 {
		return models.Month{}, errors.New("год и месяц не разделены точкой")
	}

	year, err := strconv.Atoi(parts[0])
	if err != nil {
		return models.Month{}, err
	}
	if year < 0 || year > 99 {
		return models.Month{}, errors.New("год указан не двумя цифрами")
	}
	year = 2000 + year

	month, err := strconv.Atoi(parts[1])
	if err != nil {
		return models.Month{}, err
	}
	if month < 1 || month > 12 {
		return models.Month{}, errors.New("месяц не может быть меньше одного или больше двенадцати")
	}

	validMonth := models.Month{year, month}
	canon := validMonth.String()
	if name != canon {
		errorMsg := fmt.Sprint("формат ", name, " не соответствует каноническому ", canon)
		return models.Month{}, errors.New(errorMsg)
	}

	return validMonth, nil
}
