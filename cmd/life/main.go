package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/418Coffee/life"
)

func main() {
	rand.Seed(0)
	l := life.NewGame(40, 15, true)
	for i := 0; i < 30; i++ {
		l.Tick()
		fmt.Print("\x1bc", l)
		time.Sleep(time.Second / 30)
	}
}
