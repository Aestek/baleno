package view

import "github.com/aestek/baleno/keymap"

type FakeDrawer struct {
}

func (f *FakeDrawer) Do(func())                 {}
func (f *FakeDrawer) DrawCursor(x, y int)       {}
func (f *FakeDrawer) DrawChar(x, y int, r rune) {}
func (f *FakeDrawer) KeyPresses() <-chan keymap.KeyPress {
	return nil
}
