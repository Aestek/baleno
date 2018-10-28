package view

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aestek/baleno/buffer"
	"github.com/aestek/baleno/state"
	"github.com/stretchr/testify/assert"
)

func TestCursor(t *testing.T) {
	type tcase struct {
		in        string
		initialX  int
		initialY  int
		advanceX  int
		advanceY  int
		expectedX int
		expectedY int
	}

	cases := []tcase{
		{
			in:        "",
			advanceX:  1,
			expectedX: 0,
			expectedY: 0,
		},
		{
			in:        "",
			advanceY:  1,
			expectedX: 0,
			expectedY: 0,
		},
		{
			in:        "a",
			advanceX:  1,
			expectedX: 1,
			expectedY: 0,
		},
		{
			in:        "a",
			advanceX:  2,
			expectedX: 1,
			expectedY: 0,
		},
		{
			in:        "a",
			advanceX:  -1,
			expectedX: 0,
			expectedY: 0,
		},
		{
			in:        "a\n",
			advanceX:  1,
			expectedX: 1,
			expectedY: 0,
		},
		{
			in:        "a\n",
			advanceX:  1,
			initialX:  1,
			initialY:  0,
			expectedX: 0,
			expectedY: 1,
		},
		{
			in:        "a\n",
			advanceX:  2,
			expectedX: 0,
			expectedY: 1,
		},
		{
			in:        "a\n",
			advanceY:  1,
			expectedX: 0,
			expectedY: 1,
		},
		{
			in:        "a\na",
			advanceX:  1,
			initialX:  0,
			initialY:  1,
			expectedX: 1,
			expectedY: 1,
		},
		{
			in:        "Hello\naaa",
			advanceX:  1,
			initialX:  5,
			initialY:  0,
			expectedX: 0,
			expectedY: 1,
		},
		{
			in:        "Hello\naaa",
			advanceX:  1,
			initialX:  2,
			initialY:  1,
			expectedX: 3,
			expectedY: 1,
		},
	}

	for _, c := range cases {
		b, err := buffer.NewBytes(strings.NewReader(c.in), &buffer.UTF8Encoding{})
		assert.Nil(t, err)
		v := NewBuffer(b)
		v.Attach(state.New())
		v.CursorSetX(c.initialX)
		v.CursorSetY(c.initialY)
		v.CursorAdvance(c.advanceX, c.advanceY)
		assert.Equal(t, c.expectedX, v.CursorX(), fmt.Sprintf("%+v", c))
		assert.Equal(t, c.expectedY, v.CursorY(), fmt.Sprintf("%+v", c))
		co := v.CursorOffset()
		x, y := v.OffsetToCursor(co)
		assert.Equal(t, c.expectedX, x, "cursor offset: %+v", c)
		assert.Equal(t, c.expectedY, y, "cursor offset: %+v", c)
	}
}
