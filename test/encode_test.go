package test

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/ab36245/go-msgpack"
	"github.com/ab36245/go-writer"
)

func TestEncodeArrayLengths(t *testing.T) {
	run := func(n int, e string) {
		l1 := uint32(n)
		me := msgpack.NewEncoder()
		me.PutArrayLength(l1)
		v := me.Bytes()
		encodeTest(t, v, e, false)
	}

	// fixarray
	run(10, `
	    |1 bytes
		|    0000 9a
	`)

	// 16 bit
	run(30000, `
		|3 bytes
		|    0000 dc 75 30
	`)

	// 32bit
	run(80000, `
		|5 bytes
		|    0000 dd 00 01 38 80
	`)
}

func TestEncodeBools(t *testing.T) {
	run := func(b bool, e string) {
		me := msgpack.NewEncoder()
		me.PutBool(b)
		v := me.Bytes()
		encodeTest(t, v, e, false)
	}

	run(true, `
		|1 bytes
		|    0000 c3
	`)
	run(false, `
		|1 bytes
		|    0000 c2
	`)
}

func TestEncodeBytes(t *testing.T) {
	run := func(n int, e string) {
		b := make([]byte, n)
		me := msgpack.NewEncoder()
		me.PutBytes(b)
		v := me.Bytes()
		encodeTest(t, v, e, true)

	}
	// 8 bit length
	run(240, `
		|242 bytes
		|    0000 c4 f0 00 00
	`)

	// 16 bit length
	run(2400, `
		|2403 bytes
		|    0000 c5 09 60 00 00
	`)

	// 32 bit length
	run(240000, `
		|240005 bytes
		|    0000 c6 00 03 a9 80 00 00
	`)

	// Over 32 bit length
	{
		b := make([]byte, 4294967296)
		me := msgpack.NewEncoder()
		v := me.PutBytes(b)
		e := "byte slice (4294967296 bytes) is too long to encode"
		if v == nil || v.Error() != e {
			encodeFail(t, v.Error(), e)
		}
	}
}

func TestEncodeFloat32s(t *testing.T) {
	run := func(f float32, e string) {
		me := msgpack.NewEncoder()
		me.PutFloat32(f)
		v := me.Bytes()
		encodeTest(t, v, e, false)
	}

	run(85.125, `
		|5 bytes
		|    0000 ca 42 aa 40 00
	`)
	run(85.3, `
		|5 bytes
		|    0000 ca 42 aa 99 9a
	`)
	run(0.00085125, `
		|5 bytes
		|    0000 ca 3a 5f 26 6c
	`)
}

func TestEncodeFloat64s(t *testing.T) {
	run := func(f float64, e string) {
		me := msgpack.NewEncoder()
		me.PutFloat64(f)
		v := me.Bytes()
		encodeTest(t, v, e, false)
	}

	run(85.125, `
		|9 bytes
		|    0000 cb 40 55 48 00 00 00 00 00
	`)
	run(85.3, `
		|9 bytes
		|    0000 cb 40 55 53 33 33 33 33 33
	`)
	run(0.00085125, `
		|9 bytes
		|    0000 cb 3f 4b e4 cd 74 92 79 14
	`)
}

func TestEncodeInts(t *testing.T) {
	run := func(i int64, e string) {
		me := msgpack.NewEncoder()
		me.PutInt(i)
		v := me.Bytes()
		encodeTest(t, v, e, false)
	}

	// fixint
	run(69, `
		|1 bytes
		|    0000 45
	`)
	run(-11, `
		|1 bytes
		|    0000 f5
	`)

	// 8 bit
	run(-42, `
		|2 bytes
		|    0000 d0 d6
	`)

	// 16 bit
	run(259, `
		|3 bytes
		|    0000 d1 01 03
	`)
	run(-259, `
		|3 bytes
		|    0000 d1 fe fd
	`)

	// 32 bit
	run(65538, `
		|5 bytes
		|    0000 d2 00 01 00 02
	`)
	run(-65538, `
		|5 bytes
		|    0000 d2 ff fe ff fe
	`)

	// 64 bit
	run(4294967299, `
		|9 bytes
		|    0000 d3 00 00 00 01 00 00 00 03
	`)
	run(-4294967299, `
		|9 bytes
		|    0000 d3 ff ff ff fe ff ff ff fd
	`)
}

