package buffer

import (
	"fmt"
	"io"
	"io/ioutil"
	"sync"
)

var _ Buffer = &Bytes{}

var ErrOutOfRange = fmt.Errorf("out of range")

type Bytes struct {
	contents    []rune
	indexesLock sync.Mutex
	*indexes
}

func NewBytes(b io.Reader, enc *UTF8Encoding) (*Bytes, error) {
	raw, err := ioutil.ReadAll(b)
	if err != nil {
		return nil, err
	}

	runes := enc.Decode(raw)

	bytes := &Bytes{
		contents: runes,
	}

	idx := &indexes{
		contents: func() []rune {
			return bytes.contents
		},
		indexes: make(map[IndexDef]*Index),
	}

	bytes.indexes = idx
	return bytes, nil
}

func (b *Bytes) Read(from, to int) ([]rune, error) {
	if to == -1 {
		return b.contents[from:], nil
	}
	return b.contents[from:to], nil
}

func (b *Bytes) Insert(at int, contents []rune) error {
	b.indexesLock.Lock()
	defer b.indexesLock.Unlock()

	b.contents = append(
		b.contents[:at],
		append(contents, b.contents[at:]...)...,
	)

	b.indexes.build()
	return nil
}

func (b *Bytes) Delete(at int, n int) error {
	if n < 0 {
		return ErrOutOfRange
	}
	if at < 0 {
		return ErrOutOfRange
	}
	if at+n > len(b.contents) {
		return ErrOutOfRange
	}
	b.indexesLock.Lock()
	defer b.indexesLock.Unlock()
	b.contents = append(
		b.contents[:at],
		b.contents[at+n:]...,
	)
	b.indexes.build()
	return nil
}

func (b *Bytes) ReadOnly() bool {
	return false
}

func (b *Bytes) Length() int {
	return len(b.contents)
}
