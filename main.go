package main

import (
	"bufio"
	"cmp"
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"slices"
	"strconv"
	"strings"
)

const (
	binLength     = 6
	minCardLength = 13
	maxCardLength = 19
)

type Bank struct {
	Name     string
	StartBIN int
	EndBIN   int
}

func main() {
	dbPath := flag.String("db", "banks.txt", "Путь к файлу с базой банков")
	flag.Parse()

	fmt.Println("Загрузка программы валидации банковских карт космической системы...")

	banks, err := loadBanks(*dbPath)
	if err != nil {
		log.Fatalf("Критическая ошибка: не удалось загрузить базу банков: %v", err)
	}

	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("Введите номер карты (или пустую строку для выхода): ")

		if !scanner.Scan() {
			break
		}

		rawInput := scanner.Text()
		if strings.TrimSpace(rawInput) == "" {
			fmt.Println("\nВвод завершён. Закрытие программы.")
			break
		}

		cardNumber, isValid := sanitizeAndValidate(rawInput)
		if !isValid {
			fmt.Printf("Ошибка: введён некорректный формат номера карты (допускаются только цифры, пробелы и тире, длина от %v до %v цифр\n", minCardLength, maxCardLength)
			continue
		}

		if !isValidLuhn(cardNumber) {
			fmt.Println("Ошибка: введённый номер карты не прошёл проверку по алгоритму Луна")
			continue
		}

		bin, err := parseBIN(cardNumber)
		if err != nil {
			fmt.Println("Ошибка при извлечении BIN:", err)
			continue
		}

		bankName, found := findBankByBIN(bin, banks)

		if !found {
			fmt.Printf("Внимание: не удалось определить банк для карты %v\n", cardNumber)
			continue
		}

		fmt.Printf("Карта %v относится к банку: %v\n", cardNumber, bankName)
	}

	if err := scanner.Err(); err != nil {
		log.Printf("Ошибка чтения ввода: %v\n", err)
	}
}

func loadBanks(path string) ([]Bank, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.FieldsPerRecord = 3

	var banks []Bank
	lineNum := 0

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}

		if err != nil {
			return nil, err
		}

		lineNum++

		startBIN, err := strconv.Atoi(record[1])
		if err != nil {
			return nil, fmt.Errorf("не удалось преобразовать значение startBIN банка %v в строке %v, ошибка: %v", record[0], lineNum, err)
		}

		endBIN, err := strconv.Atoi(record[2])
		if err != nil {
			return nil, fmt.Errorf("не удалось преобразовать значение endBIN банка %v в строке %v, ошибка: %v", record[0], lineNum, err)
		}

		banks = append(banks, Bank{
			Name:     record[0],
			StartBIN: startBIN,
			EndBIN:   endBIN,
		})
	}

	slices.SortFunc(banks, func(a, b Bank) int {
		return cmp.Compare(a.StartBIN, b.StartBIN)
	})

	return banks, nil
}

func sanitizeAndValidate(input string) (string, bool) {
	var sb strings.Builder
	sb.Grow(len(input))

	for _, r := range input {
		if r == ' ' || r == '-' {
			continue
		}

		if r < '0' || r > '9' {
			return "", false
		}

		sb.WriteRune(r)
	}

	cleanInput := sb.String()

	if len(cleanInput) < minCardLength || len(cleanInput) > maxCardLength {
		return "", false
	}

	return cleanInput, true
}

func isValidLuhn(cardNumber string) bool {
	var sum int
	shouldDouble := false

	for i := len(cardNumber) - 1; i >= 0; i-- {
		digit := int(cardNumber[i] - '0')

		if shouldDouble {
			digit *= 2

			if digit > 9 {
				digit -= 9
			}
		}

		sum += digit
		shouldDouble = !shouldDouble
	}

	return sum%10 == 0
}

func parseBIN(cardNumber string) (int, error) {
	if len(cardNumber) < binLength {
		return 0, fmt.Errorf("номер карты слишком короткий для извлечения BIN, длина должна быть >= %v", binLength)
	}

	return strconv.Atoi(cardNumber[:binLength])
}

func findBankByBIN(BIN int, banks []Bank) (string, bool) {
	idx, found := slices.BinarySearchFunc(banks, BIN, func(b Bank, target int) int {
		if b.StartBIN > target {
			return 1
		}

		if b.EndBIN < target {
			return -1
		}

		return 0
	})

	if found {
		return banks[idx].Name, true
	}

	return "", false
}
