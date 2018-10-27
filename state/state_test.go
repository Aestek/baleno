package state

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStateSetGet(t *testing.T) {
	s := New()

	s.Set("a.b.c", 1)
	v, ok := s.Get("a.b.c")
	assert.True(t, ok)
	assert.Equal(t, 1, v)

	s.Set("a.b", 2)
	v, ok = s.Get("a.b")
	assert.True(t, ok)
	assert.Equal(t, 2, v)
}

func TestStateAlias(t *testing.T) {
	testAlias := func(real, alias string) {
		s := New()

		s.Set(real, 42)
		err := s.Alias(alias, real)
		assert.Nil(t, err, alias)

		v, ok := s.Get(alias)
		assert.True(t, ok, alias)
		assert.Equal(t, 42, v, alias)

		s.Set(real, 43)
		v, ok = s.Get(alias)
		assert.True(t, ok, alias)
		assert.Equal(t, 43, v, alias)

		s.Set(alias, 44)
		v, ok = s.Get(real)
		assert.True(t, ok, alias)
		assert.Equal(t, 44, v, alias)
	}

	testAlias("a.b.c", "x.y.z")
	testAlias("a.b.c", "x")
	testAlias("a.b.c", "x.y.z.0")
}

func TestStateWatch(t *testing.T) {
	s := New()

	s.Set("a.b.c", 42)

	c := make(chan Event, 1)
	err := s.Watch("a.b", c)
	assert.Nil(t, err)

	s.Set("a.b.c", 43)
	evt := <-c
	assert.Equal(t, "a.b.c", evt.Key)
	assert.Equal(t, 43, evt.Value)
}
