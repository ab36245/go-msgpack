package test

import (
	"strings"
	"testing"

	"github.com/ab36245/go-msgpack"
)

func TestString(t *testing.T) {
	base := "*"
	run := func(t *testing.T, n int, e string) {
		s := strings.Repeat(base, n)
		mpe := msgpack.NewEncoder()
		mpe.PutString(s)
		mps := (mpe.AsString(10))[:len(e)]
		if mps != e {
			report(t, mps, e)
		}
		mpd := msgpack.NewDecoder(mpe.Bytes())
		a, _ := mpd.GetString()
		if a != s {
			report(t, a, s)
		}
	}

	t.Run("fixstr", func(t *testing.T) {
		t.Run("max", func(t *testing.T) {
			run(t, 31, "bf 2a 2a")
		})
	})

	t.Run("8 bit length", func(t *testing.T) {
		t.Run("min", func(t *testing.T) {
			run(t, 32, "d9 20 2a 2a")
		})
		t.Run("max", func(t *testing.T) {
			run(t, 255, "d9 ff 2a 2a")
		})
	})

	t.Run("16 bit length", func(t *testing.T) {
		t.Run("min", func(t *testing.T) {
			run(t, 256, "da 01 00 2a 2a")
		})
		t.Run("max", func(t *testing.T) {
			run(t, 65535, "da ff ff 2a 2a")
		})
	})

	// The next two tests will fail on machines with insufficient memory

	// t.Run("32 bit length", func(t *testing.T) {
	// 	t.Run("min", func(t *testing.T) {
	// 		run(t, 65536, "db 00 01 00 00 2a 2a")
	// 	})
	// 	t.Run("max", func(t *testing.T) {
	// 		run(t, 4294967295, "db ff ff ff ff 2a 2a")
	// 	})
	// })

	// t.Run("too big", func(t *testing.T) {
	// 	a := "expected an error but didn't get one"
	// 	mpe := msgpack.NewEncoder()
	// 	err := mpe.PutString(strings.Repeat(base, 4294967296))
	// 	if err != nil {
	// 		a = err.Error()
	// 	}
	// 	e := "string (4294967296 bytes) is too long to encode"
	// 	if a != e {
	// 		report(t, a, e)
	// 	}
	// })
}
