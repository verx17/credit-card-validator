package main

import (
	"cmp"
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"slices"
	"strconv"
)

const UnknownBankName = "Unknown Bank"

type Bank struct {
	Name    string
	BinFrom int
	BinTo   int
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
