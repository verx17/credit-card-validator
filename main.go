package main

import (
	"bufio"
	"cmp"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"slices"
	"strconv"
	"strings"
)

const UnknownBankName = "Unknown Bank"
const BinLength = 6

type Bank struct {
	Name    string
	BinFrom int
	BinTo   int
}

func main() {
	fmt.Println("Загрузка программы валидации банковских карт космической системы...")

	banks, err := loadBankData("banks.txt")
	if err != nil {
		log.Fatalf("Критическая ошибка: не удалось загрузить базу банков: %v", err)
	}

	reader := bufio.NewReader(os.Stdin)

	for {
		cardNumber, err := getUserInput(reader)
		if err != nil {
			if err == io.EOF {
				fmt.Println("\nВвод завершён. Закрытие программы.")
				break
			}

			fmt.Println("Ошибка чтения ввода:", err)
			continue
		}

		if len(cardNumber) == 0 {
			continue
		}

		if isValid := validateInput(cardNumber); !isValid {
			fmt.Println("Введён некорректный номер карты")
			continue
		}

		if isValid := validateLuhn(cardNumber); !isValid {
			fmt.Println("Введённый номер карты не прошёл проверку")
			continue
		}

		bin, err := extractBIN(cardNumber)
		if err != nil {
			fmt.Println("Ошибка при извлечении BIN:", err)
			continue
		}

		bankName := identifyBank(bin, banks)

		if bankName == UnknownBankName {
			fmt.Printf("Ошибка: не удалось определить Банк карты %v\n", cardNumber)
			continue
		}

		fmt.Printf("Карта %v относится к банку %v\n", cardNumber, bankName)
	}
}

func loadBankData(path string) ([]Bank, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	var banks []Bank
	lineNum := 0

	for {
		bank, err := reader.Read()
		if err == io.EOF {
			break
		}

		if err != nil {
			return nil, err
		}

		lineNum++

		binFrom, err := strconv.Atoi(bank[1])
		if err != nil {
			return nil, fmt.Errorf("Не удалось преобразовать значение binFrom банка %v в строке %v, ошибка: %v", bank[0], lineNum, err)
		}

		binTo, err := strconv.Atoi(bank[2])
		if err != nil {
			return nil, fmt.Errorf("Не удалось преобразовать значение binTo банка %v в строке %v, ошибка: %v", bank[0], lineNum, err)
		}

		banks = append(banks, Bank{
			Name:    bank[0],
			BinFrom: binFrom,
			BinTo:   binTo,
		})
	}

	slices.SortFunc(banks, func(a, b Bank) int {
		return cmp.Compare(a.BinFrom, b.BinFrom)
	})

	return banks, nil
}

func identifyBank(bin int, banks []Bank) string {
	idx, isFound := slices.BinarySearchFunc(banks, bin, func(b Bank, target int) int {
		if target < b.BinFrom {
			return 1
		}

		if target > b.BinTo {
			return -1
		}

		return 0
	})

	if isFound {
		return banks[idx].Name
	}

	return UnknownBankName
}

func extractBIN(cardNumber string) (int, error) {
	if len(cardNumber) < BinLength {
		return 0, fmt.Errorf("Номер карты слишком короткий для извлечения BIN, длина должна быть >= %v", BinLength)
	}

	return strconv.Atoi(cardNumber[:BinLength])
}

func validateLuhn(cardNumber string) bool {
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

func getUserInput(reader *bufio.Reader) (string, error) {
	fmt.Print("Введите номер карты: ")

	userInput, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}

	userInput = strings.TrimSpace(userInput)
	userInput = strings.ReplaceAll(userInput, " ", "")
	userInput = strings.ReplaceAll(userInput, "-", "")

	return userInput, nil
}

func validateInput(input string) bool {
	if len(input) < 13 || len(input) > 19 {
		return false
	}

	for _, r := range input {
		if r < '0' || r > '9' {
			return false
		}
	}

	return true
}
