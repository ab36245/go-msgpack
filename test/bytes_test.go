package test

import (
	"testing"

	"github.com/ab36245/go-msgpack"
)

func TestBytes(t *testing.T) {
	run := func(t *testing.T, n int, e string) {
		b := make([]byte, n)
		mpe := msgpack.NewEncoder()
		mpe.PutBytes(b)
		mps := (mpe.AsString(10))[:len(e)]
		if mps != e {
			report(t, mps, e)
		}
		mpd := msgpack.NewDecoder(mpe.Bytes())
		a, _ := mpd.GetBytes()
		if len(a) != n {
			report(t, len(a), n)
		}
	}

	t.Run("8 bit length", func(t *testing.T) {
		t.Run("max", func(t *testing.T) {
			run(t, 255, "c4 ff 00 00")
		})
	})

	t.Run("16 bit length", func(t *testing.T) {
		t.Run("min", func(t *testing.T) {
			run(t, 256, "c5 01 00 00 00")
		})
		t.Run("max", func(t *testing.T) {
			run(t, 65535, "c5 ff ff 00 00")
		})
	})

	t.Run("32 bit length", func(t *testing.T) {
		t.Run("min", func(t *testing.T) {
			run(t, 65536, "c6 00 01 00 00 00 00")
		})
		t.Run("max", func(t *testing.T) {
			run(t, 4294967295, "c6 ff ff ff ff 00 00")
		})
	})

	t.Run("too big", func(t *testing.T) {
		a := "expected an error but didn't get one"
		mpe := msgpack.NewEncoder()
		err := mpe.PutBytes(make([]byte, 4294967296))
		if err != nil {
			a = err.Error()
		}
		e := "byte slice (4294967296 bytes) is too long to encode"
		if a != e {
			report(t, a, e)
		}
	})
}
