package main

import (
	"fmt"

	"github.com/rendicott/gree"
)

func main() {
	// yellow := color.New(color.FgYellow).SprintFunc()
	// red := color.New(color.FgRed).SprintFunc()
	// magenta := color.New(color.FgMagenta).SprintFunc()
	a := gree.NewNode("father")
	for i := 0; i < 2; i++ {
		//description := magenta(fmt.Sprintf("child%d", i))
		description := fmt.Sprintf("child%d", i)
		b := gree.NewNode(description)
		for j := 0; j < 3; j++ {
			//description = yellow(fmt.Sprintf("%s%d", "grandchild", j))
			description = fmt.Sprintf("%s%d", "grandchild", j)
			b.NewChild(description)
		}
		a.AddChild(b)
	}
	a.NewChild("one").NewChild("two").NewChild("three").NewChild("four").NewChild("five")

	gen := a.GetGeneration(2)
	for i, c := range gen {
		c.NewChild(fmt.Sprintf("%s%d", "apple", i))
	}
	gen = a.GetGeneration(3)
	for i, c := range gen {
		c.NewChild(fmt.Sprintf("%s%d", "oranges", i))
	}
	_ = a.GetGeneration(2)[0].NewChild("bob")
	//err := a.SetPaddingAll("╳╳╳╳╳╳")
	// err := a.SetPaddingAll("===========")
	// if err != nil {
	// 	fmt.Println(err)
	// }
	di := gree.DrawInput{
		Border:  true,
		Padding: "╳╳╳╳╳╳╳╳╳╳╳╳",
		// Padding: "            ",
		Debug: true,
	}
	fmt.Println(a.DrawOptions(&di))
}
