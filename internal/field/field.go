package field

import (
	"fmt"
	"image"
	"image/color"
	"image/gif"
	"image/png"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/jimmykodes/cursor"

	"github.com/jimmykodes/a_MAZE_ing/internal/node"
	"github.com/jimmykodes/a_MAZE_ing/internal/output"
	"github.com/jimmykodes/a_MAZE_ing/internal/side"
)

var red = color.RGBA{R: 255, G: 0, B: 0, A: 255}

type Field struct {
	Width     int
	Height    int
	StartSide side.Side
	Animate   bool
	Output    output.Output
	Start     *node.Node
	End       *node.Node
	cursor    *cursor.Cursor
	Nodes     [][]*node.Node
	frames    []*image.Paletted
	scale     int

	// gen is the chosen generation function
	gen func()

	// current is the node currently being examined in the relevant search algorithm
	current *node.Node

	// available is used to represent the nodes not visited surrounding the current node
	available [4]*node.Node
}

func New(width, height, scale int, startSide side.Side, out output.Output, animate bool, function string) *Field {
	f := &Field{
		Width:     width,
		Height:    height,
		scale:     scale,
		Output:    out,
		Animate:   animate,
		StartSide: startSide,
	}
	f.cursor = cursor.New(os.Stdout)
	switch function {
	case "dfs":
		f.gen = f.dfs
	case "bfs":
		f.gen = f.bfs
	case "prim":
		f.gen = f.prim
	default:
		f.gen = func() {
			panic("invalid generation algorithm")
		}
	}

	switch startSide {
	case side.Left:
		f.Start = &node.Node{X: 0, Y: rand.Intn(height), IsStart: true}
		f.End = &node.Node{X: width - 1, Y: rand.Intn(height), IsEnd: true}
	case side.Right:
		f.Start = &node.Node{X: width - 1, Y: rand.Intn(height), IsStart: true}
		f.End = &node.Node{X: 0, Y: rand.Intn(height), IsEnd: true}
	case side.Top:
		f.Start = &node.Node{X: rand.Intn(width), Y: 0, IsStart: true}
		f.End = &node.Node{X: rand.Intn(width), Y: height - 1, IsEnd: true}
	case side.Bottom:
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

func (f *Field) Gen() {
	f.current = f.Start
	f.gen()
}

func (f *Field) animator() (animate func(), close func()) {
	init := false
	switch f.Output {
	case output.Image:
		close = func() {
			gifFile, err := os.Create("maze.gif")
			if err != nil {
				fmt.Println("error saving animation", err)
				return
			}
			defer gifFile.Close()
			delays := make([]int, len(f.frames))
			anim := &gif.GIF{Image: f.frames, Delay: delays}
			if err := gif.EncodeAll(gifFile, anim); err != nil {
				fmt.Println("error saving animation", err)
				return
			}
		}
		animate = func() {
			if err := f.WriteFrame(); err != nil {
				fmt.Println("error generating frame", err)
			}
		}
	case output.Text:
		f.cursor.AltBuffer()
		f.cursor.Hide()
		close = func() {
			f.cursor.OriginalBuffer()
			f.cursor.Show()
		}
		animate = func() {
			if init {
				f.cursor.Up(f.Height*2 + 1)
			} else {
				init = true
			}
			time.Sleep(time.Second / 10)
			if err := f.WriteFrame(); err != nil {
				fmt.Println("error generating frame", err)
			}
		}
	}

	return animate, close
}

// updateAvailable will set the available nodes from the current node into the avaiable
// array. It will return the number of nodes in the available array are relevant to this
// node.
func (f *Field) updateAvailable() int {
	count := 0
	if f.current.X > 0 {
		// not in the first column so look left
		if l := f.Nodes[f.current.Y][f.current.X-1]; !l.Visited {
			f.available[count] = l
			count++
		}
	}
	if f.current.X < f.Width-1 {
		// not in the last column so look right
		if r := f.Nodes[f.current.Y][f.current.X+1]; !r.Visited {
			f.available[count] = r
			count++
		}
	}
	if f.current.Y > 0 {
		// not in the first row, look up
		if t := f.Nodes[f.current.Y-1][f.current.X]; !t.Visited {
			f.available[count] = t
			count++
		}
	}
	if f.current.Y < f.Height-1 {
		// not in last row, look down
		if b := f.Nodes[f.current.Y+1][f.current.X]; !b.Visited {
			f.available[count] = b
			count++
		}
	}
	return count
}

func (f *Field) Repr() [][]uint8 {
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
				case side.Left:
					repr[y][x-1] = 1
				case side.Right:
					repr[y][x+1] = 1
				case side.Top:
					repr[y-1][x] = 1
				case side.Bottom:
					repr[y+1][x] = 1
				}
			} else if n.IsEnd {
				switch f.StartSide {
				case side.Left:
					repr[y][x+1] = 1
				case side.Right:
					repr[y][x-1] = 1
				case side.Top:
					repr[y+1][x] = 1
				case side.Bottom:
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

func (f *Field) WriteFrame() error {
	switch f.Output {
	case output.Text:
		f.writeText(true)
		return nil
	case output.Image:
		f.frames = append(f.frames, f.genImage(true))
		return nil
	default:
		return fmt.Errorf("invalid output type")
	}
}

func (f *Field) WriteImage(name string) error {
	img := f.genImage(false)
	imgFile, err := os.Create(name)
	if err != nil {
		return err
	}
	defer imgFile.Close()
	return png.Encode(imgFile, img)
}

func (f *Field) genImage(colorCurrent bool) *image.Paletted {
	r := f.Repr()
	img := image.NewPaletted(
		image.Rect(0, 0, f.scale*((f.Width*2)+1), f.scale*((f.Height*2)+1)),
		color.Palette{color.White, color.Black, red},
	)
	for y, row := range r {
		for x, col := range row {
			var c color.Color
			switch col {
			case 0:
				c = color.Black
			case 1:
				c = color.White
			case 2:
				if colorCurrent {
					c = color.RGBA{R: 255, G: 0, B: 0, A: 255}
				} else {
					c = color.White
				}
			}
			for i := 0; i < f.scale; i++ {
				for j := 0; j < f.scale; j++ {
					img.Set((x*f.scale)+i, (y*f.scale)+j, c)
				}
			}

		}
	}
	return img
}

func (f *Field) WriteText() {
	f.writeText(false)
}

func (f *Field) writeText(colorCurrent bool) {
	r := f.Repr()
	for _, row := range r {
		for _, col := range row {
			switch col {
			case 0:
				f.cursor.WhiteBG()
			case 1:
				f.cursor.BlackBG()
			case 2:
				if colorCurrent {
					f.cursor.RedBG()
				} else {
					f.cursor.BlackBG()
				}
			}
			fmt.Print("  ")
			f.cursor.Clear()
		}
		fmt.Println()
	}
}

func (f *Field) String() string {
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
