package lib

import (
	"fmt"
	"math/rand"
)

func Contains(elems []string, v string) bool {
	for _, s := range elems {
		if v == s {
			return true
		}
	}
	return false
}

func RandomStr(length int) string {
	str := ""
	allChars := "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"
	for i := 0; i < length; i++ {
		str += string(allChars[rand.Intn(len(allChars))])
	}
	fmt.Println(str)
	return str
}

func Base62(n int64) string {
	// Convert a number n from decimal base to base 62
	// n is a number in decimal base
	// returns a string in base 62
	var base62 string = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	var base10 int64 = 62
	var base62_int int64 = 0
	var remainder int64 = 0
	var result string = ""
	var i int = 0
	for n > 0 {
		remainder = n % base10
		n = n / base10
		base62_int = remainder
		result = string(base62[base62_int]) + result
		i++
	}
	fmt.Println(result)
	return result
}
