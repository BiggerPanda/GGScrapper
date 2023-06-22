package utility

import (
	"os"
)

// func to check if file exists

func filecheck(filename string) bool {
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return false
	}
	return true
}
