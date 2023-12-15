package gree

import (
	"testing"
	"os"
	"fmt"
)

func TestDrawSimple(t *testing.T) {
        a := NewNode("root")
        a.NewChild("child1")
        a.NewChild("child2")
        a.NewChild("child3").NewChild("grandchild1")
	got := a.Draw()
	testfile := "./testdata/TestDrawSimple.txt"
	dat, err := os.ReadFile(testfile)
	if err != nil {
		t.Errorf("error pulling expected from file '%s', error '%s'\n", testfile, err.Error())
	}
	expected := string(dat)
	if got != expected {
		fmt.Println(got)
		t.Errorf("output does not match expected in testfile %s\n", testfile)
	}
}

func TestDepth(t *testing.T) {
	a := NewNode("root")
        a.NewChild("child1")
        a.NewChild("child2")
        a.NewChild("child3").NewChild("grandchild1")
	nodes := a.GetAllDescendents()
	got := make(map[string]int)
	for _, gotNode := range nodes {
		got[gotNode.String()] = gotNode.GetDepth()
	}
	expected := map[string]int{
		"root": 0,
		"child1": 1,
		"child2": 1,
		"child3": 1,
		"grandchild1": 2,
	}
	for k, v := range expected {
		if got[k] != v {
			t.Errorf("expected '%s=%d', got '%d'",
				k,v,
				got[k],
			)
		}
	}
}

func TestDepthSimple(t *testing.T) {
	a := NewNode("root")
	b := NewNode("child1")
	c := NewNode("grandchild1")
	b.AddChild(&c)
	a.AddChild(&b)
	expected := 2
	got := c.GetDepth()
	if got != expected {
		t.Errorf("got '%d', expected '%d'\n", got, expected)
	}
}

func TestHeight(t *testing.T) {
	a := NewNode("root")
        a.NewChild("child1")
        a.NewChild("child2")
        a.NewChild("child3").NewChild("grandchild1")
	expected := 4
	got := a.getDescHeight()
	if got != expected {
		t.Errorf("expected %d, got %d", expected, got)
	}
}

func TestMaxWidth(t *testing.T) {
	a := NewNode("root")
        a.NewChild("child1")
        a.NewChild("child2")
        a.NewChild("child3").NewChild("grandchild1")
	a.Draw()
	expected := 20
	got := a.getDescMaxWidth()
	if got != expected {
		fmt.Println(a.Draw())
		t.Errorf("expected %d, got %d", expected, got)
	}
}
