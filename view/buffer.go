package view

import (
	"fmt"
	"log"

	"github.com/aestek/baleno/buffer"
	"github.com/aestek/baleno/keymap"
)

const EOL = '\n'

const eolIdx = "\n"

type Drawer interface {
	Do(func())
	DrawCursor(x, y int)
	DrawChar(x, y int, r rune)
	KeyPresses() <-chan keymap.KeyPress
}

type BufferView struct {
	drawer        Drawer
	buffer        *buffer.Bytes
	width, height int
	line, col     int
	cursor        Cursor
	keyMap        *keymap.Node
}

func NewBuffer(b *buffer.Bytes, d Drawer, w, h int) *BufferView {
	v := &BufferView{
		buffer: b,
		drawer: d,
		width:  w,
		height: h,
	}

	v.keyMap = InsertKeyMap(v)

	go func() {
		for k := range d.KeyPresses() {
			v.keyMap.Exec(k)
			v.Draw()
		}
	}()

	return v
}

func (v *BufferView) Draw() {
	v.drawer.Do(func() {
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
	})
}

func (v *BufferView) drawLine(line []rune, ln int) {
	if v.cursor.Y == ln {
		x := v.cursor.X
		if x > len(line) {
			x = len(line)
		}

		dx := x - v.col
		if dx >= 0 && dx < v.width {
			v.drawer.DrawCursor(x, ln-v.line)
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
		v.drawer.DrawChar(i, ln-v.line, c)
	}
}

func (v *BufferView) Buffer() buffer.Buffer {
	return v.buffer
}

func (v *BufferView) Cursor() Cursor {
	return v.cursor
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

		v.cursor.X = cx + x
	}

	if cy+y < 0 {
		v.cursor.Y = 0
	} else if cy+y > lc-1 {
		v.cursor.Y = lc - 1
	} else {
		v.cursor.Y = cy + y
	}

	fmt.Println(v.cursor)
}

func (v *BufferView) CursorSetX(x int) {
	v.cursor.X = x
}

func (v *BufferView) CursorSetY(y int) {
	v.cursor.Y = y
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
	v.cursor.X, v.cursor.Y = v.OffsetToCursor(cp - n)
	fmt.Println(v.cursor)
}

func (v *BufferView) cursorReal() (int, int) {
	x := v.cursor.X

	ln := v.buffer.LineLength(v.cursor.Y)
	if v.cursor.X > ln {
		x = ln
	}

	return x, v.cursor.Y
}
