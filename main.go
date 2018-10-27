package main

import (
	"bytes"
	"log"

	"github.com/aestek/baleno/buffer"
	"github.com/aestek/baleno/render"
	"github.com/aestek/baleno/view"
	"github.com/faiface/pixel/pixelgl"
)

func main() {
	pixelgl.Run(func() {
		r := render.NewRenderer(
			render.Config{
				FontPath: "/home/aestek/.dotfiles/.fonts/Inconsolata-Regular.ttf",
			},
		)

		b, err := buffer.NewBytes(
			bytes.NewBuffer([]byte("Hello\nlol")),
			&buffer.UTF8Encoding{},
		)
		if err != nil {
			log.Fatal(err)
		}

		v := view.NewBuffer(b, 1000, 1000)
		v.Draw()

		pane := view.NewViewPane(v)
		pane.X = 10
		pane.Y = 5
		r.Run(view.NewWindow(pane))
	})
}
