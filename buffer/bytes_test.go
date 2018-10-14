package buffer

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLineCount(t *testing.T) {
	type tcase struct {
		in  string
		exp int
	}

	cases := []tcase{
		{"", 1},
		{"\n", 2},
		{"aaa", 1},
		{"aaa\n", 2},
		{"aaa\nbbbadf", 2},
		{"aaa\nbbbadf\n", 3},
		{"aaa\nbbbadf\nzdazda", 3},
		{"aaa\nbbbadf\nzdazda\n", 4},
	}

	for _, c := range cases {
		b, err := NewBytes(strings.NewReader(c.in), &UTF8Encoding{})
		assert.Nil(t, err)
		assert.Equal(t, c.exp, b.LineCount())
	}
}

func TestLineLength(t *testing.T) {
	type tcase struct {
		in  string
		exp []int
	}

	cases := []tcase{
		{"", []int{0}},
		{"\n", []int{0, 0}},
		{"aaa", []int{3}},
		{"aaa\n", []int{3, 0}},
		{"aaa\nbbbadf", []int{3, 6}},
		{"aaa\nbbbadf\n", []int{3, 6, 0}},
		{"aaa\nbbbadf\nzdazda", []int{3, 6, 6}},
		{"aaa\nbbbadf\nzdazda\n", []int{3, 6, 6, 0}},
	}

	for _, c := range cases {
		b, err := NewBytes(strings.NewReader(c.in), &UTF8Encoding{})
		assert.Nil(t, err)

		for i := 0; i < b.LineCount(); i++ {
			assert.Equal(t, c.exp[i], b.LineLength(i), "%+v, line %d", c, i)
		}
	}
}
