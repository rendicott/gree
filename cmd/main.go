package main

import (
	"fmt"

	"github.com/rendicott/gree"
)

func main() {
	a := gree.NewNode("root")
	a.NewChild("child1").NewChild("grandchild1")
	a.NewChild("child2").NewChild("grandchild2")
	a.NewChild("child3").NewChild("grandchild3")
	a.NewChild("child4")
	a.NewChild("child5")
	// add child to 1st grandchild of 2nd child of root
	a.GetChild(1).GetChild(0).NewChild("whoops")

	// integrate a new tree
	b := gree.NewNode("extended")
	b.NewChild("puppy1").NewChild("grandpuppy1")
	b.NewChild("puppy2").NewChild("grandpuppy2")

	a.GetChild(2).AddChild(&b)
	// a := gree.NewNode("root")
	// a0 := gree.NewNode("child1")
	// a1 := gree.NewNode("child2")
	// a2 := gree.NewNode("child3")
	// b := gree.NewNode("grandchild1")
	// b0 := gree.NewNode("grandchild2")
	// b1 := gree.NewNode("grandchild3")
	// b2 := gree.NewNode("grandchild4")
	// b3 := gree.NewNode("grandchild5")
	// c := gree.NewNode("greatgrandchild1")
	// c1 := gree.NewNode("greatgrandchild2")
	// b0.AddChild(c)
	// a0.AddChild(b)
	// a0.AddChild(b0)
	// b1.AddChild(c1)
	// a0.AddChild(b1)
	// a1.AddChild(b2)
	// a.AddChild(a0)
	// a.AddChild(a1)
	// a2.AddChild(b3)
	// a.AddChild(a2)
	fmt.Println(a.Draw())
}
