package utils

import (
	"fmt"
	"math/rand"
	"time"
)

func GetRandomString(length int) (string, error) {
	if length <= 0 {
		return "", fmt.Errorf("%s", "invalid random string length")
	}

	rand.Seed(time.Now().UnixNano())
	digits := "0123456789"
	specials := "~=+%^*/()[]{}/!@#$?|"
	all := "ABCDEFGHIJKLMNOPQRSTUVWXYZ" +"abcdefghijklmnopqrstuvwxyz" + digits + specials
	buf := make([]byte, length)
	for i := 0; i < length; i++ {
		buf[i] = all[rand.Intn(len(all))]
	}

	return string(buf), nil
}
