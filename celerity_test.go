package celerity

import (
	"testing"
)

func TestSetEnvironment(t *testing.T) {
	{
		SetEnvironment(DEV)
		c := NewContext()
		if c.Env != DEV {
			t.Errorf("environment not set: %s", c.Env)
		}
	}
	{
		SetEnvironment(PROD)
		c := NewContext()
		if c.Env != PROD {
			t.Errorf("environment not set: %s", c.Env)
		}
	}
	SetEnvironment(DEV)
}
