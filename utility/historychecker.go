package utility

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

// func to check if file exists

func Filecheck(filename string) bool {
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return false
	}
	return true
}

func ParsePrice(price string) []float64 {
	r, _ := regexp.Compile("[+-]?([0-9]*[,])?[0-9]+")

	var prices []float64

	fmt.Println("Price: ", r.FindAllString(price, -1))
	for _, p := range r.FindAllString(price, -1) {
		p = strings.Replace(p, ",", ".", -1)
		if s, err := strconv.ParseFloat(p, 32); err == nil {
			prices = append(prices, s)
		}
	}

	return prices
}