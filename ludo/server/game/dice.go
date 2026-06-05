package game

import (
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func RollDice() int {
	return rand.Intn(6) + 1
}
