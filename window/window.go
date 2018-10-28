package window

import (
	"github.com/aestek/baleno/keymap"
	"github.com/aestek/baleno/render"
	"github.com/aestek/baleno/state"
	"github.com/aestek/baleno/view"
)

type Window struct {
	state   *state.State
	paneIdx int
}

func New(state *state.State) *Window {
	w := &Window{
		state: state,
	}
	state.Set("window.panes.root.sizing", CompSizing(RatioSideSizing(1), RatioSideSizing(1)))
	state.Set("window.panes.root.split_dir", PaneSplitHorizontal)
	state.Set("window.panes.root.view", nil)
	state.Alias("window.panes.focused", "window.panes.root")

	return w
}

func (w *Window) AddView(v view.View) {
	v.Attach(w.state.Namespace("window.panes.focused.view"))

	paneView, _ := w.state.Get("window.panes.focused.view")

	if paneView == nil {
		w.state.Set("window.panes.focused.view", v)
		return
	}

	paneSplitDir := w.state.MustGet("window.panes.focused.split_dir").(PaneSplitDir)

	w.paneIdx++
	prefix := state.K("window.panes.focused.children", w.paneIdx)
	w.state.Set(state.K(prefix, "view"), paneView)
	w.state.Set(state.K(prefix, "split_dir"), paneSplitDir)
	w.state.Set(state.K(prefix, "sizing"), HalfSizing(paneSplitDir))

	w.paneIdx++
	prefix = state.K("window.panes.focused.children", w.paneIdx)
	w.state.Set(state.K(prefix, "view"), v)
	w.state.Set(state.K(prefix, "split_dir"), paneSplitDir)
	w.state.Set(state.K(prefix, "sizing"), HalfSizing(paneSplitDir))

	w.state.Set("window.panes.focused.view", nil)
	w.state.Alias("window.panes.focused", prefix)
}

func (w *Window) HandleKeyPress(k keymap.KeyPress) {
	v, ok := w.state.Get("window.panes.focused.view")
	if !ok {
		return
	}
	v.(view.View).HandleKeyPress(k)
}

func (w *Window) Render(d render.Drawer, winWidth, winHeight int) {
	w.renderPane(d, "window.panes.root", winWidth, winHeight, 0, 0)
}

func (w *Window) renderPane(d render.Drawer, paneK string, winWidth, winHeight, offsetX, offsetY int) (int, int) {
	paneWidth, paneHeight := w.state.MustGet(state.K(paneK, "sizing")).(Sizing)(winWidth, winHeight)
	paneView := w.state.MustGet(state.K(paneK, "view"))

	if paneView == nil {
		w.state.Range(state.K(paneK, "children"), func(k string, _ interface{}) {
			ox, oy := w.renderPane(d, state.K(paneK, "children", k), paneWidth, paneHeight, offsetX, offsetY)
			offsetX += ox
			offsetY += oy
		})
	} else {
		view := paneView.(view.View)
		view.SetSize(paneWidth, paneHeight)
		b := view.Buffer()
		for y := 0; y < len(b); y++ {
			for x := 0; x < len(b[y]); x++ {
				if b[y][x].Cursor {
					d.DrawCursor(offsetX+x, offsetY+y)
				}
				if b[y][x].Char != 0 {
					d.DrawChar(offsetX+x, offsetY+y, b[y][x].Char)
				}
			}
		}
	}

	splitDir := w.state.MustGet(state.K(paneK, "split_dir")).(PaneSplitDir)

	if splitDir == PaneSplitHorizontal {
		return paneWidth, 0
	} else {
		return 0, paneHeight
	}
}
