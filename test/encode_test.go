package test

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/ab36245/go-msgpack"
	"github.com/ab36245/go-writer"
)

func TestEncodeInts(t *testing.T) {
	// fixint
	{
		me := msgpack.NewEncoder()
		me.PutInt(69)
		me.PutInt(-11)
		v := me.Bytes()
		e := `
			|2 bytes
			|    0000 45 f5
		`
		encodeTest(t, v, e)
	}

	// 8 bit
	{
		me := msgpack.NewEncoder()
		me.PutInt(130)
		me.PutInt(-42)
		v := me.Bytes()
		e := `
			|4 bytes
			|    0000 d0 82 d0 d6
		`
		encodeTest(t, v, e)
	}

	// 16 bit
	{
		me := msgpack.NewEncoder()
		me.PutInt(259)
		me.PutInt(-259)
		v := me.Bytes()
		e := `
			|6 bytes
			|    0000 d1 01 03 d1 fe fd
		`
		encodeTest(t, v, e)
	}

	// 32 bit
	{
		me := msgpack.NewEncoder()
		me.PutInt(65538)
		me.PutInt(-65538)
		v := me.Bytes()
		e := `
			|10 bytes
			|    0000 d2 00 01 00 02 d2 ff fe ff fe
		`
		encodeTest(t, v, e)
	}

	// 64 bit
	{
		me := msgpack.NewEncoder()
		me.PutInt(4294967299)
		me.PutInt(-4294967299)
		v := me.Bytes()
		e := `
			|18 bytes
			|    0000 d3 00 00 00 01 00 00 00 03 d3 ff ff ff fe ff ff
			|    0016 ff fd
		`
		encodeTest(t, v, e)
	}
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
	// fixstr
	{
		s := "short"
		me := msgpack.NewEncoder()
		me.PutString(s)
		v := me.Bytes()
		e := `
			|6 bytes
			|    0000 a5 73 68 6f 72 74
		`
		encodeTest(t, v, e)
	}

	// 8 bit length
	{
		s := "short, but longer than the minimum 31 bytes for fixstr"
		me := msgpack.NewEncoder()
		me.PutString(s)
		v := me.Bytes()
		e := `
			|56 bytes
            |    0000 d9 36 73 68 6f 72 74 2c 20 62 75 74 20 6c 6f 6e
            |    0016 67 65 72 20 74 68 61 6e 20 74 68 65 20 6d 69 6e
            |    0032 69 6d 75 6d 20 33 31 20 62 79 74 65 73 20 66 6f
            |    0048 72 20 66 69 78 73 74 72
		`
		encodeTest(t, v, e)
	}

	// 16 bit length
	{
		s := strings.Repeat("16-bit!", 40)
		me := msgpack.NewEncoder()
		me.PutString(s)
		v := me.Bytes()
		e := `
			|283 bytes
            |    0000 da 01 18 31 36 2d 62 69 74 21 31 36 2d 62 69 74
            |    0016 21 31 36 2d 62 69 74 21 31 36 2d 62 69 74 21 31
            |    0032 36 2d 62 69 74 21 31 36 2d 62 69 74 21 31 36 2d
            |    0048 62 69 74 21 31 36 2d 62 69 74 21 31 36 2d 62 69
            |    0064 74 21 31 36 2d 62 69 74 21 31 36 2d 62 69 74 21
            |    0080 31 36 2d 62 69 74 21 31 36 2d 62 69 74 21 31 36
            |    0096 2d 62 69 74 21 31 36 2d 62 69 74 21 31 36 2d 62
            |    0112 69 74 21 31 36 2d 62 69 74 21 31 36 2d 62 69 74
            |    0128 21 31 36 2d 62 69 74 21 31 36 2d 62 69 74 21 31
            |    0144 36 2d 62 69 74 21 31 36 2d 62 69 74 21 31 36 2d
            |    0160 62 69 74 21 31 36 2d 62 69 74 21 31 36 2d 62 69
            |    0176 74 21 31 36 2d 62 69 74 21 31 36 2d 62 69 74 21
            |    0192 31 36 2d 62 69 74 21 31 36 2d 62 69 74 21 31 36
            |    0208 2d 62 69 74 21 31 36 2d 62 69 74 21 31 36 2d 62
            |    0224 69 74 21 31 36 2d 62 69 74 21 31 36 2d 62 69 74
            |    0240 21 31 36 2d 62 69 74 21 31 36 2d 62 69 74 21 31
            |    0256 36 2d 62 69 74 21 31 36 2d 62 69 74 21 31 36 2d
            |    0272 62 69 74 21 31 36 2d 62 69 74 21
		`
		encodeTest(t, v, e)
	}

	// 32 bit length
	// Too long to test here
	// See decoder tests

	// Over 32 bit length
	{
		s := strings.Repeat("a", 4294967296)
		me := msgpack.NewEncoder()
		v := me.PutString(s)
		e := "string (4294967296 bytes) is too long to encode"
		if v == nil || v.Error() != e {
			es := fmt.Sprintf("%-8s %s", "expected", e)
			vs := fmt.Sprintf("%-8s %s", "actual", v)
			t.Fatalf("\n%s\n%s", es, vs)
		}
	}
}

func TestEncodeTime(t *testing.T) {
	loc, _ := time.LoadLocation("UTC")

	// timestamp32
	{
		d := time.Date(1997, 8, 28, 0, 0, 0, 0, loc)
		me := msgpack.NewEncoder()
		me.PutTime(d)
		v := me.Bytes()
		e := `
			|6 bytes
			|    0000 d6 ff 34 04 bf 80
		`
		encodeTest(t, v, e)
	}

	// timestamp64
	{
		d := time.Date(1995, 9, 12, 0, 0, 0, 420, loc)
		me := msgpack.NewEncoder()
		me.PutTime(d)
		v := me.Bytes()
		e := `
			|10 bytes
            |    0000 d7 ff 00 00 06 90 30 54 cd 80
		`
		encodeTest(t, v, e)
	}

	// timestamp96
	{
		d := time.Date(1961, 10, 19, 0, 0, 0, 420, loc)
		me := msgpack.NewEncoder()
		me.PutTime(d)
		v := me.Bytes()
		e := `
			|15 bytes
            |    0000 c7 0c ff 00 00 01 a4 ff ff ff ff f0 92 32 00
		`
		encodeTest(t, v, e)
	}
}

func encodeTest(t *testing.T, v []byte, e string) {
	vs := writer.Value(v)
	es := writer.Trim(e)
	if vs == es {
		return
	}
	es = fmt.Sprintf("%-8s %s", "expected", es)
	vs = fmt.Sprintf("%-8s %s", "actual", vs)
	t.Fatalf("\n%s\n%s", es, vs)
}
