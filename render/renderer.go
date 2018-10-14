package render

import (
	"time"

	"github.com/aestek/baleno/keymap"

	"github.com/faiface/pixel/imdraw"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
	"golang.org/x/image/colornames"
)

type Renderer struct {
	Config         Config
	win            *pixelgl.Window
	txt            *text.Text
	imd            *imdraw.IMDraw
	keyPresses     chan keymap.KeyPress
	drawQueue      []func()
	blockH, blockW float64
}

func NewRenderer(cfg Config) *Renderer {
	renderer := &Renderer{
		Config:     cfg,
		keyPresses: make(chan keymap.KeyPress),
	}

	winCfg := pixelgl.WindowConfig{
		Title:  "Baleno Editor",
		Bounds: pixel.R(0, 0, 1024, 768),
		VSync:  true,
	}
	win, err := pixelgl.NewWindow(winCfg)
	if err != nil {
		panic(err)
	}

	fontSize := 30.0

	font := loadFont(renderer.Config.FontPath, fontSize)
	atlas := text.NewAtlas(font, text.ASCII)
	txt := text.New(pixel.V(3, 0), atlas)

	txt.Color = colornames.Lightgrey

	renderer.win = win
	renderer.txt = txt
	renderer.imd = imdraw.New(nil)
	renderer.blockH = fontSize
	renderer.blockW = txt.BoundsOf("A").W()

	return renderer
}

func (r *Renderer) Run() {
	fps := time.Tick(time.Second / 120)
	for !r.win.Closed() {
		handleKeys(r.win, r.keyPresses)

		if len(r.drawQueue) > 0 {
			r.win.Clear(colornames.Ivory)
			for _, f := range r.drawQueue {
				f()
			}
			r.drawQueue = r.drawQueue[:0]
			r.win.Update()
		} else {
			r.win.UpdateInput()
		}

		<-fps
	}
}
