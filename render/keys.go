package render

import (
	"unicode/utf8"

	"github.com/aestek/baleno/keymap"
	"github.com/faiface/pixel/pixelgl"
)

func handleKeys(win *pixelgl.Window) *keymap.KeyPress {
	typed := win.Typed()

	var key keymap.Key

	switch {
	case win.JustPressed(pixelgl.KeyLeft) || win.Repeated(pixelgl.KeyLeft):
		key = keymap.Key{Code: keymap.KeyCodeLeft}
	case win.JustPressed(pixelgl.KeyRight) || win.Repeated(pixelgl.KeyRight):
		key = keymap.Key{Code: keymap.KeyCodeRight}
	case win.JustPressed(pixelgl.KeyUp) || win.Repeated(pixelgl.KeyUp):
		key = keymap.Key{Code: keymap.KeyCodeUp}
	case win.JustPressed(pixelgl.KeyDown) || win.Repeated(pixelgl.KeyDown):
		key = keymap.Key{Code: keymap.KeyCodeDown}
	case win.JustPressed(pixelgl.KeyBackspace) || win.Repeated(pixelgl.KeyBackspace):
		key = keymap.Key{Code: keymap.KeyCodeBackspace}
	case win.JustPressed(pixelgl.KeyEnter) || win.Repeated(pixelgl.KeyEnter):
		key = keymap.Key{Code: keymap.KeyCodeEnter}
	case len(typed) > 0:
		r, _ := utf8.DecodeRuneInString(typed)
		key = keymap.Key{Code: keymap.KeyCodeChar, Value: r}
	default:
		return nil
	}
	return &keymap.KeyPress{
		Key: key,
	}
}
