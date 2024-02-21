package gree

import (
	"fmt"
	"os"
	"strings"
	"testing"
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
	linesExpected := strings.Split(expected, "\n")
	linesGot := strings.Split(got, "\n")
	for i := 0; i >= len(linesExpected); i++ {
		lineGot := ""
		if len(linesGot) >= i {
			lineGot = linesGot[i]
		}
		if linesExpected[i] != lineGot {
			t.Errorf("output does not match expected in testfile %s\n", testfile)
			t.Errorf("line %d, expected '%s', got '%s'", i, linesExpected[i], lineGot)
		}
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
		"root":        0,
		"child1":      1,
		"child2":      1,
		"child3":      1,
		"grandchild1": 2,
	}
	for k, v := range expected {
		if got[k] != v {
			t.Errorf("expected '%s=%d', got '%d'",
				k, v,
				got[k],
			)
		}
	}
}

func TestDepthSimple(t *testing.T) {
	a := NewNode("root")
	b := NewNode("child1")
	c := NewNode("grandchild1")
	b.AddChild(c)
	a.AddChild(b)
	expected := 2
	got := c.GetDepth()
	if got != expected {
		t.Errorf("got '%d', expected '%d'\n", got, expected)
	}
}

func TestMaxWidth(t *testing.T) {
	a := NewNode("root")
	a.NewChild("child1")
	a.NewChild("child2")
	a.NewChild("child3").NewChild("grandchild1")
	a.Draw()
	expected := 18
	got := a.getDescMaxWidth()
	if got != expected {
		fmt.Println(a.Draw())
		t.Errorf("expected %d, got %d", expected, got)
	}
}

func TestGetGeneration(t *testing.T) {
	a := NewNode("root")
	a.NewChild("child1")
	a.NewChild("child2")
	a.NewChild("child3").NewChild("grandchild1")
	got := a.GetGeneration(1)
	expected := 3
	if len(got) != expected {
		t.Errorf("expected %d, got %d", expected, len(got))
	}
}
