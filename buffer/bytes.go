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
	indexes     map[IndexDef]*Index
}

func NewBytes(b io.Reader, enc *UTF8Encoding) (*Bytes, error) {
	raw, err := ioutil.ReadAll(b)
	if err != nil {
		return nil, err
	}

	runes := enc.Decode(raw)

	return &Bytes{
		contents: runes,
		indexes:  map[IndexDef]*Index{},
	}, nil
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

	b.buildIndexes()
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
	b.buildIndexes()
	return nil
}

func (b *Bytes) ReadOnly() bool {
	return false
}

func (b *Bytes) Length() int {
	return len(b.contents)
}

func (b *Bytes) Index(def IndexDef) []int {
	b.indexesLock.Lock()
	defer b.indexesLock.Unlock()
	_, ok := b.indexes[def]
	if !ok {
		b.indexes[def] = NewIndex(def, b.contents)
	}
	return b.indexes[def].Entries()
}

func (b *Bytes) LineCount() int {
	return len(b.Index(StartOfLineIdx))
}

func (b *Bytes) LineLength(n int) int {
	idx := b.Index(StartOfLineIdx)

	if n == 0 && len(idx) == 1 {
		return len(b.contents)
	}

	if n == len(idx)-1 {
		return len(b.contents) - idx[len(idx)-1]
	}

	return idx[n+1] - idx[n] - 1
}

func (b *Bytes) buildIndexes() {
	for _, i := range b.indexes {
		i.Build(b.contents, 0, len(b.contents))
	}
}
