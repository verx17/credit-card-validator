package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
)

func main() {
	fmt.Println("Загрузка программы валидации банковских карт космической системы...")

	banks, err := loadBankData("banks.txt")
	if err != nil {
		log.Fatal(err)
	}

	for {
		cardNumber := getUserInput()

		if len(cardNumber) == 0 {
			fmt.Println("Закрытие программы...")
			break
		}

		if isValid := validateInput(cardNumber); !isValid {
			fmt.Println("Введён некорректный номер карты")
			continue
		}

		if isValid := validateLuhn(cardNumber); !isValid {
			fmt.Println("Введённый номер карты не прошёл проверку")
			continue
		}

		bin := extractBIN(cardNumber)
		bankName := identifyBank(bin, banks)

		if bankName == EUB {
			fmt.Printf("Ошибка: не удалось определить Банк карты %v\n", cardNumber)
			continue
		}

		fmt.Printf("Карта %v относится к банку %v\n", cardNumber, bankName)
	}
}

type Bank struct {
	Name    string
	BinFrom int
	BinTo   int
}

func loadBankData(path string) (banks []Bank, err error) {
	file, err := os.Open(path)
	if err != nil {
		log.Fatal("Error opening file: ", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)

	for {
		bank, err := reader.Read()
		if err == io.EOF {
			break
		}

		if err != nil {
			return nil, err
		}

		binFrom, _ := strconv.Atoi(bank[1])
		binTo, _ := strconv.Atoi(bank[2])

		banks = append(banks, Bank{
			Name:    bank[0],
			BinFrom: binFrom,
			BinTo:   binTo,
		})
	}

	return banks, nil
}

const EUB = "Unknown Bank"

func identifyBank(bin int, banks []Bank) (BankName string) {
	for _, bank := range banks {
		if bin >= bank.BinFrom && bin <= bank.BinTo {
			return bank.Name
		}
	}

	return EUB
}

func extractBIN(cardNumber string) (bin int) {
	bin, _ = strconv.Atoi(cardNumber[:6])
	return bin
}

func validateLuhn(cardNumber string) bool {
	var sum int
	shouldDouble := false

	for i := len(cardNumber) - 1; i >= 0; i-- {
		r := rune(cardNumber[i])
		digit := int(r - '0')

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

func getUserInput() string {
	fmt.Println("Введите номер карты")

	reader := bufio.NewReader(os.Stdin)
	userInput, _ := reader.ReadString('\n')
	userInput = strings.TrimSpace(userInput)
	userInput = strings.ReplaceAll(userInput, " ", "")

	return userInput
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
