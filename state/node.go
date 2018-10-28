package state

import (
	"strings"
)

type Node struct {
	Key      string
	Parent   *Node
	Alias    *Node
	Aliased  []*Node
	Value    interface{}
	Watches  []chan Event
	Children []*Node
}

func (n *Node) Name() string {
	p := []string{}
	for ; n != nil; n = n.Parent {
		p = append(p, n.Key)
	}

	for i, j := 0, len(p)-1; i < j; i, j = i+1, j-1 {
		p[i], p[j] = p[j], p[i]
	}

	return strings.Join(p, ".")
}