func TestEncodeMapLengths(t *testing.T) {
	run := func(n int, e string) {
		l := uint32(n)
		me := msgpack.NewEncoder()
		me.PutMapLength(l)
		v := me.Bytes()
		encodeTest(t, v, e, false)
	}

	// fixmap
	run(10, `
		|1 bytes
		|    0000 8a
	`)

	// 16 bit
	run(30000, `
		|3 bytes
		|    0000 de 75 30
	`)

	// 32bit
	run(80000, `
		|5 bytes
		|    0000 df 00 01 38 80
	`)
}

func TestEncodeNil(t *testing.T) {
	me := msgpack.NewEncoder()
	me.PutNil()
	v := me.Bytes()
	e := `
			|1 bytes
			|    0000 c0
		`
	encodeTest(t, v, e, false)
}

func TestEncodeStrings(t *testing.T) {
	base := "hi!"
	run := func(n int, e string) {
		s := strings.Repeat(base, n)
		me := msgpack.NewEncoder()
		me.PutString(s)
		v := me.Bytes()
		encodeTest(t, v, e, true)
	}

	// fixstr
	run(10, `
		|31 bytes
		|    0000 be 68 69 21
	`)

	run(80, `
		|242 bytes
		|    0000 d9 f0 68 69 21
	`)

	// 16 bit length
	run(800, `
		|2403 bytes
		|    0000 da 09 60 68 69 21
	`)

	// 32 bit length
	run(80000, `
		|240005 bytes
		|    0000 db 00 03 a9 80 68 69 21
	`)

	// Over 32 bit length
	{
		s := strings.Repeat("a", 4294967296)
		me := msgpack.NewEncoder()
		v := me.PutString(s)
		e := "string (4294967296 bytes) is too long to encode"
		if v == nil || v.Error() != e {
			encodeFail(t, v.Error(), e)
		}
	}
}

func TestEncodeTimes(t *testing.T) {
	loc, _ := time.LoadLocation("UTC")
	run := func(d time.Time, e string) {
		me := msgpack.NewEncoder()
		me.PutTime(d)
		v := me.Bytes()
		encodeTest(t, v, e, false)
	}

	// timestamp32
	run(time.Date(1997, 8, 28, 0, 0, 0, 0, loc), `
		|6 bytes
		|    0000 d6 ff 34 04 bf 80
	`)

	// timestamp64
	run(time.Date(1995, 9, 12, 0, 0, 0, 420, loc), `
		|10 bytes
		|    0000 d7 ff 00 00 06 90 30 54 cd 80
	`)

	// timestamp96
	run(time.Date(1961, 10, 19, 0, 0, 0, 420, loc), `
		|15 bytes
		|    0000 c7 0c ff 00 00 01 a4 ff ff ff ff f0 92 32 00
	`)
}

func TestEncodeUints(t *testing.T) {
	run := func(i uint64, e string) {
		me := msgpack.NewEncoder()
		me.PutUint(i)
		v := me.Bytes()
		encodeTest(t, v, e, false)
	}

	// fixint
	run(69, `
		|1 bytes
		|    0000 45
	`)

	// 8 bit
	run(130, `
		|2 bytes
		|    0000 cc 82
	`)

	// 16 bit
	run(259, `
		|3 bytes
		|    0000 cd 01 03
	`)

	// 32 bit
	run(65538, `
		|5 bytes
		|    0000 ce 00 01 00 02
	`)

	// 64 bit
	run(4294967299, `
		|9 bytes
		|    0000 cf 00 00 00 01 00 00 00 03
	`)
}

func encodeTest(t *testing.T, v []byte, e string, prefix bool) {
	vs := writer.Value(v)
	es := writer.Trim(e)
	if prefix {
		vs = vs[0:len(es)]
	}
	if vs == es {
		return
	}
	encodeFail(t, vs, es)
}

func encodeFail(t *testing.T, v, e string) {
	es := fmt.Sprintf("%-8s %s", "expected", e)
	vs := fmt.Sprintf("%-8s %s", "actual", v)
	t.Fatalf("\n%s\n%s", es, vs)
}
