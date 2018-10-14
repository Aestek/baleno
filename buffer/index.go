package buffer

import "sync"

type Index struct {
	seq       []rune
	positions []int
	lock      sync.RWMutex
}

func NewIndex(seq []rune, c []rune) *Index {
	i := &Index{
		seq: seq,
	}
	i.Build(c, 0, len(c))
	return i
}

func (idx *Index) Build(c []rune, from, to int) {
	idx.lock.Lock()
	defer idx.lock.Unlock()

	idx.positions = []int{}
Next:
	for i := from; i < to; i++ {
		for j, r := range idx.seq {
			if i+j >= to {
				return
			}

			if c[i+j] != r {
				continue Next
			}
		}
		idx.positions = append(idx.positions, i)
	}
}

func (idx *Index) Entries() []int {
	idx.lock.RLock()
	defer idx.lock.RUnlock()
	return idx.positions
}
