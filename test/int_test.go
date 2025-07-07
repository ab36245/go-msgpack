package test

import (
	"testing"

	"github.com/ab36245/go-msgpack"
)

func TestInt(t *testing.T) {
	run := func(t *testing.T, i int64, e string) {
		mpe := msgpack.NewEncoder()
		mpe.PutInt(i)
		mps := mpe.AsString(-1)
		if mps != e {
			report(t, mps, e)
		}
		mpd := msgpack.NewDecoder(mpe.Bytes())
		a, _ := mpd.GetInt()
		if a != i {
			report(t, a, i)
		}
	}

	t.Run("fixint", func(t *testing.T) {
		t.Run("max positive", func(t *testing.T) {
			run(t, 127, "7f")
		})
		t.Run("min negative", func(t *testing.T) {
			run(t, -1, "ff")
		})
		t.Run("max negative", func(t *testing.T) {
			run(t, -32, "e0")
		})
	})

	t.Run("8 bit", func(t *testing.T) {
		t.Run("min negative", func(t *testing.T) {
			run(t, -33, "d0 df")
		})
		t.Run("max negative", func(t *testing.T) {
			run(t, -128, "d0 80")
		})
	})

	t.Run("16 bit", func(t *testing.T) {
		t.Run("min positive", func(t *testing.T) {
			run(t, 128, "d1 00 80")
		})
		t.Run("max positive", func(t *testing.T) {
			run(t, 32767, "d1 7f ff")
		})
		t.Run("min negative", func(t *testing.T) {
			run(t, -129, "d1 ff 7f")
		})
		t.Run("max negative", func(t *testing.T) {
			run(t, -32768, "d1 80 00")
		})
	})

	t.Run("32 bit", func(t *testing.T) {
		t.Run("min positive", func(t *testing.T) {
			run(t, 32768, "d2 00 00 80 00")
		})
		t.Run("max positive", func(t *testing.T) {
			run(t, 2147483647, "d2 7f ff ff ff")
		})
		t.Run("min negative", func(t *testing.T) {
			run(t, -32769, "d2 ff ff 7f ff")
		})
		t.Run("max negative", func(t *testing.T) {
			run(t, -2147483648, "d2 80 00 00 00")
		})
	})

	t.Run("64 bit", func(t *testing.T) {
		t.Run("min positive", func(t *testing.T) {
			run(t, 2147483648, "d3 00 00 00 00 80 00 00 00")
		})
		t.Run("min negative", func(t *testing.T) {
			run(t, -2147483649, "d3 ff ff ff ff 7f ff ff ff")
		})
	})
}
