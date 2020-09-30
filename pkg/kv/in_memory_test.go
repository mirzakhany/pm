package kv

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInMemory(t *testing.T) {
	i := Memory()
	i.SetString("a", "b")

	s, ok := i.Get("a")
	assert.True(t, ok)
	assert.Equal(t, s, "b")

	s, ok = i.Get("aa")
	assert.False(t, ok)
	assert.Empty(t, s)
}
