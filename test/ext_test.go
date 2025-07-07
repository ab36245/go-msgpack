package test

import (
	"testing"

	"github.com/ab36245/go-msgpack"
)

func TestExtUint(t *testing.T) {
	run := func(t *testing.T, ti uint8, i uint64, e string) {
		mpe := msgpack.NewEncoder()
		mpe.PutExtUint(ti, i)
		mps := mpe.AsString(-1)
		if mps != e {
			report(t, mps, e)
		}
		mpd := msgpack.NewDecoder(mpe.Bytes())
		ta, a, _ := mpd.GetExtUint()
		if ta != ti {
			report(t, ta, ti)
		}
		if a != i {
			report(t, a, i)
		}
	}

	t.Run("fixext1", func(t *testing.T) {
		t.Run("max", func(t *testing.T) {
			run(t, 42, 255, "d4 2a ff")
		})
	})

	t.Run("fixext2", func(t *testing.T) {
		t.Run("min", func(t *testing.T) {
			run(t, 42, 256, "d5 2a 01 00")
		})
		t.Run("max", func(t *testing.T) {
			run(t, 42, 65535, "d5 2a ff ff")
		})
	})

	t.Run("32 bit", func(t *testing.T) {
		t.Run("min", func(t *testing.T) {
			run(t, 42, 65536, "d6 2a 00 01 00 00")
		})
		t.Run("max", func(t *testing.T) {
			run(t, 42, 4294967295, "d6 2a ff ff ff ff")
		})
	})

	t.Run("64 bit", func(t *testing.T) {
		t.Run("min", func(t *testing.T) {
			run(t, 42, 4294967296, "d7 2a 00 00 00 01 00 00 00 00")
		})
	})
}
