package utils

import (
	"fmt"
	"regexp"
	"strings"
	"time"
)

func RemoveInvalidCharacters(text string) string {
	// Regular expression to find invalid characters
	reg := regexp.MustCompile("[^\x00-\x7F]+")

	// Remove invalid characters from the text using the regular expression
	cleanedText := reg.ReplaceAllString(text, "")

	textReturn := strings.ToLower(cleanedText)

	return textReturn
}

// convertDateToStandardFormat try to identify the format of the input data and convert it to YYYY-MM-DD
func ConvertDateToStandardFormat(dateStr string) (string, error) {
	formats := []string{"02/01/2006", "2006-31-01", "01/02/2006"}

	for _, format := range formats {
		if date, err := time.Parse(format, dateStr); err == nil {
			return date.Format("2006-01-02"), nil
		}
	}

	return "", fmt.Errorf("formato de data desconhecido: %s", dateStr)
}
