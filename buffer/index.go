package buffer

import "sync"

var StartOfLineIdx = IndexDef{
	Search:     "\n",
	After:      true,
	ForceFirst: true,
}

var EndOfLineIdx = IndexDef{
	Search: "\n",
}

type IndexDef struct {
	Search     string
	After      bool
	ForceFirst bool
}

type Index struct {
	positions []int
	lock      sync.RWMutex
	def       IndexDef
}

func NewIndex(def IndexDef, c []rune) *Index {
	i := &Index{
		def: def,
	}
	i.Build(c, 0, len(c))
	return i
}

func (idx *Index) Build(c []rune, from, to int) {
	idx.lock.Lock()
	defer idx.lock.Unlock()

	seq := []rune(idx.def.Search)
	idx.positions = []int{}
	if idx.def.ForceFirst {
		idx.positions = append(idx.positions, 0)
	}

Next:
	for i := from; i < to; i++ {
		for j, r := range seq {
			if i+j >= to {
				return
			}

			if c[i+j] != r {
				continue Next
			}
		}
		pos := i
		if idx.def.After {
			pos++
		}
		idx.positions = append(idx.positions, pos)
	}
}

func (idx *Index) Entries() []int {
	idx.lock.RLock()
	defer idx.lock.RUnlock()
	return idx.positions
}

type indexes struct {
	indexes  map[IndexDef]*Index
	lock     sync.Mutex
	contents func() []rune
}

func (i *indexes) Index(def IndexDef) []int {
	i.lock.Lock()
	defer i.lock.Unlock()
	_, ok := i.indexes[def]
	if !ok {
		i.indexes[def] = NewIndex(def, i.contents())
	}
	return i.indexes[def].Entries()
}

func (i *indexes) LineCount() int {
	return len(i.Index(StartOfLineIdx))
}

func (i *indexes) LineLength(n int) int {
	content := i.contents()
	idx := i.Index(StartOfLineIdx)

	if n == 0 && len(idx) == 1 {
		return len(content)
	}

	if n == len(idx)-1 {
		return len(content) - idx[len(idx)-1]
	}

	return idx[n+1] - idx[n] - 1
}

func (i *indexes) build() {
	content := i.contents()
	for _, i := range i.indexes {
		i.Build(content, 0, len(content))
	}
}
