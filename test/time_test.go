package test

import (
	"testing"
	"time"

	"github.com/ab36245/go-msgpack"
)

func TestTime(t *testing.T) {
	loc, _ := time.LoadLocation("UTC")
	run := func(t *testing.T, d time.Time, e string) {
		mpe := msgpack.NewEncoder()
		mpe.PutTime(d)
		mps := mpe.AsString(-1)
		if mps != e {
			report(t, mps, e)
		}
		mpd := msgpack.NewDecoder(mpe.Bytes())
		a, _ := mpd.GetTime()
		if a != d {
			report(t, a, d)
		}
	}

	t.Run("timestamp32", func(t *testing.T) {
		run(t,
			time.Date(1997, 8, 28, 0, 0, 0, 0, loc),
			"d6 ff 34 04 bf 80",
		)
	})

	t.Run("timestamp64", func(t *testing.T) {
		run(t,
			time.Date(1995, 9, 12, 0, 0, 0, 420000, loc),
			"d7 ff 00 19 a2 80 30 54 cd 80",
		)
	})

	t.Run("timestamp96", func(t *testing.T) {
		run(t,
			time.Date(1961, 10, 19, 0, 0, 0, 420000, loc),
			"c7 0c ff 00 06 68 a0 ff ff ff ff f0 92 32 00",
		)
	})
}
