// Package gree provides a Node struct to which
// children can be retrieved and added. Calling
// the Draw() method on a Node returns the 'tree'
// like string representation of the Node and its
// children
//
// Example:
//
//  func main() {
//      a := gree.NewNode("root")
//	a.NewChild("child1")
//      a.NewChild("child2").NewChild("grandchild1")
//	a.Draw()
//  }
//
// Displays
//
//  root
//      ├── child1
//      └── child2
//          └── grandchild1
//
// 
package gree

import (
	"fmt"
)

// Node contains methods for adding/retrieving children
// and rendering a tree drawing.
type Node struct {
	parent   *Node
	children []*Node
	contents string
}

// GetChild returns a pointer to the y'th child
// of the Node. If the y'th child does not exist
// a nil pointer is returned.
func (n *Node) GetChild(y int) (dc *Node) {
	for i, c := range n.children {
		if y == i {
			return c
		}
	}
	return nil
}

// NewNode returns a new node with contents of
// the passed string.
func NewNode(contents string) Node {
	return Node{
		contents: contents,
	}
}

// String() satisfies the Stringer interface
func (n Node) String() string {
	return n.contents
}

// NewChild adds a child with contents of the passed
// string to this Node's children. It returns the pointer
// to the new Node. This can be discarded or used for chaining
// methods in literals (e.g., a.NewChild("foo").NewChild("bar"))
func (n *Node) NewChild(contents string) *Node {
	nn := Node{contents: contents}
	n.AddChild(&nn)
	return &nn
}

// AddChild adds the given Node to the children
// of the current Node
func (n *Node) AddChild(nc *Node) {
	nc.parent = n
	n.children = append(n.children, nc)
}

// Draw returns a string of the rendered tree for this
// Node as if this node is root
func (n Node) Draw() string {
	return n.draw("", "", false, false, false, false)
}

// draw is meant to be a recursive function passing knowledge about parent relationships
// down as function args through iterations.
func (n *Node) draw(rendering, padding string, amSibling, amLastSibling, parentIsSibling, parentIsLastSibling bool) string {
	var decorator string
	if amLastSibling {
		decorator = "└── "
	} else {
		decorator = "├── "
	}
	if n.parent != nil {
		if parentIsSibling && !parentIsLastSibling {
			padding += "│   "
		} else {
			padding += "    "
		}
		rendering += fmt.Sprintf("\n%s%s%s", padding, decorator, n)
	} else {
		rendering = n.String()
	}
	size := len(n.children)
	for i, child := range n.children {
		as := true            // am sibling
		als := false          // am last sibling
		pis := amSibling      // parent is sibling
		pils := amLastSibling // parent is last sibling
		if i == (size - 1) {  // last element
			als = true
		}
		rendering = child.draw(rendering, padding, as, als, pis, pils)
	}
	return rendering
}
