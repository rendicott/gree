// Package gree provides a Node struct to which
// children can be retrieved and added. Calling
// the Draw() method on a Node returns the 'tree'
// like string representation of the Node and its
// children
//
// Example:
//
//	 func main() {
//	     a := gree.NewNode("root")
//		a.NewChild("child1")
//	     a.NewChild("child2").NewChild("grandchild1")
//		a.Draw()
//	 }
//
// Displays
//
//	root
//	    ├── child1
//	    └── child2
//	        └── grandchild1
package gree

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"unicode/utf8"
)

// Node contains methods for adding/retrieving children
// and rendering a tree drawing.
type Node struct {
	parent   *Node
	children []*Node

	// Contents is the string identifier for thise node
	// and is what will be displayed
	contents string
	// Padding determines how many spaces for
	// each indentation, defaults to "   " (3 spaces)
	padding             string
	depth               int
	amLastSibling       bool
	amSibling           bool
	parentIsSibling     bool
	parentIsLastSibling bool
	decorator           string
	repr                string
	isRoot              bool
	aligned             bool
}

// Debug returns a string with all of this node's
// unexposed properties.
func (n *Node) Debug() string {
	var parent string
	if n.parent != nil {
		parent = n.parent.contents
	}
	return fmt.Sprintf(
		"%s: %d\n"+
			"%s: %v\n"+
			"%s: %v\n"+
			"%s: %v\n"+
			"%s: %v\n"+
			"%s: '%s'\n"+
			"%s: '%s'\n"+
			"%s: '%s'\n"+
			"%s: '%s'\n"+
			"%s: '%s'\n",
		"depth", n.depth,
		"amLastSibling", n.amLastSibling,
		"amSibling", n.amSibling,
		"parentIsSibling", n.parentIsSibling,
		"parentIsLastSibling", n.parentIsLastSibling,
		"padding", strings.Replace(n.padding, " ", "-", -1),
		"decorator", strings.Replace(n.decorator, " ", "-", -1),
		"repr", strings.Replace(n.repr, " ", "-", -1),
		"contents", n.contents,
		"parent", parent,
	)
}

type collector struct {
	results []*Node
}

func (c *collector) add(n *Node) {
	c.results = append(c.results, n)
}

