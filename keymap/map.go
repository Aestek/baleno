package keymap

type KeyMatch func(KeyPress) bool

type Node struct {
	Match    KeyMatch
	Action   func(k KeyPress)
	Children []*Node
}

func (n *Node) Exec(k KeyPress) {
	if !n.Match(k) {
		return
	}

	if n.Action != nil {
		n.Action(k)
	}

	for _, c := range n.Children {
		c.Exec(k)
	}
}
