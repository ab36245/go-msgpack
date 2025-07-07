package test

import (
	"testing"

	"github.com/ab36245/go-msgpack"
)

func TestNil(t *testing.T) {
	t.Run("put", func(t *testing.T) {
		mpe := msgpack.NewEncoder()
		mpe.PutNil()
		mps := mpe.AsString(-1)
		e := "c0"
		if mps != e {
			report(t, mps, e)
		}
	})
}
