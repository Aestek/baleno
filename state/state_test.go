package state

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestStateSetGet(t *testing.T) {
	s := New()

	s.Set("a.b.c", 1)
	v, ok := s.Get("a.b.c")
	require.True(t, ok)
	require.Equal(t, 1, v)

	s.Set("a.b", 2)
	v, ok = s.Get("a.b")
	require.True(t, ok)
	require.Equal(t, 2, v)
}

func TestStateAlias(t *testing.T) {
	testAlias := func(real, alias string) {
		s := New()

		s.Set(real, 42)
		err := s.Alias(alias, real)
		require.Nil(t, err, alias)

		v, ok := s.Get(alias)
		require.True(t, ok, alias)
		require.Equal(t, 42, v, alias)

		s.Set(real, 43)
		v, ok = s.Get(alias)
		require.True(t, ok, alias)
		require.Equal(t, 43, v, alias)

		s.Set(alias, 44)
		v, ok = s.Get(real)
		require.True(t, ok, alias)
		require.Equal(t, 44, v, alias)
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
	require.Nil(t, err)

	s.Set("a.b.c", 43)
	evt := <-c
	require.Equal(t, "a.b.c", evt.Key)
	require.Equal(t, 43, evt.Value)

	c = make(chan Event, 1)
	err = s.Watch("a.b.c", c)
	require.Nil(t, err)

	s.Set("a.b.c", 44)
	evt = <-c
	require.Equal(t, "a.b.c", evt.Key)
	require.Equal(t, 44, evt.Value)
}

func TestStateWatchAlias(t *testing.T) {
	s := New()

	s.Set("a.b.c", 42)
	err := s.Alias("alias", "a.b.c")
	require.Nil(t, err)

	c := make(chan Event, 1)
	err = s.Watch("alias", c)
	require.Nil(t, err)

	s.Set("a.b.c", 43)
	evt := <-c
	require.Equal(t, "a.b.c", evt.Key)
	require.Equal(t, 43, evt.Value)

	s.Set("a.b.d", "lol")
	err = s.Alias("alias", "a.b.d")
	require.Nil(t, err)

	s.Set("a.b.d", "lul")
	evt = <-c
	require.Equal(t, "a.b.d", evt.Key)
	require.Equal(t, "lul", evt.Value)
}

func TestNamespace(t *testing.T) {
	s := New()

	s.Set("a.b.c", 42)

	s2 := s.Namespace("a")
	require.Equal(t, 42, s2.MustGet("b.c"))

	s.Set("a.b.c", 43)
	require.Equal(t, 43, s2.MustGet("b.c"))
}
