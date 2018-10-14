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
	indexes     map[string]*Index
}

func NewBytes(b io.Reader, enc *UTF8Encoding) (*Bytes, error) {
	raw, err := ioutil.ReadAll(b)
	if err != nil {
		return nil, err
	}

	runes := enc.Decode(raw)

	return &Bytes{
		contents: runes,
		indexes:  map[string]*Index{},
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

func (b *Bytes) Index(search string) []int {
	b.indexesLock.Lock()
	defer b.indexesLock.Unlock()
	_, ok := b.indexes[search]
	if !ok {
		b.indexes[search] = NewIndex([]rune(search), b.contents)
	}
	return b.indexes[search].Entries()
}

func (b *Bytes) LineCount() int {
	return len(b.Index("\n")) + 1
}

func (b *Bytes) LineLength(n int) int {
	idx := b.Index("\n")
	lc := b.LineCount()
	if lc == 1 && n != 0 {
		return 0
	}
	if lc == 1 {
		return len(b.contents)
	}
	if n == 0 {
		return idx[0]
	}
	if n == len(idx) {
		return len(b.contents) - 1 - idx[len(idx)-1]
	}
	if n > len(idx) {
		return 0
	}
	return idx[n] - 1 - idx[n-1]
}

func (b *Bytes) buildIndexes() {
	for _, i := range b.indexes {
		go i.Build(b.contents, 0, len(b.contents))
	}
}
