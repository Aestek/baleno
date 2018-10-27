package view

import (
	"github.com/aestek/baleno/keymap"
)

type View interface {
	Buffer() DrawBuffer
	SetSize(x, y int)
	HandleKeyPress(k keymap.KeyPress)
}
