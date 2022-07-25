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
