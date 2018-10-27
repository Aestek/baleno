package view

import (
	"github.com/aestek/baleno/keymap"
	"github.com/aestek/baleno/render"
)

type PaneSplitDir int

const (
	PaneSplitVertical PaneSplitDir = iota
	PaneSplitHorizontal
)

type Pane struct {
	SplitDir PaneSplitDir
	Children []*Pane

	View View
	X, Y int
}

func NewViewPane(view View) *Pane {
	return &Pane{
		View: view,
	}
}

type Window struct {
	Root    *Pane
	Focused *Pane
}

func NewWindow(root *Pane) *Window {
	return &Window{
		Root:    root,
		Focused: root,
	}
}

func (w *Window) HandleKeyPress(k keymap.KeyPress) {
	w.Focused.View.HandleKeyPress(k)
}

func (w *Window) Render(d render.Drawer) {
	w.renderPane(d, w.Root)
}

func (w *Window) renderPane(d render.Drawer, pane *Pane) {
	if pane.View == nil {
		for _, p := range pane.Children {
			w.renderPane(d, p)
		}
	} else {
		b := pane.View.Buffer()
		for y := 0; y < len(b); y++ {
			for x := 0; x < len(b[y]); x++ {
				if b[y][x].Cursor {
					d.DrawCursor(pane.X+x, pane.Y+y)
				}
				if b[y][x].Char != 0 {
					d.DrawChar(pane.X+x, pane.Y+y, b[y][x].Char)
				}
			}
		}
	}
}
