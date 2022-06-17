package rope

import (
	"io"
)

// A Reader provides an implementation of io.Reader for ropes.
type Reader struct {
	rope     Rope
	position int64
}

// Return a new reader attached to the given rope.
func NewReader(rope Rope) *Reader {
	return rope.Reader()
}

// Read implements the standard Read interface:
// it reads data from the rope, populating p, and returns
// the number of bytes actually read.
func (reader *Reader) Read(p []byte) (n int, err error) {
	n, err = reader.rope.ReadAt(p, reader.position)
	if err == nil {
		reader.position += int64(n)
	}
	return
}

// Return a new Reader attached to this rope.
func (rope Rope) Reader() *Reader {
	return &Reader{rope: rope}
}

// ReadAt implements the standard ReadAt interface:
// it reads len(p) bytes from offset off into p, and returns
// the number of bytes actually read. If n < len(p), err will
// explain the shortfall.
func (rope Rope) ReadAt(p []byte, off int64) (n int, err error) {
	o := int(off)
	for n < len(p) && o+n < rope.Length() {
		leaf, at := rope.leafForOffset(o + n)
		n += copy(p[n:], leaf.content[at:])
	}

	if n < len(p) {
		err = io.EOF
	}

	return
}
