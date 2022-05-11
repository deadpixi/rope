package rope

import (
	"testing"
)

func expectString(a, b string, t *testing.T) {
	if a != b {
		t.Fatalf("expected '%v', got '%v'", a, b)
	}
}

func expectInt(a, b int, t *testing.T) {
	if a != b {
		t.Fatalf("expected %v, got %v", a, b)
	}
}

func TestConcat(t *testing.T) {
	rope := NewString("foo")
	expectString("foo", rope.String(), t)

	rope = rope.ConcatString("bar")
	expectString("foobar", rope.String(), t)

	expectInt(6, rope.Length(), t)
}

func TestInsert(t *testing.T) {
	rope := NewString("hello")
	rope = rope.InsertString(rope.Length(), "world")
	rope = rope.InsertString(5, ", ")

	expectString("hello, world", rope.String(), t)

	rope = rope.InsertString(rope.Length(), "!")
	expectString("hello, world!", rope.String(), t)
}

func TestSplit(t *testing.T) {
	rope := NewString("how now")
	left, right := rope.Split(3)
	expectString("how", left.String(), t)
	expectString(" now", right.String(), t)
}

func TestDelete(t *testing.T) {
	rope := NewString("how now brown cow")
	rope = rope.Delete(8, 6)

	expectString("how now cow", rope.String(), t)
}

func TestBalance(t *testing.T) {
	rope := NewString("hello")
	for i := 0; i < 32; i++ {
		rope = &Rope{
			length: 5,
			depth:  i + 1,
			left:   rope,
			right:  New(),
		}
	}

	if rope.isBalanced() {
		t.Fatalf("expected rope to be unbalanced")
	}

	rope = rebalance(rope)

	if !rope.isBalanced() {
		t.Fatalf("expected rope to be balanced")
	}

	expectString("hello", rope.String(), t)
	expectInt(5, rope.Length(), t)
}

func TestBigString(t *testing.T) {
	rope := New()
	for i := 0; i < 1048576/32; i++ {
		rope = rope.ConcatString("aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa")
	}

	if !rope.isBalanced() {
		t.Fatalf("expected rope to be balanced")
	}

	for i := 1; i < 11; i++ {
		rope = rope.InsertString(1048576/i, "foo")
	}

	if !rope.isBalanced() {
		t.Fatalf("expected rope to be balanced")
	}
}
