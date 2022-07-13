package util

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

const alphabet = "abcdefghijklmnopqrstuvwxyz"

func init() {
	rand.Seed(time.Now().UnixNano())
}

func RandomFloat(min, max int64) float64 {
	return float64(min +  rand.Int63n(max - min + 1)) + rand.Float64()
}

func RandomString(n int) string {
	var sb strings.Builder
	k := len(alphabet)

	for i := 0; i < n; i++ {
		c := alphabet[rand.Intn(k)]
		sb.WriteByte(c)
	}

	return sb.String()
}

func RandomFullName() string {
	return RandomString(6)
}

func RandomMoney() string {
	return fmt.Sprintf("%f",RandomFloat(0, 1000))
}

func RandomEmail() string {
	return fmt.Sprintf("%s@email.com", RandomString(6))
}