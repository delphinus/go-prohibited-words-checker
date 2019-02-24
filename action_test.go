package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewWalk(t *testing.T) {
	a := assert.New(t)
	for _, c := range []struct {
		re []string
		ok bool
	}{
		{re: []string{"invalid[re"}, ok: false},
		{re: []string{"valid[re]"}, ok: true},
	} {
		done := prepareValidConfig(t)
		defer done()
		Config.Ignores = c.re
		if _, err := NewWalker(); c.ok {
			a.NoError(err)
		} else {
			a.Error(err)
			t.Logf("found err: %s", err)
		}
	}
}
