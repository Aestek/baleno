package view

import (
	"github.com/aestek/baleno/keymap"
	"github.com/aestek/baleno/state"
)

type View interface {
	Buffer() DrawBuffer
	SetSize(x, y int)
	Attach(s *state.State)
	HandleKeyPress(k keymap.KeyPress)
}
