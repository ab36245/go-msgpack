package test

import (
	"testing"
	"time"

	"github.com/ab36245/go-msgpack"
	"github.com/ab36245/go-writer"
)

func TestEncodeInts(t *testing.T) {
	me := msgpack.NewEncoder()
	me.PutInt(0x45)
	me.PutInt(0x1A4)
	v := me.Bytes()
	e := `
		|4 bytes
    	|    0000 45 cd 01 a4
	`
	encodeTest(t, v, e)
}

func TestEncodeNil(t *testing.T) {
	me := msgpack.NewEncoder()
	me.PutInt(69)
	me.PutNil()
	me.PutString("follows nil")
	v := me.Bytes()
	e := `
		|14 bytes
        |    0000 45 c0 ab 66 6f 6c 6c 6f 77 73 20 6e 69 6c
	`
	encodeTest(t, v, e)
}

func TestEncodeString(t *testing.T) {
	me := msgpack.NewEncoder()
	me.PutString("malaka")
	v := me.Bytes()
	e := `
		|7 bytes
    	|    0000 a6 6d 61 6c 61 6b 61
	`
	encodeTest(t, v, e)
}

func TestEncodeTime(t *testing.T) {
	mp := msgpack.NewEncoder()

	when, _ := time.Parse(time.DateOnly, "1997-08-28")
	secs := when.UnixMilli() / 1000
	mp.PutInt(secs)
	mp.PutTime(when)

	v := mp.Bytes()
	e := `
		|11 bytes
        |    0000 ce 34 04 bf 80 d6 ff 34 04 bf 80
	`
	encodeTest(t, v, e)
}

func encodeTest(t *testing.T, v any, e string) {
	s := writer.Value(v)
	e = writer.Trim(e)
	if s != e {
		t.Fatalf("\nexpected %s\nactual %s", e, s)
	}
}
