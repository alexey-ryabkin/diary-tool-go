package main

import (
	"bufio"
	"fmt"
	"os"
	"time"

	"github.com/alexey-ryabkin/diary-tool-go/internal/fs"
	"github.com/alexey-ryabkin/diary-tool-go/internal/models"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run() error {
	var (
		currentMonth, nextMonth models.Month
		timeNow                 time.Time
		res                     fs.Result
	)

	if err := printLastFullMonth(); err != nil {
		fmt.Printf("Определить последний созданный месяц невозможно: %v\n", err)
	}

	timeNow = time.Now()
	currentMonth = models.Month{
		Year:  timeNow.Year(),
		Month: int(timeNow.Month()),
	}
	nextMonth = currentMonth.Next()
	fmt.Printf("1 - создать текущий месяц - %v\n", currentMonth)
	fmt.Printf("2 (по умолчанию) - создать следующий месяц - %v\n", nextMonth)

	consoleReader := bufio.NewReader(os.Stdin)
	input, _, err := consoleReader.ReadRune()
	if err != nil {
		return err
	}

	wd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("Определить рабочую директорию невозможно, создать файлики негде: %w", err)
	}

	if input == '1' {
		fmt.Println("Будет создан текущий месяц.")
		res, err = fs.CreateNewFiles(currentMonth, wd)
	} else {
		fmt.Println("Будет создан следующий месяц.")
		res, err = fs.CreateNewFiles(nextMonth, wd)
	}
	if err != nil {
		return err
	}

	fmt.Printf("Успешно создано файлов - %v/%v:\n", len(res.Created), len(res.Target))
	for _, file := range res.Created {
		fmt.Println(file)
	}
	fmt.Printf("Файлов уже было ранее - %v:\n", len(res.Existed))
	for _, file := range res.Existed {
		fmt.Println(file)
	}
	fmt.Printf("Лишних файлов в папке - %v:\n", len(res.Garbage))
	for _, file := range res.Garbage {
		fmt.Println(file)
	}

	return nil
}

// printLastFullMonth выводит на экран имя последнего корректно созданного месяца.
func printLastFullMonth() error {
	workingDir, err := os.Getwd()
	if err != nil {
		return err
	}
	lastFullMonth, err := fs.GetLatestFullMonth(workingDir)
	if err != nil {
		return err
	}
	fmt.Printf("Последний полностью созданный месяц - %v\n", lastFullMonth)
	return nil
}
