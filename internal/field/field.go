package field

import (
	"fmt"
	"math/rand"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/jimmykodes/cursor"

	"github.com/jimmykodes/a_MAZE_ing/internal/node"
	"github.com/jimmykodes/a_MAZE_ing/internal/output"
)

type Side int

const (
	Left Side = iota
	Right
	Top
	Bottom
)

type Field struct {
	Width     int
	Height    int
	StartSide Side
	Animate   bool
	Output    output.Output
	Start     *node.Node
	End       *node.Node
	current   *node.Node
	cursor    *cursor.Cursor
	Nodes     [][]*node.Node
}

func New(width, height int, startSide Side, out output.Output, animate bool) *Field {
	f := &Field{
		Width:     width,
		Height:    height,
		Output:    out,
		Animate:   animate,
		StartSide: startSide,
	}
	if out == output.Text {
		f.cursor = cursor.New(os.Stdout)
	}
	switch startSide {
	case Left:
		f.Start = &node.Node{X: 0, Y: rand.Intn(height), IsStart: true}
		f.End = &node.Node{X: width - 1, Y: rand.Intn(height), IsEnd: true}
	case Right:
		f.Start = &node.Node{X: width - 1, Y: rand.Intn(height), IsStart: true}
		f.End = &node.Node{X: 0, Y: rand.Intn(height), IsEnd: true}
	case Top:
		f.Start = &node.Node{X: rand.Intn(width), Y: 0, IsStart: true}
		f.End = &node.Node{X: rand.Intn(width), Y: height - 1, IsEnd: true}
	case Bottom:
		f.Start = &node.Node{X: rand.Intn(width), Y: height - 1, IsStart: true}
		f.End = &node.Node{X: rand.Intn(width), Y: 0, IsEnd: true}
	}
	f.Nodes = make([][]*node.Node, height)
	for y := 0; y < height; y++ {
		f.Nodes[y] = make([]*node.Node, width)
		for x := 0; x < width; x++ {
			if y == f.Start.Y && x == f.Start.X {
				f.Nodes[y][x] = f.Start
			} else if y == f.End.Y && x == f.End.X {
				f.Nodes[y][x] = f.End
			} else {
				f.Nodes[y][x] = &node.Node{X: x, Y: y}
			}
		}
	}
	return f
}

func (f Field) Gen() {
	f.current = f.Start
	f.dfs()
}

// dfs is a depth-first-search maze generation method
func (f Field) dfs() {
	var (
		available = make([]*node.Node, 4)
		count     = 0
		once      sync.Once
		deferFunc func()
		init      bool
	)
	defer func() {
		if deferFunc != nil {
			deferFunc()
		}
	}()
	for {
		f.current.Visited = true
		if f.Animate {
			if f.Output == output.Text {
				once.Do(func() {
					f.cursor.AltBuffer()
					f.cursor.Hide()
					deferFunc = func() {
						f.cursor.OriginalBuffer()
						f.cursor.Show()
					}
				})
				if init {
					f.cursor.Up(f.Height*2 + 1)
				} else {
					init = true
				}
				f.WriteFrame()
				time.Sleep(time.Second / 60)
			}
		}

		if f.current.IsEnd {
			f.current = f.current.Parent
			continue
		}
		// reset count per loop
		count = 0
		if f.current.X > 0 {
			// not in the first column so look left
			if l := f.Nodes[f.current.Y][f.current.X-1]; !l.Visited {
				available[count] = l
				count++
			}
		}
		if f.current.X < f.Width-1 {
			// not in the last column so look right
			if r := f.Nodes[f.current.Y][f.current.X+1]; !r.Visited {
				available[count] = r
				count++
			}
		}
		if f.current.Y > 0 {
			// not in the first row, look up
			if t := f.Nodes[f.current.Y-1][f.current.X]; !t.Visited {
				available[count] = t
				count++
			}
		}
		if f.current.Y < f.Height-1 {
			// not in last row, look down
			if b := f.Nodes[f.current.Y+1][f.current.X]; !b.Visited {
				available[count] = b
				count++
			}
		}

		if count == 0 {
			if p := f.current.Parent; p != nil {
				f.current = p
				continue
			} else {
				break
			}
		}

		next := available[rand.Intn(count)]
		next.Parent = f.current
		f.current = next
	}
}

func (f Field) Repr() [][]uint8 {
	repr := make([][]uint8, (f.Height*2)+1)
	for i := 0; i < (f.Height*2)+1; i++ {
		repr[i] = make([]uint8, (f.Width*2)+1)
		for x := range repr[i] {
			repr[i][x] = 0
		}
	}

	for _, nodes := range f.Nodes {
		for _, n := range nodes {
			if !n.Visited && !n.IsEnd && !n.IsStart {
				continue
			}
			x := (n.X * 2) + 1
			y := (n.Y * 2) + 1
			if f.current != nil && n.X == f.current.X && n.Y == f.current.Y {
				repr[y][x] = 2
			} else {
				repr[y][x] = 1
			}
			if n.IsStart {
				switch f.StartSide {
				case Left:
					repr[y][x-1] = 1
				case Right:
					repr[y][x+1] = 1
				case Top:
					repr[y-1][x] = 1
				case Bottom:
					repr[y+1][x] = 1
				}
			} else if n.IsEnd {
				switch f.StartSide {
				case Left:
					repr[y][x+1] = 1
				case Right:
					repr[y][x-1] = 1
				case Top:
					repr[y+1][x] = 1
				case Bottom:
					repr[y-1][x] = 1
				}
			}
			if n.Parent != nil {
				if n.Parent.X == n.X {
					// in the same column
					if n.Parent.IsAbove(n) {
						repr[y-1][x] = 1
					} else {
						repr[y+1][x] = 1
					}
				} else {
					// in the same row
					if n.Parent.IsLeft(n) {
						repr[y][x-1] = 1
					} else {
						repr[y][x+1] = 1
					}
				}
			}
		}
	}
	return repr
}

func (f Field) WriteFrame() {
	if f.Output == output.Text {
		r := f.Repr()
		for _, row := range r {
			for _, col := range row {
				switch col {
				case 0:
					f.cursor.WhiteBG()
				case 1:
					f.cursor.BlackBG()
				case 2:
					f.cursor.RedBG()
				}
				fmt.Print("  ")
				f.cursor.Clear()
			}
			fmt.Println()
		}
	}
}

func (f Field) String() string {
	repr := f.Repr()
	intermediate := make([]string, len(repr))
	for i, r := range repr {
		s := make([]string, len(r))
		for i2, i3 := range r {
			var place string
			switch i3 {
			case 0:
				place = "X"
			case 1:
				place = " "
			case 2:
				place = "0"
			}
			s[i2] = place
		}
		intermediate[i] = strings.Join(s, " ")
	}
	return strings.Join(intermediate, "\n")
}
