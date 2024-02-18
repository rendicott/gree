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
	"sort"
	"strconv"
	"strings"
	"unicode/utf8"
)

// Node contains methods for adding/retrieving children
// and rendering a tree drawing.
type Node struct {
	parent   *Node
	lineage  []*Node // lineage is the parent and all of the parent's parents
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
	parentIsRoot        bool
	decorator           string
	repr                string
	isRoot              bool
	aligned             bool
	x1                  int
	x2                  int
	done                bool
	index               int
	count               counter
	labelDrawn          bool
}

func (n *Node) showLineage() (repr string) {
	repr += n.String() + ": "
	for _, p := range n.lineage {
		if p != nil {
			repr += p.String() + ","
		}
	}
	return repr
}

type counter struct {
	count int
}

func (c *counter) add() {
	c.count++
}

func (c *counter) get() int {
	return c.count
}

// setx1 sets the x1 property of this node and auto
// recalculates x2 based on the contents
func (n *Node) setx1(x int) {
	n.x1 = x
	n.x2 = n.x1 + utf8.RuneCountInString(n.contents)
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

func (n *Node) relateAsRoot() {
	n.isRoot = true
	n.relate(&n.count, false, true, false, false, nil)
}

func (n *Node) getDescMaxWidth() (max int) {
	// first have to relate before getDescMaxWidth works properly, yuck
	n.relateAsRoot()
	all := n.GetAllDescendents()
	for _, desc := range all {
		lots := desc.x2 + utf8.RuneCountInString(desc.decorator) + utf8.RuneCountInString(desc.padding)
		if lots > max {
			max = lots
		}
	}
	return max
}

// NewNode returns a new node with contents of
// the passed string.
func NewNode(contents string) *Node {
	n := Node{}
	n.SetContents(contents)
	n.setPadding("   ")
	return &n
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
	sort.Slice(all, func(i, j int) bool {
		return all[i].index < all[j].index
	})
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

type rrow struct {
	contents map[int]rune
	width    int
}

func (r *rrow) setRowI(i int, ru rune, override bool) {
	if r.width >= i {
		if override && r.contents[i] != 0 {
			r.contents[i] = ru
		} else if r.contents[i] == 0 {
			r.contents[i] = ru
		}
	}
}

func (r *rrow) appendString(afterI int, s string) {
	istr := []rune(s)
	for i := 0; i <= r.width; i++ {
		if i == afterI {
			for j := 0; j < len(istr); j++ {
				r.setRowI(i+j, istr[j], false)
			}
			break
		}
	}
}

func (r rrow) toRunes() []rune {
	return []rune(r.str())
}

func (r rrow) str() string {
	var results []rune
	for i := 0; i <= r.width; i++ {
		results = append(results, r.contents[i])
	}
	return string(results)
}

func newRrow(width int) *rrow {
	nrr := rrow{
		contents: make(map[int]rune, width),
		width:    width,
	}
	return &nrr
}

func vbar() rune {
	return []rune("│")[0]
}

func (n *Node) render(width, rightAlign int) (row *rrow) {
	fmt.Printf("rendering node '%s' with rightAlign = %d\n", n.String(), rightAlign)
	row = newRrow(width)
	for x := 0; x <= width; x++ {
		for _, p := range n.lineage {
			if x == p.x1 {
				if !p.amLastSibling {
					row.setRowI(x, vbar(), false)
				}
			}
		}
		if rightAlign != 0 && x == rightAlign {
			row.appendString(x, n.decorator+n.String())
		} else if x == n.x1 && rightAlign == 0 {
			row.appendString(x, n.decorator+n.String())
		} else {
			row.setRowI(x, n.padRune(), false)
		}
	}
	return row
}

func (n Node) padRune() rune {
	return []rune(firstRuneChar(n.padding))[0]
}

func (n *Node) DrawOptions(di *DrawInput) (rendering string) {
	n.relateAsRoot() // set key properties of nodes
	width := n.getDescMaxWidth()
	bmp := make(map[int][]rune)
	desc := n.GetAllDescendents()
	longest := n.getLongestNodeLabel()
	var rightAlign int
	if di.Align {
		rightAlign = width - longest
	}
	// draw root first
	bmp[0] = n.render(width, 0).toRunes()
	// now draw descendents
	for i := 1; i <= len(desc); i++ {
		cn := desc[i-1]
		bmp[i] = cn.render(width, rightAlign).toRunes()
	}
	// build string
	var pre strings.Builder
	// order our map
	keys := make([]int, 0)
	for k, _ := range bmp {
		keys = append(keys, k)
	}
	sort.Ints(keys)
	for _, k := range keys {
		line := bmp[k]
		pre.Write([]byte(string(line)))
		pre.Write([]byte("\n"))
	}
	rendering = pre.String()
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

func firstRuneChar(s string) (char string) {
	for i, w := 0, 0; i < len(s); i += w {
		runeValue, _ := utf8.DecodeRuneInString(s[i:])
		return string(runeValue)
	}
	return " "
}

func horo() rune {
	return []rune(horos())[0]
}
func horos() string {
	return "─"
}

// relate is meant to be a recursive function passing knowledge about parent relationships
// it sets node properties to be used later for drawing purposes
func (n *Node) relate(count *counter, amSibling, amLastSibling, parentIsSibling, parentIsLastSibling bool, parent *Node) {
	if n.done {
		return
	}
	n.index = count.get()
	count.add()
	n.amLastSibling = amLastSibling
	n.amSibling = amSibling
	n.parentIsLastSibling = parentIsLastSibling
	n.parentIsSibling = parentIsSibling
	size := len(n.children)
	if !n.isRoot {
		if n.amLastSibling {
			n.decorator = "└" + strings.Repeat(horos(), len(n.padding)-1) + " "
		} else {
			n.decorator = "├" + strings.Repeat(horos(), len(n.padding)-1) + " "
		}
	}
	if parent != nil {
		if parent.isRoot {
			n.setx1(parent.x1)
		} else {
			n.setx1(parent.x1 + utf8.RuneCountInString(n.decorator))
		}
		n.lineage = append(n.parent.lineage, n.parent)
		if parent.isRoot {
			n.parentIsRoot = true
		}
	}
	for i, child := range n.children {
		as := true            // am sibling
		als := false          // am last sibling
		pis := amSibling      // parent is sibling
		pils := amLastSibling // parent is last sibling
		if i == (size - 1) {  // last element
			als = true
		}
		child.relate(count, as, als, pis, pils, n)
	}
	n.done = true
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
		col.add(nn)
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
