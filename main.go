package main

import (
	"bytes"
	"log"

	"github.com/aestek/baleno/state"
	"github.com/aestek/baleno/window"

	"github.com/aestek/baleno/buffer"
	"github.com/aestek/baleno/render"
	"github.com/aestek/baleno/view"
	"github.com/faiface/pixel/pixelgl"
)

func main() {
	pixelgl.Run(func() {
		state := state.New()

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

		win := window.New(state)
		win.AddView(view.NewStatusView(state))
		win.AddView(view.NewBuffer(b))

		r.Run(win)
	})
}