// GetDepth returns the node's depth. Only valid after
// a GetAllDescendents call has been run
func (n *Node) GetDepth() int {
	return n.depth
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

func (n *Node) getDescHeight() int {
	return len(n.GetAllDescendents())
}

func (n *Node) getDescMaxWidth() (max int) {
	// first have to draw before getDescMaxWidth works properly, yuck
	border := false
	_ = n.draw("", "", 0, border, 0)
	all := n.GetAllDescendents()
	for _, desc := range all {
		runeCount := utf8.RuneCountInString(desc.repr)
		//runeCountDec := utf8.RuneCountInString(desc.decorator)
		//fmt.Printf("'%s' is len %d\n", rend, runeCount)
		//fmt.Printf("decorator '%s' is len %d\n", desc.decorator, runeCountDec)
		if runeCount > max {
			max = runeCount
		}
	}
	return max
}

// NewNode returns a new node with contents of
// the passed string.
func NewNode(contents string) Node {
	n := Node{}
	n.SetContents(contents)
	n.setPadding("   ")
	return n
}

// String() satisfies the Stringer interface
func (n Node) String() string {
	return n.contents
}

// SetContents sets new contents for this node
func (n *Node) SetContents(newContents string) {
	n.contents = newContents
}

// setPadding sets new padding for this node. Warning:
// setting padding for individual nodes can cause odd
// display characteristics.
func (n *Node) setPadding(padding string) error {
	if len(padding) < 1 {
		return errors.New("padding must be greater than len(1)")
	}
	n.padding = padding
	return nil
}

// SetPadding sets new padding for this node
// and all of it's descendents
func (n *Node) SetPaddingAll(padding string) (err error) {
	n.padding = padding
	for _, node := range n.GetAllDescendents() {
		err = node.setPadding(padding)
		if err != nil {
			return err
		}
	}
	return err
}

// GetAllDescendents gets all descendents of this node
// and returns a slice of pointers. Useful
// for updating them.
func (n *Node) GetAllDescendents() (all []*Node) {
	all = append(all, n.children...)
	for _, child := range n.children {
		all = append(all, child.GetAllDescendents()...)
	}
	return all
}

// NewChild adds a child with contents of the passed
// string to this Node's children. It returns the pointer
// to the new Node. This can be discarded or used for chaining
// methods in literals (e.g., a.NewChild("foo").NewChild("bar"))
func (n *Node) NewChild(contents string) *Node {
	nn := Node{}
	nn.SetContents(contents)
	nn.setPadding("   ")
	n.AddChild(&nn)
	return &nn
}

// AddChild adds the given Node to the children
// of the current Node
func (n *Node) AddChild(nc *Node) {
	nc.parent = n
	nc.depth = n.depth + 1
	n.children = append(n.children, nc)
	n.updateDepths()
}

func (n *Node) updateDepths() {
	newDepth := 0
	parent := n.parent
	for parent != nil {
		newDepth += 1
		parent = parent.parent
	}
	n.depth = newDepth
	for _, child := range n.children {
		child.updateDepths()
	}
}

// func (n Node) DrawWrap(di *DrawInput) (rendering string) {
// 	_ = n.draw("", "", 0, false)
// 	n.drawn = true
// 	maxWidth := n.getDescMaxWidth()
// 	fmt.Printf("overall maxWidth = %d\n", maxWidth)
// 	di.Border = false
// 	for _, child := range n.children {
// 		fmt.Println(child.DrawOptions(di))
// 		child.drawn = true
// 		fmt.Printf("child '%s' maxWidth = %d\n", child, child.getDescMaxWidth())
// 		for _, c := range child.children {
// 			fmt.Println(c.DrawOptions(di))
// 			c.drawn = true
// 			fmt.Printf("grandchild '%s' maxWidth = %d\n", c, c.getDescMaxWidth())
// 		}
// 	}
// 	return rendering
// }

// DrawInput holds input options for the DrawOptions method
type DrawInput struct {
	Border     bool   // whether or not to draw a border
	Debug      bool   // whether or not to add debug info to output
	Padding    string // rendered padding for this and child nodes
	Align      bool   // whether to align all labels
	startIndex int
}

// Draw sets default input options and returns a string
// of the rendered tree for this Node as if this node is root
func (n Node) Draw() (rendering string) {
	di := DrawInput{
		Border:  false,
		Debug:   false,
		Padding: n.padding,
		Align:   false,
	}
	rendering = n.DrawOptions(&di)
	return rendering
}

func (n Node) getLongestNodeLabel() (max int) {
	descs := n.GetAllDescendents()
	for _, desc := range descs {
		lenLabel := utf8.RuneCountInString(desc.String())
		if lenLabel > max {
			max = lenLabel
		}
	}
	return max
}

func numRunesToStartIdex(rendering string, startIndex int) (count int) {
	lines := strings.Split(rendering, "\n")
	for _, line := range lines {
		if len(line) >= startIndex {
			toStart := line[0:startIndex]
			tcount := utf8.RuneCountInString(toStart)
			if tcount > count {
				count = tcount
			}
		}
	}
	return count
}

// DrawOptions takes input options and returns a string
// of the rendered tree for this Node as if this node is root
func (n Node) DrawOptions(di *DrawInput) (rendering string) {
	if di.Padding != "" {
		n.SetPaddingAll(di.Padding)
	}
	debug := di.Debug
	border := di.Border
	// set this node as root
	n.isRoot = true
	n.relate(false, false, false, false)
	maxWidth := n.getDescMaxWidth()
	tempRendering := n.draw("", "", maxWidth, border, 0)
	//lines := strings.Split(tempRendering, "\n")
	//// remove paddingn from 1st generation
	//for _, line := range lines {
	//	rendering += " " + strings.TrimPrefix(line, n.padding) + "\n"
	//}
	if di.Align && !n.aligned {
		// need to find rightmost diagram char and then redraw
		startIndex := getAlignCol(tempRendering)
		if border {
			startIndex++
		}
		// adjust maxwidth by longest node label - len(longestdesc)
		lenLongestLabel := n.getLongestNodeLabel()
		fmt.Println(lenLongestLabel)
		maxWidth = numRunesToStartIdex(tempRendering, startIndex) + lenLongestLabel
		tempRendering = n.draw("", "", maxWidth, border, startIndex)
	}
	if border {
		topBorder := "┏" + strings.Repeat("━", maxWidth-1) + "┓"
		botBorder := "┗" + strings.Repeat("━", maxWidth-1) + "┛"
		rendering = fmt.Sprintf("%s%s\n%s", topBorder, tempRendering, botBorder)
	} else {
		rendering = tempRendering
	}
	if debug {
		rendering += drawRuler(maxWidth)
	}
	return rendering
}

func findEndOfDiagramLine(line string) int {
	runes := []rune("└━")
	controlRune := runes[0]
	horoRune := runes[1]
	spaceByte := []byte(" ")[0]
	lineLength := utf8.RuneCountInString(line)
	for i, char := range line {
		if char == controlRune {
			// fmt.Printf("hit controlrune '%v' at index %d\n", char, i)
			for j := i; j <= lineLength; j++ {
				// fmt.Printf("line[%d] = '%v'\n", j, rune(line[j]))
				if rune(line[j]) == horoRune {
					// fmt.Printf("skipping horoRune '%v' at index '%d'\n", horoRune, j)
					continue
				} else if line[j] == spaceByte {
					return j
				}
			}
		}
	}
	return 0
}

func getAlignCol(rendering string) (max int) {
	lines := strings.Split(rendering, "\n")
	for _, line := range lines {
		endIndex := findEndOfDiagramLine(line)
		// fmt.Printf("got endInex = %d for line '%s'\n", endIndex, line)
		if endIndex > max {
			max = endIndex
		}
	}
	return max
}

// drawRuler adds a ruler with column identifiers
// every 5 ticks. It tries to keep labels lined
// up with tick marks
func drawRuler(maxWidth int) (ruler string) {
	ruler += "\n"
	for i := 0; i <= maxWidth; i++ {
		if i%5 == 0 {
			ruler += "|"
		} else {
			ruler += "."
		}
	}
	ruler += "\n"
	var skipCount int
	for i := 0; i <= maxWidth; i++ {
		if i%5 == 0 {
			label := strconv.Itoa(i)
			labelWidth := utf8.RuneCountInString(label)
			skipCount += labelWidth - 1
			ruler += label
		} else {
			if skipCount == 0 {
				ruler += " "
			} else {
				skipCount--
			}
		}
	}
	ruler += "\n"
	return ruler
}

// draw builds the rendering string rcursively
func (n *Node) draw(rendering, padding string, maxWidth int, border bool, startIndex int) string {
	var decorator string
	horo := "─"
	if n.amLastSibling {
		decorator = "└" + strings.Repeat(horo, len(n.padding)-1)
	} else {
		decorator = "├" + strings.Repeat(horo, len(n.padding)-1)
	}
	n.decorator = decorator
	var repr string
	if border {
		repr = "┃"
	}
	repr += n.String()
	// grab a fillchar from the user defined padding instead of
	// just always using space
	var fillChar string
	if len(n.padding) > 0 {
		fillChar = firstRuneChar(n.padding)
	} else {
		fillChar = " "
	}
	if n.parent != nil && !n.isRoot {
		if n.depth > 1 {
			if n.parentIsSibling && !n.parentIsLastSibling {
				padding += "│" + n.padding
			} else if n.parent.isRoot {
				padding += ""
			} else {
				padding += fillChar + n.padding
			}
		}
		if border {
			repr = fmt.Sprintf("┃%s%s", padding, decorator)
		} else {
			repr = fmt.Sprintf("%s%s", padding, decorator)
		}
	}
	// pad for align
	if len(repr) < startIndex && startIndex != 0 && !n.isRoot {
		diff := startIndex - len(repr)
		repr = fmt.Sprintf("%s%s", repr, strings.Repeat(horo, diff))
	}
	// finally, add contents of node text
	if !n.isRoot { // indicates root
		repr += " " + n.String()
	}
	currWidth := utf8.RuneCountInString(repr)
	if currWidth < maxWidth {
		fill := maxWidth - currWidth
		repr += strings.Repeat(fillChar, fill)
	}
	if border {
		repr += "┃"
	}
	rendering += fmt.Sprintf("\n%s", repr)
	n.repr = repr
	for _, child := range n.children {
		rendering = child.draw(rendering, padding, maxWidth, border, startIndex)
	}
	return rendering
}

func firstRuneChar(s string) (char string) {
	for i, w := 0, 0; i < len(s); i += w {
		runeValue, _ := utf8.DecodeRuneInString(s[i:])
		return string(runeValue)
	}
	return " "
}

// relate is meant to be a recursive function passing knowledge about parent relationships
// it sets node properties to be used later for drawing purposes
func (n *Node) relate(amSibling, amLastSibling, parentIsSibling, parentIsLastSibling bool) {
	n.amLastSibling = amLastSibling
	n.amSibling = amSibling
	n.parentIsLastSibling = parentIsLastSibling
	n.parentIsSibling = parentIsSibling
	size := len(n.children)
	for i, child := range n.children {
		as := true            // am sibling
		als := false          // am last sibling
		pis := amSibling      // parent is sibling
		pils := amLastSibling // parent is last sibling
		if i == (size - 1) {  // last element
			als = true
		}
		child.relate(as, als, pis, pils)
	}
}

func (n *Node) dive(depth int) int {
	if len(n.children) > 0 {
		depth += 1
		for _, child := range n.children {
			var d int
			if d = child.dive(depth); d > depth {
				depth = d
			}
		}
	}
	return depth
}

func (n *Node) diveRetrieve(depth, desired int, col *collector) {
	// if desired is -1 then we'll just set depth and
	// add ourselves to collector
	if desired == -1 {
		nn := NewNode(n.contents)
		nn.setPadding(n.padding)
		nn.parent = n.parent
		nn.children = append(nn.children, n.children...)
		nn.depth = depth
		col.add(&nn)
	}

	// if this node's children are the desired depth then
	// add them to the collector and return
	if (depth+1 == desired) && (col != nil) && len(n.children) != 0 {
		for _, c := range n.children {
			col.add(c)
		}
		return
	}

	// otherwise, dig deeper
	if len(n.children) > 0 {
		depth += 1
		for _, child := range n.children {
			child.diveRetrieve(depth, desired, col)
		}
	}
}

// NumChildren returns the number of children
// this node has
func (n *Node) NumChildren() int {
	return len(n.children)
}

// GetGeneration gets all the children of the y'th
// generation of this node
func (n *Node) GetGeneration(y int) []*Node {
	col := collector{}
	var depth int
	n.diveRetrieve(depth, y, &col)
	return col.results
}

// MaxDepth returns the maximum depth of descendents
// and child descendents
func (n *Node) MaxDepth() (maxDepth int) {
	for _, child := range n.children {
		depth := 1
		depth = child.dive(depth)
		if depth > maxDepth {
			maxDepth = depth
		}
	}
	return maxDepth
}
