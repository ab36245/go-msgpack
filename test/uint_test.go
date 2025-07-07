package test

import (
	"testing"

	"github.com/ab36245/go-msgpack"
)

func TestUint(t *testing.T) {
	run := func(t *testing.T, i uint64, e string) {
		mpe := msgpack.NewEncoder()
		mpe.PutUint(i)
		mps := mpe.AsString(-1)
		if mps != e {
			report(t, mps, e)
		}
		mpd := msgpack.NewDecoder(mpe.Bytes())
		a, _ := mpd.GetUint()
		if a != i {
			report(t, a, i)
		}
	}

	t.Run("fixint", func(t *testing.T) {
		t.Run("max", func(t *testing.T) {
			run(t, 127, "7f")
		})
	})

	t.Run("8 bit", func(t *testing.T) {
		t.Run("min", func(t *testing.T) {
			run(t, 128, "cc 80")
		})
		t.Run("max", func(t *testing.T) {
			run(t, 255, "cc ff")
		})
	})

	t.Run("16 bit", func(t *testing.T) {
		t.Run("min", func(t *testing.T) {
			run(t, 256, "cd 01 00")
		})
		t.Run("max", func(t *testing.T) {
			run(t, 65535, "cd ff ff")
		})
	})

	t.Run("32 bit", func(t *testing.T) {
		t.Run("min", func(t *testing.T) {
			run(t, 65536, "ce 00 01 00 00")
		})
		t.Run("max", func(t *testing.T) {
			run(t, 4294967295, "ce ff ff ff ff")
		})
	})

	t.Run("64 bit", func(t *testing.T) {
		t.Run("min", func(t *testing.T) {
			run(t, 4294967296, "cf 00 00 00 01 00 00 00 00")
		})
	})
}
