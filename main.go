package main

import (
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

type side int

const (
	left side = iota
	right
	top
	bottom
)

var (
	width     = 10
	height    = 10
	startSide = top
)

type Node struct {
	x, y    int
	start   bool
	end     bool
	visited bool
	parent  *Node
}

func (n Node) IsLeft(n2 *Node) bool {
	return n.x < n2.x
}

func (n Node) IsAbove(n2 *Node) bool {
	return n.y < n2.y
}

func (n Node) String() string {
	return fmt.Sprintf("(%d, %d)", n.x, n.y)
}

type Field [][]*Node

func (f Field) String() string {
	repr := make([][]string, (height*2)+1)
	for i := 0; i < (height*2)+1; i++ {
		repr[i] = make([]string, (width*2)+1)
		for x := range repr[i] {
			repr[i][x] = "X"
		}
	}

	for _, nodes := range f {
		for _, node := range nodes {
			x := (node.x * 2) + 1
			y := (node.y * 2) + 1
			repr[y][x] = " "
			if node.start {
				switch startSide {
				case left:
					repr[y][x-1] = " "
				case right:
					repr[y][x+1] = " "
				case top:
					repr[y-1][x] = " "
				case bottom:
					repr[y+1][x] = " "
				}
			} else if node.end {
				switch startSide {
				case left:
					repr[y][x+1] = " "
				case right:
					repr[y][x-1] = " "
				case top:
					repr[y+1][x] = " "
				case bottom:
					repr[y-1][x] = " "
				}
			}
			if node.parent != nil {
				if node.parent.x == node.x {
					// in the same column
					if node.parent.IsAbove(node) {
						repr[y-1][x] = " "
					} else {
						repr[y+1][x] = " "
					}
				} else {
					// in the same row
					if node.parent.IsLeft(node) {
						repr[y][x-1] = " "
					} else {
						repr[y][x+1] = " "
					}
				}
			}
		}
	}
	intermediate := make([]string, len(repr))
	for i, s := range repr {
		intermediate[i] = strings.Join(s, " ")
	}
	return strings.Join(intermediate, "\n")
}

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
	var (
		start *Node
		end   *Node
	)
	switch startSide {
	case left:
		start = &Node{x: 0, y: rand.Intn(height), start: true}
		end = &Node{x: width - 1, y: rand.Intn(height), end: true}
	case right:
		start = &Node{x: width - 1, y: rand.Intn(height), start: true}
		end = &Node{x: 0, y: rand.Intn(height), end: true}
	case top:
		start = &Node{x: rand.Intn(width), y: 0, start: true}
		end = &Node{x: rand.Intn(width), y: height - 1, end: true}
	case bottom:
		start = &Node{x: rand.Intn(width), y: height - 1, start: true}
		end = &Node{x: rand.Intn(width), y: 0, end: true}
	}
	var field Field
	field = make([][]*Node, height)
	for i := 0; i < height; i++ {
		field[i] = make([]*Node, width)
		for j := 0; j < width; j++ {
			if i == start.y && j == start.x {
				field[i][j] = start
			} else if i == end.y && j == end.x {
				field[i][j] = end
			} else {
				field[i][j] = &Node{x: j, y: i}
			}
		}
	}

	var (
		available = make([]*Node, 4)
		current   *Node
		count     = 0
	)
	current = start
	for {
		current.visited = true
		if current.end {
			current = end.parent
			continue
		}
		// reset count per loop
		count = 0
		if current.x > 0 {
			// not in the first column so look left
			if l := field[current.y][current.x-1]; !l.visited {
				available[count] = l
				count++
			}
		}
		if current.x < width-1 {
			// not in the last column so look right
			if r := field[current.y][current.x+1]; !r.visited {
				available[count] = r
				count++
			}
		}
		if current.y > 0 {
			// not in the first row, look up
			if t := field[current.y-1][current.x]; !t.visited {
				available[count] = t
				count++
			}
		}
		if current.y < height-1 {
			// not in last row, look down
			if b := field[current.y+1][current.x]; !b.visited {
				available[count] = b
				count++
			}
		}

		if count == 0 {
			if p := current.parent; p != nil {
				current = p
				continue
			} else {
				break
			}
		}

		next := available[rand.Intn(count)]
		next.parent = current
		current = next
	}
	fmt.Println(field)
}
