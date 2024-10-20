package utils

import (
	"encoding/json"
	"fmt"
	"strings"
)

func StructToString(data interface{}) string {
	result, err := json.Marshal(data)
	if err != nil {
		return ""
	}

	return string(result)
}

func FormatPhoneNumber(phone string) string {
	phone = strings.TrimLeft(phone, " +")

	// Check if the phone number starts with "62"
	if strings.HasPrefix(phone, "62") {
		return phone
	}

	// Check if the phone number starts with "+62"
	if strings.HasPrefix(phone, "+62") {
		// Remove the "+" character and return the number
		return strings.TrimPrefix(phone, "+")
	}

	if strings.Contains(phone, "0") {
		// Replace the "0" with "62" as the prefix and return the number
		return fmt.Sprintf("62%s", strings.Replace(phone, "0", "", 1))
	}

	// If the phone number doesn't start with "62" or "+62", add "62" as the prefix
	return fmt.Sprintf("62%s", phone)
}
