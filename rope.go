package rope

import (
	"strings"
)

const (
	shortLength = 256
	maxDepth    = 64
)

var fibonacci []int

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

var empty = &Rope{leaf: true}

func New() *Rope {
	return empty
}

func NewString(s string) *Rope {
	return &Rope{leaf: true, contents: s}
}

func init() {
	first := 0
	second := 1

	for c := 0; c < maxDepth+3; c++ {
		next := 0
		if c <= 1 {
			next = c
		} else {
			next = first + second
			first = second
			second = next
		}
		fibonacci = append(fibonacci, next)
	}
}

type Rope struct {
	leaf          bool
	contents      string
	length, depth int
	left, right   *Rope
}

// WALKERS
type walker interface {
	visit(*Rope)
}

// LEAF COLLECTION
// A walker that collects leaves.
// Note that we could just do this recursively, but by doing it this way
// we can avoid some pressure on the garbage collector compared to the
// naive approach, and also lets us use the existing walking infrastructure.
type leafCollector struct {
	leaves []*Rope
}

func (walker *leafCollector) visit(node *Rope) {
	if node.leaf {
		walker.leaves = append(walker.leaves, node)
	}
}

// GENERIC OPERATIONS
// These operations work regardless of the type of rope (leaf, function, etc).
func rebalance(rope *Rope) *Rope {
	if rope.isBalanced() {
		return rope
	}

	walker := &leafCollector{}
	rope.walk(walker)

	return merge(walker.leaves, 0, len(walker.leaves))
}

func merge(leaves []*Rope, start, end int) *Rope {
	length := end - start
	switch length {
	case 1:
		return leaves[start]
	case 2:
		return leaves[start].Append(leaves[start+1])
	default:
		mid := start + length/2
		return merge(leaves, start, mid).Append(merge(leaves, mid, end))
	}
}

func (node *Rope) Insert(other *Rope, at int) *Rope {
	if at == 0 {
		return other.Append(node)
	}

	if at == node.Length() {
		return node.Append(other)
	}

	left, right := node.Split(at)
	return rebalance(left.Append(other).Append(right))
}

func (node *Rope) InsertString(at int, s string) *Rope {
	return node.Insert(NewString(s), at)
}

func (node *Rope) Append(other *Rope) *Rope {
	switch {
	case node.Length() == 0:
		return other
	case other.Length() == 0:
		return node
	case node.Length()+other.Length() <= shortLength:
		return NewString(node.String() + other.String())
	default:
		return rebalance(&Rope{
			length: node.Length() + other.Length(),
			depth:  max(node.depth, other.depth) + 1,
			left:   node,
			right:  other,
		})
	}
}

func (node *Rope) AppendString(s string) *Rope {
	return node.Append(NewString(s))
}

func (node *Rope) Delete(at, length int) *Rope {
	left, right := node.Split(at)
	_, newRight := right.Split(length)
	return rebalance(left.Append(newRight))
}

func (node *Rope) Equal(other *Rope) bool {
	if node == other {
		return true
	}

	if node.Length() != other.Length() {
		return false
	}

	// FIXME - this could be made so much more efficient
	for i := 0; i < node.Length(); i++ {
		if node.Index(i) != other.Index(i) {
			return false
		}
	}

	return true
}

func (node *Rope) Length() int {
	if node.leaf {
		return len(node.contents)
	}
	return node.length
}

func (node *Rope) Index(at int) byte {
	if node.leaf {
		return node.contents[at]
	}

	if at < node.left.Length() {
		return node.left.Index(at)
	}

	return node.right.Index(at - node.left.Length())
}

func (node *Rope) String() string {
	return node.Report(0, node.Length())
}

func (node *Rope) Split(at int) (*Rope, *Rope) {
	if node.leaf {
		return NewString(node.contents[0:at]), NewString(node.contents[at:])
	}

	if at < node.left.Length() {
		left, right := node.left.Split(at)
		return left, right.Append(node.right)
	}

	if at > node.left.Length() {
		left, right := node.right.Split(at - node.left.Length())
		return node.left.Append(left), right
	}

	return node.left, node.right
}

func (node *Rope) Report(at, length int) string {
	if node.leaf {
		return node.contents[at : at+length]
	}

	var b strings.Builder
	if at < node.left.Length() {
		leftLength := length
		if leftLength > node.left.Length() {
			leftLength = node.left.Length()
		}
		b.WriteString(node.left.Report(at, leftLength))

		at = at - leftLength
		length = length - leftLength
	}

	if length > 0 {
		b.WriteString(node.right.Report(at, length))
	}

	return b.String()
}

func (node *Rope) isBalanced() bool {
	if node.depth >= len(fibonacci)-2 {
		return false
	}

	return node.leaf || fibonacci[node.depth+2] <= node.Length()
}

func (node *Rope) walk(w walker) {
	if node.left != nil {
		node.left.walk(w)
	}

	w.visit(node)

	if node.right != nil {
		node.right.walk(w)
	}
}
