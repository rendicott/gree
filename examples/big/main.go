package main

import (
	"fmt"

	"github.com/rendicott/gree"
)

func main() {
	a := gree.NewNode("root")
	for i := 0; i < 2; i++ {
		description := fmt.Sprintf("child%d", i)
		b := gree.NewNode(description)
		for j := 0; j < 3; j++ {
			description = fmt.Sprintf("grandchild%d", j)
			b.NewChild(description)
		}
		a.AddChild(&b)
	}
	a.NewChild("one").NewChild("two").NewChild("three").NewChild("four").NewChild("five")

	gen := a.GetGeneration(2)
	for i, c := range gen {
		c.NewChild(fmt.Sprintf("%s%d", "carrot", i))
	}
	gen = a.GetGeneration(3)
	for i, c := range gen {
		c.NewChild(fmt.Sprintf("%s%d", "banana", i))
	}
	_ = a.GetGeneration(2)[0].NewChild("bob")
	//err := a.SetPaddingAll("╳╳╳╳")
	di := gree.DrawInput{
		Border: true,
		//Padding: "          ",
		//Align: true,
	}
	fmt.Println(a.DrawOptions(&di))
}
