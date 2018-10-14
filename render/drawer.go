package render

import (
	"golang.org/x/image/colornames"

	"github.com/aestek/baleno/keymap"
	"github.com/faiface/pixel/imdraw"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
)

type drawerDoer struct {
	d *Drawer
}

type Drawer struct {
	win            *pixelgl.Window
	txt            *text.Text
	imd            *imdraw.IMDraw
	keyPresses     chan keymap.KeyPress
	blockW, blockH float64
	addToQueue     func(func())
}

func (d *Drawer) Do(fn func()) {
	d.addToQueue(fn)
}

func (d *Drawer) DrawCursor(x, y int) {
	d.imd.Clear()

	d.imd.Color = colornames.Rosybrown
	d.imd.EndShape = imdraw.SharpEndShape
	d.imd.Push(d.at(x, y))
	d.imd.Push(d.at(x+1, y))
	d.imd.Push(d.at(x+1, y+1).Add(pixel.V(0, -3)))
	d.imd.Push(d.at(x, y+1).Add(pixel.V(0, -3)))
	d.imd.Polygon(0)

	d.imd.Draw(d.win)
}

func (d *Drawer) DrawChar(x, y int, r rune) {
	d.txt.Color = colornames.Black
	d.txt.Clear()
	d.txt.WriteRune(r)
	d.txt.Draw(d.win, pixel.IM.Moved(d.at(x, y+1)))
}

func (d *Drawer) KeyPresses() <-chan keymap.KeyPress {
	return d.keyPresses
}

func (d *Drawer) at(x, y int) pixel.Vec {
	return pixel.V(
		float64(x)*d.blockW,
		d.win.Bounds().H()-float64(y)*d.blockH,
	)
}
