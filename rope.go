package rope

import (
	"strings"
)

const (
	maxDepth    = 64
	maxLeafSize = 4096
)

type Rope struct {
	content       string
	length, depth int
	left, right   *Rope
}

func New() *Rope {
	return &Rope{}
}

func NewString(s string) *Rope {
	return &Rope{content: s, length: len(s)}
}

func (rope *Rope) Append(other *Rope) *Rope {
	switch {
	case rope.length == 0:
		return other
	case other.length == 0:
		return rope
	case rope.length+other.length <= maxLeafSize:
		return NewString(rope.String() + other.String())
	default:
		depth := rope.depth
		if other.depth > depth {
			depth = other.depth
		}
		return (&Rope{
			length: rope.length + other.length,
			depth:  depth + 1,
			left:   rope,
			right:  other,
		}).rebalance()
	}
}

func (rope *Rope) AppendString(other string) *Rope {
	return rope.Append(NewString(other))
}

func (rope *Rope) Delete(offset, length int) *Rope {
	if length == 0 || offset == rope.length {
		return rope
	}

	left, right := rope.Split(offset)
	_, newRight := right.Split(length)
	return left.Append(newRight).rebalance()
}

func (rope *Rope) Equal(other *Rope) bool {
	if rope == other {
		return true
	}

	if rope.length != other.length {
		return false
	}

	for i := 0; i < rope.length; i++ {
		if rope.Index(i) != other.Index(i) {
			return false
		}
	}

	return true
}

func (rope *Rope) Index(at int) byte {
	if rope.isLeaf() {
		return rope.content[at]
	}

	if at < rope.left.length {
		return rope.left.Index(at)
	}

	return rope.right.Index(at - rope.left.length)
}

func (rope *Rope) Insert(at int, other *Rope) *Rope {
	if at == 0 {
		return other.Append(rope)
	}

	if at == rope.length {
		return rope.Append(other)
	}

	left, right := rope.Split(at)
	return left.Append(other).Append(right).rebalance()
}

func (rope *Rope) InsertString(at int, other string) *Rope {
	return rope.Insert(at, NewString(other))
}

func (rope *Rope) Length() int {
	return rope.length
}

func (rope *Rope) Split(at int) (*Rope, *Rope) {
	if rope.isLeaf() {
		return NewString(rope.content[0:at]), NewString(rope.content[at:])
	}

	if at < rope.left.length {
		left, right := rope.left.Split(at)
		return left, right.Append(rope.right)
	}

	if at > rope.left.length {
		left, right := rope.right.Split(at - rope.left.length)
		return rope.left.Append(left), right
	}

	return rope.left, rope.right
}

func (rope *Rope) String() string {
	if rope.isLeaf() {
		return rope.content
	}

	var builder strings.Builder
	rope.walk(func(node *Rope) {
		if node.isLeaf() {
			builder.WriteString(node.content)
		}
	})

	return builder.String()
}

func (rope *Rope) walk(callback func(*Rope)) {
	if rope.isLeaf() {
		callback(rope)
	} else {
		rope.left.walk(callback)
		callback(rope)
		rope.right.walk(callback)
	}
}

func (rope *Rope) isBalanced() bool {
	if rope.depth >= len(fibonacci)-2 {
		return false
	}

	return rope.isLeaf() || fibonacci[rope.depth+2] <= rope.length
}

func (rope *Rope) isLeaf() bool {
	return rope.left == nil
}

func (rope *Rope) leafForOffset(at int) (*Rope, int) {
	if rope.isLeaf() {
		return rope, at
	}

	if at < rope.left.length {
		return rope.left.leafForOffset(at)
	}

	return rope.right.leafForOffset(at - rope.left.length)
}

func (rope *Rope) rebalance() *Rope {
	if rope.isBalanced() {
		return rope
	}

	var leaves []*Rope
	rope.walk(func(node *Rope) {
		if node.isLeaf() {
			leaves = append(leaves, node)
		}
	})

	return merge(leaves, 0, len(leaves))
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
