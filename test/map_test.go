package test

import (
	"testing"

	"github.com/ab36245/go-msgpack"
)

func TestMap(t *testing.T) {
	run := func(t *testing.T, n uint32, e string) {
		mpe := msgpack.NewEncoder()
		mpe.PutMapLength(n)
		mps := mpe.AsString(-1)
		if mps != e {
			report(t, mps, e)
		}
		mpd := msgpack.NewDecoder(mpe.Bytes())
		a, _ := mpd.GetMapLength()
		if a != n {
			report(t, a, n)
		}
	}

	t.Run("fixmap", func(t *testing.T) {
		run(t, 15, "8f")
	})

	t.Run("16 bit", func(t *testing.T) {
		t.Run("min", func(t *testing.T) {
			run(t, 16, "de 00 10")
		})
		t.Run("max", func(t *testing.T) {
			run(t, 65535, "de ff ff")
		})
	})

	t.Run("32 bit", func(t *testing.T) {
		t.Run("min", func(t *testing.T) {
			run(t, 65536, "df 00 01 00 00")
		})
		t.Run("max", func(t *testing.T) {
			run(t, 4294967295, "df ff ff ff ff")
		})
	})
}
