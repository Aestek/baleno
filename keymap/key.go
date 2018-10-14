package keymap

type KeyCode int

const (
	KeyCodeEnter KeyCode = iota
	KeyCodeBackspace
	KeyCodeChar

	KeyCodeLeft
	KeyCodeRight
	KeyCodeUp
	KeyCodeDown
)

type Key struct {
	Code  KeyCode
	Value rune
}

type KeyPress struct {
	Key   Key
	Ctrl  bool
	Shift bool
	Alt   bool
}
