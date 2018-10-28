package render

import (
	"time"
	"unicode"

	"github.com/aestek/baleno/keymap"

	"github.com/faiface/pixel/imdraw"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
	"golang.org/x/image/colornames"
)

type Drawer interface {
	DrawCursor(x, y int)
	DrawChar(x, y int, r rune)
}

type Renderable interface {
	Render(d Drawer, w, h int)
	HandleKeyPress(k keymap.KeyPress)
}

type Renderer struct {
	Config         Config
	win            *pixelgl.Window
	txt            *text.Text
	imd            *imdraw.IMDraw
	blockH, blockW float64
}

func NewRenderer(cfg Config) *Renderer {
	renderer := &Renderer{
		Config: cfg,
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
	atlasRunes := []rune{}
	for _, r := range unicode.Scripts {
		atlasRunes = append(atlasRunes, text.RangeTable(r)...)
	}
	atlas := text.NewAtlas(font, atlasRunes)
	txt := text.New(pixel.V(3, 0), atlas)

	txt.Color = colornames.Lightgrey

	renderer.win = win
	renderer.txt = txt
	renderer.imd = imdraw.New(nil)
	renderer.blockH = fontSize
	renderer.blockW = txt.BoundsOf("A").W()

	return renderer
}

func (r *Renderer) Run(drawable Renderable) {
	fps := time.Tick(time.Second / 60)
	for !r.win.Closed() {
		kp := handleKeys(r.win)
		if kp != nil {
			drawable.HandleKeyPress(*kp)
		}
		r.win.Clear(colornames.Ivory)
		size := r.win.Bounds()
		drawable.Render(r, int(size.W()/r.blockW), int(size.H()/r.blockH))
		r.win.Update()
		<-fps
	}
}
