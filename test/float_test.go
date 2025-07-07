package test

import (
	"testing"

	"github.com/ab36245/go-msgpack"
)

func TestFloat32(t *testing.T) {
	run := func(t *testing.T, f float32, e string) {
		mpe := msgpack.NewEncoder()
		mpe.PutFloat32(f)
		mps := mpe.AsString(-1)
		if mps != e {
			report(t, mps, e)
		}
		mpd := msgpack.NewDecoder(mpe.Bytes())
		a, _ := mpd.GetFloat()
		if float32(a) != f {
			report(t, a, f)
		}
	}

	t.Run("85.125", func(t *testing.T) {
		run(t, 85.125, "ca 42 aa 40 00")
	})

	t.Run("85.3", func(t *testing.T) {
		run(t, 85.3, "ca 42 aa 99 9a")
	})

	t.Run("0.00085125", func(t *testing.T) {
		run(t, 0.00085125, "ca 3a 5f 26 6c")
	})
}

func TestFloat64(t *testing.T) {
	run := func(t *testing.T, f float64, e string) {
		mpe := msgpack.NewEncoder()
		mpe.PutFloat64(f)
		mps := mpe.AsString(-1)
		if mps != e {
			report(t, mps, e)
		}
		mpd := msgpack.NewDecoder(mpe.Bytes())
		a, _ := mpd.GetFloat()
		if a != f {
			report(t, a, f)
		}
	}

	t.Run("85.125", func(t *testing.T) {
		run(t, 85.125, "cb 40 55 48 00 00 00 00 00")
	})

	t.Run("85.3", func(t *testing.T) {
		run(t, 85.3, "cb 40 55 53 33 33 33 33 33")
	})

	t.Run("0.00085125", func(t *testing.T) {
		run(t, 0.00085125, "cb 3f 4b e4 cd 74 92 79 14")
	})
}
