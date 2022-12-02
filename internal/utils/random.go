package utils

import (
	"math/rand"
	"time"
)

func GetRandomItem(messages []string) string {
	rand.Seed(time.Now().Unix())
	n := rand.Int() % len(messages)
	return messages[n]
}
