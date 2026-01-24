package main

import (
	"testing"
)

func TestValidateLuhn(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"1111222233334444", true},  // Scenario 3 from README (Valid Luhn)
		{"4000001234567890", false}, // Scenario 2 from README (Invalid Luhn)
		{"5555444433332226", true},  // Corrected Scenario 1
		{"5555444433332222", false}, // Original (wrong) Scenario 1
		{"79927398713", true},       // Wikipedia example
		{"79927398710", false},
	}

	for _, tt := range tests {
		result := validateLuhn(tt.input)
		if result != tt.expected {
			t.Errorf("validateLuhn(%s) = %v; want %v", tt.input, result, tt.expected)
		}
	}
}

func TestValidateInput(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"1234567890123", true},       // 13 digits
		{"1234567890123456789", true}, // 19 digits
		{"123456789012", false},       // 12 digits (too short)
		{"12345678901234567890", false}, // 20 digits (too long)
		{"1234567890abc", false},      // non-digits
		{"", false},                   // empty
	}

	for _, tt := range tests {
		result := validateInput(tt.input)
		if result != tt.expected {
			t.Errorf("validateInput(%s) = %v; want %v", tt.input, result, tt.expected)
		}
	}
}

func TestExtractBIN(t *testing.T) {
	tests := []struct {
		input     string
		wantBin   int
		expectErr bool
	}{
		{"1234567890123", 123456, false},
		{"987654321", 987654, false},
		{"12345", 0, true}, // Too short
	}

	for _, tt := range tests {
		bin, err := extractBIN(tt.input)
		if tt.expectErr {
			if err == nil {
				t.Errorf("extractBIN(%s) expected error, got nil", tt.input)
			}
		} else {
			if err != nil {
				t.Errorf("extractBIN(%s) unexpected error: %v", tt.input, err)
			}
			if bin != tt.wantBin {
				t.Errorf("extractBIN(%s) = %v; want %v", tt.input, bin, tt.wantBin)
			}
		}
	}
}

func TestIdentifyBank(t *testing.T) {
	banks := []Bank{
		{"Bank A", 100000, 199999},
		{"Bank B", 300000, 399999},
		{"Bank C", 500000, 500000}, // Single BIN
	}

	tests := []struct {
		bin      int
		expected string
	}{
		{150000, "Bank A"},
		{100000, "Bank A"},
		{199999, "Bank A"},
		{350000, "Bank B"},
		{500000, "Bank C"},
		{200000, UnknownBankName},
		{400000, UnknownBankName},
		{999999, UnknownBankName},
		{0, UnknownBankName},
	}

	for _, tt := range tests {
		result := identifyBank(tt.bin, banks)
		if result != tt.expected {
			t.Errorf("identifyBank(%d) = %s; want %s", tt.bin, result, tt.expected)
		}
	}
}
