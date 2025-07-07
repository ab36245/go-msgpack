package test

import (
	"testing"

	"github.com/ab36245/go-msgpack"
)

func TestBool(t *testing.T) {
	run := func(t *testing.T, b bool, e string) {
		mpe := msgpack.NewEncoder()
		mpe.PutBool(b)
		mps := mpe.AsString(-1)
		if mps != e {
			report(t, mps, e)
		}
		mpd := msgpack.NewDecoder(mpe.Bytes())
		a, _ := mpd.GetBool()
		if a != b {
			report(t, a, b)
		}
	}

	t.Run("false", func(t *testing.T) {
		run(t, false, "c2")
	})
	t.Run("true", func(t *testing.T) {
		run(t, true, "c3")
	})
}
