package gree

import (
	"fmt"
)

type Node struct {
	parent *Node
	children []*Node
	contents string
}

func NewNode(contents string) Node {
	return Node{
		contents: contents,
	}
}

func (n Node) String() string {
	return n.contents
}

func (n *Node) NewChild(contents string) {
	nn := Node{contents: contents}
	n.AddChild(nn)
}

func (n *Node) AddChild(nc Node) {
	nc.parent = n
	n.children = append(n.children, &nc)
}

func (n Node) Draw() string {
	return n.draw("", "", 0, false, false, false, false)
}

func (n *Node) draw(rendering, padding string, decoratorType int, amSibling, amLastSibling, parentIsSibling, parentIsLastSibling bool) string {
	var decorator string
	switch decoratorType {
		case 1:
			decorator = "├── "
		case 2:
			decorator = "└── "
		default:
			decorator = ""
	}
	if n.parent != nil {
		if parentIsSibling && !parentIsLastSibling {
			padding += "│   "
		} else {
			padding += "    "
		}
		rendering += fmt.Sprintf("\n%s%s%s", padding, decorator, n)
	} else {
		rendering = fmt.Sprintf("%s", n)
	}
	size := len(n.children)
	for i, child := range n.children{
		dt := 0 // decorator type
		as := true // am sibling
		als := false // am last sibling
		pis := false // parent is sibling
		pils := amLastSibling // parent is last sibling
		switch i {
			case size - 1: // last element
				als = true
				dt = 2
			default:
				dt = 1
		}
		if amSibling {
			pis = true
		}
		if amLastSibling {
			als = true
		}
		rendering = fmt.Sprintf("%s", child.draw(rendering, padding, dt, as, als, pis, pils))
	}
	return rendering
}

