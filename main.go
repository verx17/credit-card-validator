package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

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
