package state

import (
	"fmt"
	"strings"
	"sync"
)

type Node struct {
	Key      string
	Alias    *Node
	Aliased  []*Node
	Value    interface{}
	Watches  []chan Event
	Children []*Node
}

type Event struct {
	Key   string
	Value interface{}
}

type State struct {
	root *Node
	lock *sync.RWMutex
}

func New() *State {
	return &State{
		root: &Node{},
		lock: &sync.RWMutex{},
	}
}

func (s *State) set(k string, v interface{}) *Node {
	parts := strings.Split(k, ".")
	path := []*Node{s.root}
	node := s.root
	partIdx := 0
Run:
	for _, n := range node.Children {
		if n.Key != parts[partIdx] {
			continue
		}

		if n.Alias != nil {
			n = n.Alias
		}

		path = append(path, n)
		path = append(path, n.Aliased...)

		if partIdx == len(parts)-1 {
			n.Value = v
			evt := Event{Key: k, Value: v}
			for _, n := range path {
				for _, w := range n.Watches {
					select {
					case w <- evt:
					default:
					}
				}
			}
			return n
		} else {
			node = n
			partIdx++
			goto Run
		}
	}
	node.Children = append(node.Children, &Node{
		Key: parts[partIdx],
	})
	goto Run
}

func (s *State) Set(k string, v interface{}) {
	s.lock.Lock()
	defer s.lock.Unlock()

	s.set(k, v)
}

func (s *State) get(k string, followAliases bool) *Node {
	parts := strings.Split(k, ".")
	node := s.root
	partIdx := 0
Run:
	for _, n := range node.Children {
		if n.Key != parts[partIdx] {
			continue
		}

		if n.Alias != nil && followAliases {
			n = n.Alias
		}

		if partIdx == len(parts)-1 {
			return n
		} else {
			node = n
			partIdx++
			goto Run
		}
	}
	return nil
}

func (s *State) Get(k string) (interface{}, bool) {
	s.lock.RLock()
	defer s.lock.RUnlock()

	node := s.get(k, true)
	if node == nil {
		return nil, false
	}

	return node.Value, node.Value != nil
}

func (s *State) MustGet(k string) interface{} {
	v, _ := s.Get(k)
	return v
}

func (s *State) Range(k string, cb func(k string, value interface{})) {
	node := s.get(k, true)
	if node == nil {
		return
	}

	for _, c := range node.Children {
		cb(c.Key, c.Value)
	}
}

func (s *State) Delete(k string) {
	s.Set(k, nil)
}

func (s *State) Alias(from string, to string) error {
	s.lock.Lock()
	defer s.lock.Unlock()

	toNode := s.get(to, true)
	if toNode == nil {
		return fmt.Errorf("state: alias: target key `%s` does not exist", to)
	}

	fromNode := s.get(from, false)
	if fromNode == nil {
		fromNode = s.set(from, nil)
	} else if fromNode.Alias != nil {
		aliases := fromNode.Alias.Aliased
		for i := range aliases {
			if aliases[i] == fromNode {
				aliases[i], aliases[len(aliases)-1] = aliases[len(aliases)-1], aliases[i]
				aliases = aliases[:len(aliases)-1]
				break
			}
		}
	}

	fromNode.Alias = toNode
	toNode.Aliased = append(toNode.Aliased, fromNode)

	return nil
}

func (s *State) Watch(k string, c chan Event) error {
	s.lock.Lock()
	defer s.lock.Unlock()

	node := s.get(k, false)
	if node == nil {
		return fmt.Errorf("state: watch: key `%s` does not exist", k)
	}

	node.Watches = append(node.Watches, c)
	return nil
}

func (s *State) Namespace(k string) *State {
	s.lock.Lock()
	defer s.lock.Unlock()

	node := s.get(k, true)
	if node == nil {
		node = s.set(k, nil)
	}

	return &State{root: node, lock: s.lock}
}
