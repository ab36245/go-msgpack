package test

import (
	"testing"

	"github.com/ab36245/go-msgpack"
)

func TestArray(t *testing.T) {
	run := func(t *testing.T, n uint32, e string) {
		mpe := msgpack.NewEncoder()
		mpe.PutArrayLength(n)
		mps := mpe.AsString(-1)
		if mps != e {
			report(t, mps, e)
		}
		mpd := msgpack.NewDecoder(mpe.Bytes())
		a, _ := mpd.GetArrayLength()
		if a != n {
			report(t, a, n)
		}
	}

	t.Run("fixarray", func(t *testing.T) {
		run(t, 15, "9f")
	})

	t.Run("16 bit", func(t *testing.T) {
		t.Run("min", func(t *testing.T) {
			run(t, 16, "dc 00 10")
		})
		t.Run("max", func(t *testing.T) {
			run(t, 65535, "dc ff ff")
		})
	})

	t.Run("32 bit", func(t *testing.T) {
		t.Run("min", func(t *testing.T) {
			run(t, 65536, "dd 00 01 00 00")
		})
		t.Run("max", func(t *testing.T) {
			run(t, 4294967295, "dd ff ff ff ff")
		})
	})
}
