package node

import (
	"fmt"
)

type Node struct {
	X       int
	Y       int
	IsStart bool
	IsEnd   bool
	Visited bool
	Parent  *Node
	// Weights
	Weights map[*Node]float64
}

func (n Node) IsLeft(n2 *Node) bool {
	return n.X < n2.X
}

func (n Node) IsAbove(n2 *Node) bool {
	return n.Y < n2.Y
}

func (n Node) String() string {
	return fmt.Sprintf("(%d, %d)", n.X, n.Y)
}
