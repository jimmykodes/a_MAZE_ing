package main

import (
	"flag"
	"fmt"
	"math/rand"
	"time"

	"github.com/jimmykodes/a_MAZE_ing/internal/field"
	"github.com/jimmykodes/a_MAZE_ing/internal/output"
)

var (
	width      int
	height     int
	scale     int
	seed       int64
	outputType string
	startSide  string
	animate    bool
)

func main() {
	flag.BoolVar(&animate, "animate", false, "animate generation process")
	flag.StringVar(&startSide, "start", "t", "side of maze to start on [t | b | l | r]")
	flag.StringVar(&outputType, "output", "image", "output type [image | text]")
	flag.IntVar(&width, "width", 10, "maze width")
	flag.IntVar(&height, "height", 10, "maze height")
	flag.IntVar(&scale, "scale", 1, "maze scale")
	flag.Int64Var(&seed, "seed", time.Now().Unix(), "maze seed value")
	flag.Parse()
	rand.Seed(seed)
	var ss field.Side
	switch startSide {
	case "t", "top":
		ss = field.Top
	case "b", "bottom":
		ss = field.Bottom
	case "l", "left":
		ss = field.Left
	case "r", "right":
		ss = field.Right
	default:
		fmt.Println("invalid start side")
		return
	}

	var out output.Output
	switch outputType {
	case "image":
		out = output.Image
	case "text":
		out = output.Text
	}

	maze := field.New(width, height, scale, ss, out, animate)
	maze.Gen()
	maze.WriteText()
	err := maze.WriteImage("maze.png")
	if err != nil {
		fmt.Println("error saving maze", err)
		return
	}
}
