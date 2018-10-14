package view

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aestek/baleno/buffer"
	"github.com/stretchr/testify/assert"
)

func TestCursor(t *testing.T) {
	type tcase struct {
		in       string
		initial  Cursor
		advanceX int
		advanceY int
		expected Cursor
	}

	cases := []tcase{
		{
			in:       "",
			advanceX: 1,
			expected: Cursor{0, 0},
		},
		{
			in:       "",
			advanceY: 1,
			expected: Cursor{0, 0},
		},
		{
			in:       "a",
			advanceX: 1,
			expected: Cursor{1, 0},
		},
		{
			in:       "a",
			advanceX: 2,
			expected: Cursor{1, 0},
		},
		{
			in:       "a",
			advanceX: -1,
			expected: Cursor{0, 0},
		},
		{
			in:       "a\n",
			advanceX: 1,
			expected: Cursor{1, 0},
		},
		{
			in:       "a\n",
			advanceX: 1,
			initial:  Cursor{1, 0},
			expected: Cursor{0, 1},
		},
		{
			in:       "a\n",
			advanceX: 2,
			expected: Cursor{0, 1},
		},
		{
			in:       "a\n",
			advanceY: 1,
			expected: Cursor{0, 1},
		},
		{
			in:       "a\na",
			advanceX: 1,
			initial:  Cursor{0, 1},
			expected: Cursor{1, 1},
		},
		{
			in:       "Hello\naaa",
			advanceX: 1,
			initial:  Cursor{5, 0},
			expected: Cursor{0, 1},
		},
		{
			in:       "Hello\naaa",
			advanceX: 1,
			initial:  Cursor{2, 1},
			expected: Cursor{3, 1},
		},
	}

	for _, c := range cases {
		b, err := buffer.NewBytes(strings.NewReader(c.in), &buffer.UTF8Encoding{})
		assert.Nil(t, err)
		v := NewBuffer(b, &FakeDrawer{}, 1000, 1000)
		v.cursor = c.initial
		v.CursorAdvance(c.advanceX, c.advanceY)
		assert.Equal(t, c.expected, v.cursor, fmt.Sprintf("%+v", c))
		co := v.CursorOffset()
		x, y := v.OffsetToCursor(co)
		assert.Equal(t, c.expected, Cursor{x, y}, "cursor offset: %+v", c)
	}
}
