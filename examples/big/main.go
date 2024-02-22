package main

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/rendicott/gree"
)

func main() {
	a := gree.NewNode("root")
	// add a bunch of children and grandchildren
	for i := 0; i < 2; i++ {
		description := fmt.Sprintf("child%d", i)
		b := gree.NewNode(description)
		b.SetColorMagenta()
		for j := 0; j < 3; j++ {
			description = fmt.Sprintf("%s%d", "grandchild", j)
			newGrandchild := b.NewChild(description).SetColorYellow()
			fmt.Printf("added newGrandchild %s at depth %d\n", newGrandchild.String(), newGrandchild.GetDepth())
		}
		a.AddChild(b)
	}
	// add a new lineage of children and grandchildren in one line
	a.NewChild("one").NewChild("two").NewChild("three").NewChild("four").NewChild("five")

	// check on the depth of one of the previously created granchildren that used to be depth 1
	gc := a.GetChild(0).GetChild(0)
	fmt.Printf("previously added grandchild '%s' now has depth %d\n", gc.String(), gc.GetDepth())

	// retrieve a generation by index and add children to each of the children in that generation
	gen := a.GetGeneration(2)
	for i, c := range gen {
		c.NewChild(fmt.Sprintf("%s%d", "apple", i)).SetColorRed()
	}

	// now grab that newly added generation and add children to each of the children in that generation
	// setting their color to a custom layered fatih/color attribute
	gen = a.GetGeneration(3)
	for i, c := range gen {
		c.NewChild(fmt.Sprintf("%s%d", "oranges", i)).SetColor(color.BgBlue).SetColor(color.FgHiRed)
	}

	// randombly target a generation and index and add a child to it
	a.GetGeneration(2)[0].NewChild("bob")

	// use every custom input option
	di := gree.DrawInput{
		Border:  true,
		Padding: "  ",
		Debug:   true,
	}
	fmt.Println(a.DrawOptions(&di))
}
