package main

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/rendicott/gree"
)

func main() {
	root := gree.NewNode("root")
	root.NewChild("sometimes you just have really long").
		NewChild("node names and you can't figure out a way to make them").
		NewChild("shorter and you also want them to have really nice").
		NewChild("colors so you can see them in your tree").SetColorRed().
		NewChild("but you just can't figure out how you got into this dystopian situation in the first place").
		NewChild("and for that situation the package helps you by truncating really long strings that will spill over the border like thiiiiiiiiiiiiiiiiiiiiiiiiiiiiiis").SetColor(color.BgBlue)
	root.NewChild("child2")
	fmt.Println(root.DrawOptions(&gree.DrawInput{
		Border:  true,
		Padding: "           ",
	}))
}
