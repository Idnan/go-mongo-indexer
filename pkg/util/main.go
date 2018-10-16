package util

import (
	"encoding/json"
	"fmt"
)

// Find a given string in an array of strings
func StringInSlice(needle string, haystack []string) bool {
	for _, elem := range haystack {
		if elem == needle {
			return true
		}
	}

	return false
}

// Converts data to json string
func JsonEncode(data interface{}) string {
	b, _ := json.Marshal(data)

	return string(b)
}

// Print to stdout with green color
func PrintGreen(text string) {
	fmt.Printf("\033[1;32m%s\033[0m", text)
}

// Print to stdout with red color
func PrintRed(text string) {
	fmt.Printf("\033[1;31m%s\033[0m", text)
}

// Print to stdout with bold font
func PrintBold(text string) {
	fmt.Printf("\033[1m%s\033[0m", text)
}
