package main

import (
	"fmt"
	"strconv"
)

const BinLength = 6

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

func extractBIN(cardNumber string) (int, error) {
	if len(cardNumber) < BinLength {
		return 0, fmt.Errorf("Номер карты слишком короткий для извлечения BIN, длина должна быть >= %v", BinLength)
	}

	return strconv.Atoi(cardNumber[:BinLength])
}
