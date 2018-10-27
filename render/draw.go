package render

import (
	"golang.org/x/image/colornames"

	"github.com/faiface/pixel/imdraw"

	"github.com/faiface/pixel"
)

const (
	cursorMarginRight = 3
	cursorMarginTop   = 5
)

func (r *Renderer) DrawCursor(x, y int) {
	r.imd.Clear()

	r.imd.Color = colornames.Rosybrown
	r.imd.EndShape = imdraw.SharpEndShape

	margin := pixel.V(cursorMarginRight, -cursorMarginTop)
	r.imd.Push(r.at(x, y).Add(margin))
	r.imd.Push(r.at(x+1, y).Add(margin))
	r.imd.Push(r.at(x+1, y+1).Add(margin))
	r.imd.Push(r.at(x, y+1).Add(margin))
	r.imd.Polygon(0)

	r.imd.Draw(r.win)
}

func (r *Renderer) DrawChar(x, y int, c rune) {
	r.txt.Color = colornames.Black
	r.txt.Clear()
	r.txt.WriteRune(c)
	r.txt.Draw(r.win, pixel.IM.Moved(r.at(x, y+1)))
}

func (r *Renderer) at(x, y int) pixel.Vec {
	return pixel.V(
		float64(x)*r.blockW,
		r.win.Bounds().H()-float64(y)*r.blockH,
	)
}
