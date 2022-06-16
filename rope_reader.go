package rope

import (
	"io"
)

type Reader struct {
	rope     *Rope
	position int64
}

func (reader *Reader) Read(p []byte) (n int, err error) {
	n, err = reader.rope.ReadAt(p, reader.position)
	if err == nil {
		reader.position += int64(n)
	}
	return
}

func (rope *Rope) Reader() *Reader {
	return &Reader{rope: rope}
}

func (rope *Rope) ReadAt(p []byte, off int64) (n int, err error) {
	for n < len(p) && n < rope.Length() {
		leaf, at := rope.leafForOffset(int(off))
		n += copy(p[n:], leaf.content[at:])
	}

	if n < len(p) {
		err = io.EOF
	}

	return
}
