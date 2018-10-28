package view

import (
	"log"
	"sync"

	"github.com/aestek/baleno/buffer"
	"github.com/aestek/baleno/keymap"
	"github.com/aestek/baleno/state"
)

type BufferView struct {
	drawBuffer     DrawBuffer
	drawBufferLock sync.Mutex
	buffer         buffer.Buffer
	width, height  int
	line, col      int
	keyMap         *keymap.Node
	state          *state.State
}

func NewBuffer(b *buffer.Bytes) *BufferView {
	v := &BufferView{
		buffer: b,
	}

	v.keyMap = InsertKeyMap(v)

	return v
}

func (v *BufferView) HandleKeyPress(k keymap.KeyPress) {
	v.keyMap.Exec(k)
	v.draw()
}

func (v *BufferView) draw() {
	v.drawBufferLock.Lock()
	defer v.drawBufferLock.Unlock()

	v.drawBuffer.Clear()

	idx := v.buffer.Index(buffer.StartOfLineIdx)
	contents, err := v.buffer.Read(0, -1)
	if err != nil {
		log.Println(err)
		return
	}
	if len(idx) == 0 {
		return
	}

	for i := 0; i < len(idx)-1; i++ {
		line := contents[idx[i]:idx[i+1]]
		v.drawLine(line, i)
	}
	v.drawLine(contents[idx[len(idx)-1]:], len(idx)-1)
}

func (v *BufferView) drawLine(line []rune, ln int) {
	if v.CursorY() == ln {
		x := v.CursorX()
		if x > len(line) {
			x = len(line)
		}

		dx := x - v.col
		if dx >= 0 && dx < v.width {
			v.drawBuffer.SetCursor(x, ln-v.line)
		}
	}

	if v.col >= len(line) {
		return
	}
	until := v.col + v.width
	if until > len(line) {
		until = len(line)
	}

	for i, c := range line[v.col:until] {
		v.drawBuffer.SetChar(i, ln-v.line, c)
	}
}

func (v *BufferView) Buffer() DrawBuffer {
	v.drawBufferLock.Lock()
	defer v.drawBufferLock.Unlock()
	return v.drawBuffer
}

func (v *BufferView) SetSize(w, h int) {
	if w == v.width && h == v.height {
		return
	}
	v.drawBufferLock.Lock()
	defer v.draw()
	defer v.drawBufferLock.Unlock()
	v.width = w
	v.height = h
	v.drawBuffer = NewDrawBuffer(w, h)
}

func (v *BufferView) CursorAdvance(x, y int) {
	cx, cy := v.cursorReal()
	lc := v.buffer.LineCount()

	if x != 0 {
		for cx+x < 0 {
			if cy == 0 {
				cx = 0
				x = 0
				break
			}
			x += cx + 1
			cy--
			cx = v.buffer.LineLength(cy)
		}

		for cx+x > v.buffer.LineLength(cy) {
			if cy == lc-1 {
				cx = v.buffer.LineLength(cy)
				x = 0
				break
			}
			x -= v.buffer.LineLength(cy) - cx + 1
			cy++
			cx = 0
		}

		v.CursorSetX(cx + x)
	}

	if cy+y < 0 {
		v.CursorSetY(0)
	} else if cy+y > lc-1 {
		v.CursorSetY(lc - 1)
	} else {
		v.CursorSetY(cy + y)
	}
}

func (v *BufferView) CursorX() int {
	return v.state.MustGet("cursor.x").(int)
}

func (v *BufferView) CursorY() int {
	return v.state.MustGet("cursor.y").(int)
}

func (v *BufferView) CursorSetX(x int) {
	v.state.Set("cursor.x", x)
}

func (v *BufferView) CursorSetY(y int) {
	v.state.Set("cursor.y", y)
}

func (v *BufferView) CursorOffset() int {
	idx := v.buffer.Index(buffer.StartOfLineIdx)

	x, y := v.cursorReal()
	return idx[y] + x
}

func (v *BufferView) OffsetToCursor(n int) (int, int) {
	idx := v.buffer.Index(buffer.StartOfLineIdx)

	y, x := 0, 0
	for i, pos := range idx {
		if pos > n {
			break
		}

		y = i
		x = n - pos
	}
	return x, y
}

func (v *BufferView) DeleteBack(n int) {
	cp := v.CursorOffset()
	if n < 0 {
		return
	}
	if n > cp {
		n = cp
	}
	v.buffer.Delete(cp-n, n)
	ox, oy := v.OffsetToCursor(cp - n)
	v.CursorSetX(ox)
	v.CursorSetY(oy)
}

func (v *BufferView) cursorReal() (int, int) {
	x := v.CursorX()
	y := v.CursorY()

	ln := v.buffer.LineLength(y)
	if x > ln {
		x = ln
	}

	return x, y
}

func (v *BufferView) Attach(s *state.State) {
	v.state = s
	s.Set("cursor.x", 0)
	s.Set("cursor.y", 0)
}
