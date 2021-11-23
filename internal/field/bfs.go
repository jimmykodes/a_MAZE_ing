package field

import (
	"math/rand"

	"github.com/jimmykodes/a_MAZE_ing/internal/node"
)

// bfs is a breadth-first-search maze generation method
//
// It isn't very good. Creates a lot of straight corridors and leaves some empty spaces.
// Not sure how much of this is my implementation vs the algorithm as a whole.
// Will have to investigate later.
func (f *Field) bfs() {
	f.Start.Visited = true
	stack := []*node.Node{f.Start}
	var (
		count     = 0
		deferFunc func()
		animate   func()
	)
	defer func() {
		if deferFunc != nil {
			deferFunc()
		}
	}()
	for {
		if f.Animate {
			if animate == nil {
				animate, deferFunc = f.animator()
			}
			animate()
		}
		// grab the front element of the stack
		f.current = stack[0]

		count = f.updateAvailable()
		rand.Shuffle(count, func(i, j int) {
			f.available[i], f.available[j] = func() (*node.Node, *node.Node) { return f.available[j], f.available[i] }()
		})
		num := rand.Intn(count + 1)
		if num == 0 && count > 0 {
			// if we randomly selected 0, but we have more than 0 items in the count, make sure we are selecting
			// at least one of these
			num = 1
		}
		for i := 0; i < num; i++ {
			n := f.available[i]
			n.Parent = f.current
			n.Visited = true
			stack = append(stack, f.available[i])
		}
		if len(stack) == 1 {
			break
		}

		stack = stack[1:]
	}
}
