package rope

import (
	"io"
	"strings"
)

const (
	// If a modification would result in a string node shorter than this,
	// we simply create a single leaf node.
	shortLength = 256

	// The maximum depth of the tree. 64 is far deeper than any tree will
	// ever be.
	maxDepth = 64
)

var fibonacci []int

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

var empty = NewString("")

func New() Rope {
	return empty
}

func NewString(s string) Rope {
	return leaf(s)
}

func NewReader(r Rope) *Reader {
	return &Reader{rope: r}
}

type Rope interface {
	// Append another rope to this rope, returning the new, concatenated rope.
	Append(Rope) Rope

	// Append a string to this rope, returning the new, concatenated rope.
	AppendString(string) Rope

	// Delete length bytes starting at offset, returning the new rope.
	Delete(offset, length int) Rope

	// Returns true if this rope is equal to other.
	Equal(other Rope) bool

	// Return the byte stored at index.
	Index(index int) byte

	// Insert rope other at the given index, returning the new rope.
	Insert(index int, other Rope) Rope

	// Insert string s at index i, returning the new rope.
	InsertString(i int, s string) Rope

	// Return the length of the rope, in bytes.
	Length() int

	// Read contents at offset.
	ReadAt(p []byte, off int64) (n int, err error)

	// Split the rope at index i, returning the left and right sides.
	Split(i int) (Rope, Rope)

	// Convert the rope to a string. Note that, for large ropes, this can
	// be an expensive operation.
	String() string

	isBalanced() bool
	rebalance() Rope
	treeDepth() int
	walk(w walker)
}

// WALKERS
type walker interface {
	visit(Rope)
}

// LEAF COLLECTION
// A walker that collects leaves.
// Note that we could just do this recursively, but by doing it this way
// we can avoid some pressure on the garbage collector compared to the
// naive approach, and also lets us use the existing walking infrastructure.
type leafCollector struct {
	leaves []leaf
}

func (walker *leafCollector) visit(node Rope) {
	if leaf, ok := node.(leaf); ok {
		walker.leaves = append(walker.leaves, leaf)
	}
}

// STRINGIFICATION
type stringifier struct {
	builder strings.Builder
}

func (walker *stringifier) visit(node Rope) {
	if leaf, ok := node.(leaf); ok {
		walker.builder.WriteString(leaf.String())
	}
}

// GENERIC IMPLEMENTATIONS
// Most of the algorithms in the Rope interface can be generically implemented
// in terms of other operations, at least for leaves and concatenation nodes.
// We factor these out here.
func delete(node Rope, offset, length int) Rope {
	left, right := node.Split(offset)
	_, newRight := right.Split(length)
	return left.Append(newRight).rebalance()
}

func equal(node, other Rope) bool {
	if node == other {
		return true
	}

	if node.Length() != other.Length() {
		return false
	}

	// FIXME - this could be made so much more efficient.
	for i := 0; i < node.Length(); i++ {
		if node.Index(i) != other.Index(i) {
			return false
		}
	}

	return true
}

func insert(at int, node, other Rope) Rope {
	if at == 0 {
		return other.Append(node).rebalance()
	}

	if at == node.Length() {
		return node.Append(other).rebalance()
	}

	left, right := node.Split(at)
	return left.Append(other).Append(right).rebalance()
}

func join(node, other Rope) Rope {
	switch {
	case node.Length() == 0:
		return other
	case other.Length() == 0:
		return node
	case node.Length()+other.Length() <= shortLength:
		return NewString(node.String() + other.String())
	default:
		depth := node.treeDepth()
		if other.treeDepth() > depth {
			depth = other.treeDepth()
		}
		return concat{
			length: node.Length() + other.Length(),
			depth:  depth + 1,
			left:   node,
			right:  other,
		}.rebalance()
	}
}

// LEAF NODES
type leaf string

func (node leaf) Append(other Rope) Rope {
	return join(node, other)
}

func (node leaf) AppendString(other string) Rope {
	return join(node, NewString(other))
}

func (node leaf) Delete(offset, length int) Rope {
	return delete(node, offset, length)
}

func (node leaf) Equal(other Rope) bool {
	return equal(node, other)
}

func (node leaf) Index(i int) byte {
	return node[i]
}

func (node leaf) Insert(i int, other Rope) Rope {
	return insert(i, node, other)
}

func (node leaf) InsertString(i int, other string) Rope {
	return insert(i, node, NewString(other))
}

func (node leaf) Length() int {
	return len(node)
}

func (node leaf) ReadAt(p []byte, off int64) (n int, err error) {
	length := int64(len(p))
	nodeLength := int64(len(node))

	if off+length >= nodeLength {
		err = io.EOF
	}

	n = copy(p, []byte(node[off:]))
	return
}

func (node leaf) Split(at int) (Rope, Rope) {
	return leaf(node[:at]), leaf(node[at:])
}

func (node leaf) String() string {
	return string(node)
}

func (node leaf) isBalanced() bool {
	return true
}

func (node leaf) rebalance() Rope {
	return node
}

func (node leaf) treeDepth() int {
	return 0
}

func (node leaf) walk(w walker) {
	w.visit(node)
}

// CONCAT NODES
type concat struct {
	contents      string
	length, depth int
	left, right   Rope
}

func (node concat) Append(other Rope) Rope {
	return join(node, other)
}

func (node concat) AppendString(other string) Rope {
	return join(node, NewString(other))
}

func (node concat) Delete(offset, length int) Rope {
	return delete(node, offset, length)
}

func (node concat) Equal(other Rope) bool {
	return equal(node, other)
}

func (node concat) Index(at int) byte {
	if at < node.left.Length() {
		return node.left.Index(at)
	}

	return node.right.Index(at - node.left.Length())
}

func (node concat) Insert(i int, other Rope) Rope {
	return insert(i, node, other)
}

func (node concat) InsertString(i int, other string) Rope {
	return insert(i, node, NewString(other))
}

func (node concat) Length() int {
	return node.length
}

func (node concat) ReadAt(p []byte, off int64) (n int, err error) {
	length := int64(len(p))
	nodeLength := int64(node.Length())

	if off+length >= nodeLength {
		err = io.EOF
	}

	if off < int64(node.left.Length()) {
		n, _ = node.left.ReadAt(p, off)
	}

	if n < len(p) {
		n2, _ := node.right.ReadAt(p[n:], 0)
		n += n2
	}
	return
}

func (node concat) Split(at int) (Rope, Rope) {
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

func (node concat) String() string {
	s := stringifier{}
	node.walk(&s)
	return s.builder.String()
}

func (node concat) isBalanced() bool {
	return node.depth < len(fibonacci)-2 && fibonacci[node.depth+2] <= node.Length()
}

func (node concat) rebalance() Rope {
	if node.isBalanced() {
		return node
	}

	walker := &leafCollector{}
	node.walk(walker)

	var merge func(leaves []leaf, start, end int) Rope
	merge = func(leaves []leaf, start, end int) Rope {
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

	return merge(walker.leaves, 0, len(walker.leaves))
}

func (node concat) treeDepth() int {
	return node.depth
}

func (node concat) walk(w walker) {
	node.left.walk(w)
	w.visit(node)
	node.right.walk(w)
}

// READERS
type Reader struct {
	rope   Rope
	offset int64
}

func (reader *Reader) Read(p []byte) (n int, err error) {
	n, err = reader.rope.ReadAt(p, reader.offset)
	reader.offset += int64(n)
	return
}
