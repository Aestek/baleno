package view

type DrawBuffer [][]DrawBlock

type DrawBlock struct {
	Cursor   bool
	Selected bool
	Char     rune
}

func NewDrawBuffer(width, height int) DrawBuffer {
	b := make(DrawBuffer, height)
	for i := 0; i < height; i++ {
		b[i] = make([]DrawBlock, width)
	}
	return b
}

func (d DrawBuffer) SetCursor(x, y int) {
	d[y][x].Cursor = true
}

func (d DrawBuffer) SetChar(x, y int, c rune) {
	d[y][x].Char = c
}

func (d DrawBuffer) Clear() {
	for y := 0; y < len(d); y++ {
		for x := 0; x < len(d[y]); x++ {
			d[y][x] = DrawBlock{}
		}
	}
}
