package rope

import (
	"bufio"
	"io"
	"strings"
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

func TestAppend(t *testing.T) {
	rope := NewString("foo")
	expectString("foo", rope.String(), t)

	rope = rope.AppendString("bar")
	expectString("foobar", rope.String(), t)

	expectInt(6, rope.Length(), t)
}

func TestInsert(t *testing.T) {
	rope := NewString("hello")
	rope = rope.InsertString(rope.Length(), "world").InsertString(5, ", ")

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

func TestSlice(t *testing.T) {
	rope := NewString("hello")
	expectString("ell", string(rope.Slice(1, 4)), t)

	rope1 := NewString("hel")
	rope2 := NewString("lo")
	rope3 := Rope{
		depth:  1,
		length: rope1.length + rope2.length,
		left:   &rope1,
		right:  &rope2,
	}

	expectString("hello", rope3.String(), t)
	expectString("ell", string(rope3.Slice(1, 4)), t)
	expectString("ello", string(rope3.Slice(1, 40)), t)
}

func TestDelete(t *testing.T) {
	rope := NewString("how now brown cow")
	rope = rope.Delete(8, 6)

	expectString("how now cow", rope.String(), t)
}

func TestEqual(t *testing.T) {
	rope := NewString("how now brown cow")
	rope = rope.Delete(8, 6)

	if !rope.Equal(NewString("how now cow")) {
		t.Fatalf("expected ropes to be equal")
	}

	rope1 := NewString("hel")
	rope2 := NewString("lo")
	rope3 := Rope{
		depth:  1,
		length: rope1.length + rope2.length,
		left:   &rope1,
		right:  &rope2,
	}

	rope4 := NewString("hel")
	rope5 := NewString("lo")
	rope6 := Rope{
		depth:  1,
		length: rope1.length + rope2.length,
		left:   &rope4,
		right:  &rope5,
	}

	if !rope6.Equal(rope3) {
		t.Fatalf("expected ropes to be equal")
	}

	expectString("hello", rope6.String(), t)
	expectString("hello", rope3.String(), t)

	ropeX := NewString(strings.Repeat("A", 4097)).AppendString(strings.Repeat("A", 1137))
	ropeY := NewString(strings.Repeat("A", 1137)).AppendString(strings.Repeat("A", 4097))
	if !ropeX.Equal(ropeY) {
		t.Fatalf("expected ropes to be equal")
	}

	ropeX = ropeX.AppendString("X")
	ropeY = ropeY.AppendString("Y")
	if ropeX.Equal(ropeY) {
		t.Fatalf("expected ropes not to be equal")
	}
}

func refer(rope Rope) *Rope {
	return &rope
}

func TestBalance(t *testing.T) {
	rope := NewString("hello")
	for i := 0; i < 32; i++ {
		newRope := New()
		rope = Rope{
			length: 5,
			depth:  i + 1,
			left:   refer(rope), // this is such a hack
			right:  &newRope,
		}
	}

	expectString("hello", rope.String(), t)

	if rope.isBalanced() {
		t.Fatalf("expected rope to be unbalanced")
	}

	rope = rope.Rebalance()

	if !rope.isBalanced() {
		t.Fatalf("expected rope to be balanced")
	}

	expectString("hello", rope.String(), t)
	expectInt(5, rope.Length(), t)
}

func TestBigString(t *testing.T) {
	rope := New()
	for i := 0; i < 1048576; i++ {
		rope = rope.AppendString("aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa")
	}

	rope = rope.Rebalance()
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

func TestReadAt(t *testing.T) {
	scottishPlay := `She should have died hereafter;
There would have been a time for such a word.
— To-morrow, and to-morrow, and to-morrow,
Creeps in this petty pace from day to day,
To the last syllable of recorded time;
And all our yesterdays have lighted fools
The way to dusty death. Out, out, brief candle!
Life's but a walking shadow, a poor player
That struts and frets his hour upon the stage
And then is heard no more. It is a tale
Told by an idiot, full of sound and fury
Signifying nothing.`

	rope := NewString("")
	scanner := bufio.NewScanner(strings.NewReader(scottishPlay))
	for scanner.Scan() {
		rope = rope.AppendString(scanner.Text())
	}

	buf := make([]byte, 1000)
	_, err := rope.ReadAt(buf, 120)
	if err != io.EOF {
		t.Fatalf("expected EOF error")
	}

	expectString("Creeps in this petty pace from day to day", string(buf)[:41], t)

	buf = make([]byte, 41)
	_, err = rope.ReadAt(buf, 120)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	expectString("Creeps in this petty pace from day to day", string(buf), t)
}
