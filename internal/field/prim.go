package field

import (
	"math/rand"

	"github.com/jimmykodes/a_MAZE_ing/internal/node"
)

// prim implements Prim's Algorithm for maze generation
func (f *Field) prim() {
	for y, nodes := range f.Nodes {
		for x, n := range nodes {
			if n.Weights == nil {
				n.Weights = make(map[*node.Node]float64)
			}
			// loop throw nodes and set their left and top weights, and then
			// the right and bottom sides of the nodes left and top nodes
			// this should add weights connecting all nodes
			if y > 0 {
				// not top row, add top
				w := rand.Float64()
				n.Weights[f.Nodes[y-1][x]] = w
				// make sure above node has same weight looking down
				f.Nodes[y-1][x].Weights[n] = w
			}
			if x > 0 {
				// not first column, add left
				w := rand.Float64()
				n.Weights[f.Nodes[y][x-1]] = w
				f.Nodes[y][x-1].Weights[n] = w
			}
		}
	}
	ends := map[*node.Node]struct{}{
		f.Start: {},
	}
	animate, stop := f.animator()
	defer stop()
	for {
		var (
			parent *node.Node
			next   *node.Node
			// setting to 2.0 because all weights are random floats between 0, 1.0 so guarantees
			// the first node checked will be less than initial value
			weight = 2.0
			c      = 0
		)
		for n := range ends {
			f.current = n
			c = f.updateAvailable()
			if c == 0 {
				// nothing else available for this node, remove it
				delete(ends, n)
				continue
			}
			for i := 0; i < c; i++ {
				if w := n.Weights[f.available[i]]; w < weight {
					next = f.available[i]
					parent = n
				}
			}
		}
		if len(ends) == 0 {
			break
		}
		// set current to parent for animation
		f.current = parent
		next.Parent = parent
		next.Visited = true
		ends[next] = struct{}{}
		animate()
	}
}
