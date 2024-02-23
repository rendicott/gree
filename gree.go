// Package gree provides a Node struct to which
// children can be retrieved and added. Calling
// the Draw() method on a Node returns the 'tree'
// like string representation of the Node and its
// children
//
// Example:
//
//	 func main() {
//		 a := gree.NewNode("root")
//		 a.NewChild("child1")
//		 a.NewChild("child2").NewChild("grandchild1")
//		 fmt.Println(a.Draw())
//	 }
//
// Displays
//
//	root
//	├── child1
//	└── child2
//	    └── grandchild1
//
// The package provides many convenient methods for
// retrieving children by generation, getting descendent
// depth, setting display padding, and setting colors.
//
// The package also exposes the DrawOptions method for
// more fine grained control over the display.
//
// Any node from which the Draw*() methods are called
// will be considered the root node for display purposes.
package gree

import (
	"errors"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/fatih/color"
	"github.com/google/uuid"
)

// Node contains methods for adding/retrieving children
// and rendering a tree drawing.
type Node struct {
	parent   *Node
	lineage  []*Node // lineage is the parent and all of the parent's parents
	children []*Node
	id       uuid.UUID

	// Contents is the string identifier for thise node
	// and is what will be displayed
	contents         string
	contentsColored  string
	colored          bool
	contentFontWidth int
	contentLength    int
	// Padding determines how many spaces for
	// each indentation, defaults to "   " (3 spaces)
	padding             string
	depth               int
	amLastSibling       bool
	amSibling           bool
	parentIsSibling     bool
	parentIsLastSibling bool
	parentIsRoot        bool
	isRoot              bool
	x1                  int
	x2                  int
	done                bool
	index               int
	count               counter
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

// GetID returns the UUID of the node.
// Useful for identifying unique nodes when
// many have the same contents.
func (n *Node) GetID() string {
	return n.id.String()
}

// setx1 sets the x1 property of this node and auto
// recalculates x2 based on the contents
func (n *Node) setx1(x int) {
	n.x1 = x
	n.x2 = n.x1 + utf8.RuneCountInString(n.contents)
	n.contentLength = utf8.RuneCountInString(n.contents)
}

// SetColorMagenta sets the color of the node to magenta
func (n *Node) SetColorMagenta() *Node {
	return n.SetColor(color.FgMagenta)
}

// SetColorGreen sets the color of the node to green
func (n *Node) SetColorYellow() *Node {
	return n.SetColor(color.FgYellow)
}

// SetColorRed sets the color of the node to red
func (n *Node) SetColorRed() *Node {
	return n.SetColor(color.FgRed)
}

// SetColor sets the color of the node to the passed fatih/color attribute
// Requires that the caller import fatih/color and reference their color.Attribute
func (n *Node) SetColor(fatihcolor color.Attribute) *Node {
	if n.contentsColored == "" {
		n.contentsColored = n.contents
	}
	n.contentsColored = color.New(fatihcolor).Sprint(n.contentsColored)
	n.colored = true
	return n
}

type collector struct {
	results []*Node
}

func (c *collector) add(n *Node) {
	c.results = append(c.results, n)
}

// GetDepth returns this node's depth. Depths are updated
// as nodes are added.
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

func (n *Node) relateAsRoot() {
	n.isRoot = true
	n.relate(&n.count, false, true, false, false, n)
}

func (n *Node) getDescMaxWidth() (max int) {
	// first have to relate before getDescMaxWidth works properly, yuck
	n.relateAsRoot()
	all := n.GetAllDescendents()
	for _, dec := range all {
		declen := dec.x2 + utf8.RuneCountInString(dec.padding)
		if declen > max {
			max = declen
		}
	}
	return max
}

// NewNode returns a new node with contents of
// the passed string. Please do not use color formatted
// strings and instead use the provided SetColor* methods.
func NewNode(contents string) *Node {
	n := Node{
		id: uuid.New(),
	}
	n.SetContents(contents)
	n.setPadding("   ")
	return &n
}

// String returns a string, satisfying the Stringer interface
func (n Node) String() string {
	return n.contents
}

// SetContents sets new contents for this node. Please
// do not use color formatted strings and instead use the provided SetColor* methods.
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
// and all of it's descendents.
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

const (
	blankUUID string = "00000000-0000-0000-0000-000000000000"
)

// NewChild adds a child with contents of the passed
// string to this Node's children. It returns the pointer
// to the new Node. This can be discarded or used for chaining
// methods in literals (e.g., a.NewChild("foo").NewChild("bar"))
//
// Please do not use color formatted strings and instead use the provided SetColor* methods.
func (n *Node) NewChild(contents string) *Node {
	if n.id.String() == blankUUID {
		n.id = uuid.New()
	}
	nn := n.AddChild(NewNode(contents))
	return nn
}

// AddChild adds the given Node to the children
// of the current Node
func (n *Node) AddChild(nc *Node) *Node {
	if n.id.String() == blankUUID {
		n.id = uuid.New()
	}
	nc.parent = n
	nc.depth = n.depth + 1
	n.children = append(n.children, nc)
	n.updateDepths()
	return nc
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
	Border  bool   // whether or not to draw a border
	Debug   bool   // whether or not to add debug info to output
	Padding string // rendered padding for this and child nodes
}

// Draw sets default input options and returns a string
// of the rendered tree for this Node as if this node is root
func (n *Node) Draw() (rendering string) {
	di := DrawInput{
		Border:  false,
		Debug:   false,
		Padding: n.padding,
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

func (n *Node) render(width int, border bool) (row *rrow) {
	if n.colored {
		width = width + (utf8.RuneCountInString(n.contentsColored) - utf8.RuneCountInString(n.contents))
	}
	row = newRrow(width)
	for x := 0; x <= width; x++ {
		if (x == 0 || x == width) && border {
			row.setRowI(x, vbar(), true)
		}
		for _, p := range n.lineage {
			if x == p.x1 {
				if !p.amLastSibling && !p.isRoot {
					row.setRowI(x, vbar(), false)
				}
			}
		}
		if x == n.x1 {
			if n.colored {
				row.appendString(x, n.genDecorator(0)+n.contentsColored)
			} else {
				row.appendString(x, n.genDecorator(0)+n.String())
			}
		} else {
			row.setRowI(x, n.padRune(), false)
		}
	}
	return row
}

func (n *Node) genDecorator(decLength int) string {
	if n.isRoot {
		return ""
	}
	length := utf8.RuneCountInString(n.padding) - 1
	if decLength != 0 {
		length = decLength
	}
	if n.amLastSibling && !n.isRoot {
		return sibCharLastS() + strings.Repeat(horos(), length) + " "
	} else {
		return sibCharS() + strings.Repeat(horos(), length) + " "
	}
}

func (n Node) padRune() rune {
	return []rune(firstRuneChar(n.padding))[0]
}

// maybe we could handle chars with greater font width later
func (n *Node) setFontWidth() {
	n.contentFontWidth = utf8.RuneCountInString(n.contents)
	// couldn't figure out a good way to do this
	// for _, runeValue := range n.String() {
	// 	if unicode.Is(unicode.Han, runeValue) {
	// 		n.contentFontWidth += 4
	// 	}
	//}
}

func genTopBorder(width int) string {
	top := fmt.Sprintf("┌%s┐", strings.Repeat(horos(), width-1))
	return top
}

func genBottomBorder(width int) string {
	bottom := fmt.Sprintf("└%s┘", strings.Repeat(horos(), width-1))
	return bottom
}

func (n *Node) shiftAllRight(amount int) {
	n.setx1(n.x1 + amount)
	for _, desc := range n.GetAllDescendents() {
		desc.setx1(desc.x1 + amount)
	}
}

// DrawOptions takes a DrawInput struct with desired parameters
// and returns the tree formatted string.
func (n *Node) DrawOptions(di *DrawInput) (rendering string) {
	if di.Padding != "" {
		n.SetPaddingAll(di.Padding)
	}
	n.relateAsRoot() // set key properties of nodes
	bmp := make(map[int][]rune)
	width := n.getDescMaxWidth()
	if di.Border {
		width += 3
		n.shiftAllRight(2)
	}
	desc := n.GetAllDescendents()
	// draw root first
	bmp[0] = n.render(width, di.Border).toRunes()
	// now draw descendents
	for i := 1; i <= len(desc); i++ {
		cn := desc[i-1]
		cn.setFontWidth()
		bmp[i] = cn.render(width, di.Border).toRunes()
	}
	// build string
	var pre strings.Builder
	if di.Border {
		pre.Write([]byte(genTopBorder(width)))
		pre.Write([]byte("\n"))
	}
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
	if di.Border {
		pre.Write([]byte(genBottomBorder(width)))
		pre.Write([]byte("\n"))
	}
	if di.Debug {
		pre.Write([]byte(drawRuler(width)))
	}
	rendering = pre.String()
	return rendering
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

func horos() string {
	return "─"
}

func sibCharLastS() string {
	return "└"
}

func sibCharS() string {
	return "├"
}

func cleanLineage(input []*Node) (output []*Node) {
	for _, n := range input {
		if n != nil {
			output = append(output, n)
		}
	}
	return output
}

// relate is meant to be a recursive function passing knowledge about parent relationships
// it sets node properties to be used later for drawing purposes
func (n *Node) relate(count *counter, amSibling, amLastSibling, parentIsSibling, parentIsLastSibling bool, parent *Node) {
	n.index = count.get()
	count.add()
	n.amLastSibling = amLastSibling
	n.amSibling = amSibling
	n.parentIsLastSibling = parentIsLastSibling
	n.parentIsSibling = parentIsSibling
	size := len(n.children)
	if parent != nil {
		if parent.isRoot {
			n.setx1(parent.x1)
			n.parentIsRoot = true
		} else {
			n.setx1(parent.x1 + utf8.RuneCountInString(n.padding) + 1)
		}
		n.lineage = make([]*Node, len(parent.lineage)+1)
		for _, ancestor := range parent.lineage {
			if ancestor != nil {
				n.lineage = append(n.lineage, ancestor)
			}
		}
		n.lineage = append(n.lineage, parent)
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
	n.lineage = cleanLineage(n.lineage)
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
