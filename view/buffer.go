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
	contents, err := v.buffer.Read(0, -1)
	if err != nil {
		log.Println(err)
		return
	}

	v.drawer.Do(func() {
		lines(contents, func(ln int, line []rune) bool {
			if ln < v.line {
				return true
			}
			if ln > v.line+v.height {
				return false
			}

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
				return true
			}
			until := v.col + v.width
			if until > len(line) {
				until = len(line)
			}

			for i, c := range line[v.col:until] {
				v.drawer.DrawChar(i, ln-v.line, c)
			}

			return true
		})
	})
}

func (b *BufferView) Buffer() buffer.Buffer {
	return b.buffer
}

func (b *BufferView) Cursor() Cursor {
	return b.cursor
}

func (b *BufferView) CursorAdvance(x, y int) {
	cx, cy := b.cursorReal()
	lc := b.buffer.LineCount()

	if x != 0 {
		for cx+x < 0 {
			if cy == 0 {
				cx = 0
				x = 0
				break
			}
			x += cx + 1
			cy--
			cx = b.buffer.LineLength(cy)
		}

		for cx+x > b.buffer.LineLength(cy) {
			if cy == lc-1 {
				cx = b.buffer.LineLength(cy)
				x = 0
				break
			}
			x -= b.buffer.LineLength(cy) - cx + 1
			cy++
			cx = 0
		}

		b.cursor.X = cx + x
	}

	if cy+y < 0 {
		b.cursor.Y = 0
	} else if cy+y > lc-1 {
		b.cursor.Y = lc - 1
	} else {
		b.cursor.Y = cy + y
	}

	fmt.Println(b.cursor)
}

func (b *BufferView) CursorSetX(x int) {
	b.cursor.X = x
}

func (b *BufferView) CursorSetY(y int) {
	b.cursor.Y = y
}

func (b *BufferView) CursorOffset() int {
	x, y := b.cursorReal()
	if y == 0 {
		return x
	}
	idx := b.buffer.Index(eolIdx)
	return idx[y-1] + 1 + x
}

func (b *BufferView) OffsetToCursor(n int) (int, int) {
	y, x := 0, 0
	idx := b.buffer.Index(eolIdx)
	lastPos := 0
	for i, pos := range idx {
		if n > pos {
			y = i + 1
			lastPos = pos + 1
		} else {
			break
		}
	}
	x = n - lastPos
	return x, y
}

func (b *BufferView) DeleteBack(n int) {
	cp := b.CursorOffset()
	if n < 0 {
		return
	}
	if n > cp {
		n = cp
	}
	b.buffer.Delete(cp-n, n)
	b.cursor.X, b.cursor.Y = b.OffsetToCursor(cp - n)
	fmt.Println(b.cursor)
}

func (b *BufferView) cursorReal() (int, int) {
	x := b.cursor.X

	ln := b.buffer.LineLength(b.cursor.Y)
	if b.cursor.X > ln {
		x = ln
	}

	return x, b.cursor.Y
}
