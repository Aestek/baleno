package view

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLinesEmpty(t *testing.T) {
	called := false
	lines([]rune{}, func(int, []rune) bool {
		called = true
		return true
	})
	assert.False(t, called)
}

func TestOneLine(t *testing.T) {
	calledLines := [][]rune{}
	calledLine := -1
	lines([]rune("Hello!"), func(l int, c []rune) bool {
		calledLines = append(calledLines, c)
		calledLine = l
		return true
	})
	assert.Equal(t, [][]rune{
		[]rune("Hello!"),
	}, calledLines)
	assert.Equal(t, 0, calledLine)
}

func TestTwoLinesLastEmpty(t *testing.T) {
	calledLines := [][]rune{}
	calledLine := -1
	lines([]rune("Hello!\n"), func(l int, c []rune) bool {
		calledLines = append(calledLines, c)
		calledLine = l
		return true
	})
	assert.Equal(t, [][]rune{
		[]rune("Hello!"),
		[]rune{},
	}, calledLines)
	assert.Equal(t, 1, calledLine)
}

func TestTwoLine(t *testing.T) {
	calledLines := [][]rune{}
	calledLineNb := []int{}
	lines([]rune("HEY\nHello!"), func(l int, c []rune) bool {
		calledLines = append(calledLines, c)
		calledLineNb = append(calledLineNb, l)
		return true
	})
	assert.Equal(t, [][]rune{
		[]rune("HEY"),
		[]rune("Hello!"),
	}, calledLines)
	assert.Equal(t, []int{0, 1}, calledLineNb)
}
