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
