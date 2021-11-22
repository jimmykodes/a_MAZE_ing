package main

import (
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"time"

	"github.com/jimmykodes/a_MAZE_ing/internal/field"
	"github.com/jimmykodes/a_MAZE_ing/internal/output"
)

var (
	width     = 10
	height    = 10
	startSide = field.Top
	out       = output.Text
)

func main() {
	rand.Seed(time.Now().Unix())
	a := os.Args[1:]
	switch len(a) {
	case 1:
		b, err := strconv.Atoi(a[0])
		if err != nil {
			fmt.Println("invalid input, must be a number")
			return
		}
		width = b
		height = b
	case 2:
		var err error
		width, err = strconv.Atoi(a[0])
		if err != nil {
			fmt.Println("invalid input, must be a number")
			return
		}
		height, err = strconv.Atoi(a[1])
		if err != nil {
			fmt.Println("invalid input, must be a number")
			return
		}
	}
	maze := field.New(width, height, startSide, out, true)
	maze.Gen()
	fmt.Println(maze)
}
