package util

import (
	"encoding/json"
	"fmt"
)

func StringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func JsonEncode(data interface{}) string {
	b, _ := json.Marshal(data)

	return string(b)
}

func PrintGreen(text string) {
	fmt.Printf("\033[1;32m%s\033[0m", text)
}

func PrintRed(text string) {
	fmt.Printf("\033[1;31m%s\033[0m", text)
}
