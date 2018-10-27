package view

import "github.com/aestek/baleno/keymap"

func InsertKeyMap(view *BufferView) *keymap.Node {
	return &keymap.Node{
		Match: func(keymap.KeyPress) bool {
			return true
		},
		Children: []*keymap.Node{
			{
				Match: func(k keymap.KeyPress) bool {
					return k.Key.Code == keymap.KeyCodeChar
				},
				Action: func(k keymap.KeyPress) {
					view.buffer.Insert(view.CursorOffset(), []rune{k.Key.Value})
					view.CursorAdvance(1, 0)
				},
			},
			{
				Match: func(k keymap.KeyPress) bool {
					return k.Key.Code == keymap.KeyCodeLeft
				},
				Action: func(k keymap.KeyPress) {
					view.CursorAdvance(-1, 0)
				},
			},
			{
				Match: func(k keymap.KeyPress) bool {
					return k.Key.Code == keymap.KeyCodeRight
				},
				Action: func(k keymap.KeyPress) {
					view.CursorAdvance(1, 0)
				},
			},
			{
				Match: func(k keymap.KeyPress) bool {
					return k.Key.Code == keymap.KeyCodeUp
				},
				Action: func(k keymap.KeyPress) {
					view.CursorAdvance(0, -1)
				},
			},
			{
				Match: func(k keymap.KeyPress) bool {
					return k.Key.Code == keymap.KeyCodeDown
				},
				Action: func(k keymap.KeyPress) {
					view.CursorAdvance(0, 1)
				},
			},
			{
				Match: func(k keymap.KeyPress) bool {
					return k.Key.Code == keymap.KeyCodeBackspace
				},
				Action: func(k keymap.KeyPress) {
					view.DeleteBack(1)
				},
			},
			{
				Match: func(k keymap.KeyPress) bool {
					return k.Key.Code == keymap.KeyCodeEnter
				},
				Action: func(k keymap.KeyPress) {
					view.buffer.Insert(view.CursorOffset(), []rune{'\n'})
					view.CursorAdvance(0, 1)
					view.CursorSetX(0)
				},
			},
		},
	}
}
