package useCase

import (
	"fmt"
	"math/rand"
	"time"
)

func GenerateClientID() string {
	timestamp := time.Now().UnixNano() / int64(time.Millisecond)
	randomNumbers := RandSeq(19)
	return fmt.Sprintf("%d-%s", timestamp, randomNumbers)
}

func RandSeq(n int) string {
	const letters = "0123456789"
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
