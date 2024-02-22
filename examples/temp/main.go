package main

import (
	"fmt"

	"github.com/rendicott/gree"
)

func main() {
	a := gree.NewNode("root")
	a.NewChild("child1")
	a.NewChild("child2").NewChild("grandchild1")
	fmt.Println(a.Draw())
}
