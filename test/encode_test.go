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
	// fixarray
	{
		l1 := uint32(10)
		me := msgpack.NewEncoder()
		me.PutArrayLength(l1)
		v := me.Bytes()
		e := `
			|1 bytes
			|    0000 9a
		`
		encodeTest(t, v, e, false)
	}

	// 16 bit
	{
		l1 := uint32(30000)
		me := msgpack.NewEncoder()
		me.PutArrayLength(l1)
		v := me.Bytes()
		e := `
		    |3 bytes
            |    0000 dc 75 30
		`
		encodeTest(t, v, e, false)
	}

	// 32bit
	{
		l1 := uint32(80000)
		me := msgpack.NewEncoder()
		me.PutArrayLength(l1)
		v := me.Bytes()
		e := `
			|5 bytes
			|    0000 dd 00 01 38 80
		`
		encodeTest(t, v, e, false)
	}
}

func TestEncodeBools(t *testing.T) {
	{
		me := msgpack.NewEncoder()
		me.PutBool(true)
		me.PutBool(false)
		v := me.Bytes()
		e := `
			|2 bytes
			|    0000 c3 c2
		`
		encodeTest(t, v, e, false)
	}
}

func TestEncodeBytes(t *testing.T) {
	// 8 bit length
	{
		b := make([]byte, 240)
		me := msgpack.NewEncoder()
		me.PutBytes(b)
		v := me.Bytes()
		e := `
			|242 bytes
            |    0000 c4 f0 00 00
		`
		encodeTest(t, v, e, true)
	}

	// 16 bit length
	{
		b := make([]byte, 2400)
		me := msgpack.NewEncoder()
		me.PutBytes(b)
		v := me.Bytes()
		e := `
			|2403 bytes
            |    0000 c5 09 60 00 00
		`
		encodeTest(t, v, e, true)
	}

	// 32 bit length
	{
		b := make([]byte, 240000)
		me := msgpack.NewEncoder()
		me.PutBytes(b)
		v := me.Bytes()
		e := `
			|240005 bytes
            |    0000 c6 00 03 a9 80 00 00
		`
		encodeTest(t, v, e, true)
	}

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

func TestEncodeFloats(t *testing.T) {
	t.Fatal("not implemeted")
}

func TestEncodeInts(t *testing.T) {
	// fixint
	{
		i1 := int64(69)
		i2 := int64(-11)
		me := msgpack.NewEncoder()
		me.PutInt(i1)
		me.PutInt(i2)
		v := me.Bytes()
		e := `
			|2 bytes
			|    0000 45 f5
		`
		encodeTest(t, v, e, false)
	}

	// 8 bit
	{
		i1 := int64(-42)
		me := msgpack.NewEncoder()
		me.PutInt(i1)
		v := me.Bytes()
		e := `
			|2 bytes
			|    0000 d0 d6
		`
		encodeTest(t, v, e, false)
	}

	// 16 bit
	{
		i1 := int64(259)
		i2 := int64(-259)
		me := msgpack.NewEncoder()
		me.PutInt(i1)
		me.PutInt(i2)
		v := me.Bytes()
		e := `
			|6 bytes
			|    0000 d1 01 03 d1 fe fd
		`
		encodeTest(t, v, e, false)
	}

	// 32 bit
	{
		i1 := int64(65538)
		i2 := int64(-65538)
		me := msgpack.NewEncoder()
		me.PutInt(i1)
		me.PutInt(i2)
		v := me.Bytes()
		e := `
			|10 bytes
			|    0000 d2 00 01 00 02 d2 ff fe ff fe
		`
		encodeTest(t, v, e, false)
	}

	// 64 bit
	{
		i1 := int64(4294967299)
		i2 := int64(-4294967299)
		me := msgpack.NewEncoder()
		me.PutInt(i1)
		me.PutInt(i2)
		v := me.Bytes()
		e := `
			|18 bytes
			|    0000 d3 00 00 00 01 00 00 00 03 d3 ff ff ff fe ff ff
			|    0016 ff fd
		`
		encodeTest(t, v, e, false)
	}
}

func TestEncodeMapLengths(t *testing.T) {
	// fixmap
	{
		l1 := uint32(10)
		me := msgpack.NewEncoder()
		me.PutMapLength(l1)
		v := me.Bytes()
		e := `
			|1 bytes
			|    0000 8a
		`
		encodeTest(t, v, e, false)
	}

	// 16 bit
	{
		l1 := uint32(30000)
		me := msgpack.NewEncoder()
		me.PutMapLength(l1)
		v := me.Bytes()
		e := `
		    |3 bytes
            |    0000 de 75 30
		`
		encodeTest(t, v, e, false)
	}

	// 32bit
	{
		l1 := uint32(80000)
		me := msgpack.NewEncoder()
		me.PutMapLength(l1)
		v := me.Bytes()
		e := `
			|5 bytes
			|    0000 df 00 01 38 80
		`
		encodeTest(t, v, e, false)
	}
}

func TestEncodeNil(t *testing.T) {
	t.Fatal("not implemeted")
}

func TestEncodeStrings(t *testing.T) {
	base := "hi!"
	// fixstr
	{
		s := strings.Repeat(base, 10)
		me := msgpack.NewEncoder()
		me.PutString(s)
		v := me.Bytes()
		e := `
			|31 bytes
            |    0000 be 68 69 21
		`
		encodeTest(t, v, e, true)
	}

	// 8 bit length
	{
		s := strings.Repeat(base, 80)
		me := msgpack.NewEncoder()
		me.PutString(s)
		v := me.Bytes()
		e := `
			|242 bytes
            |    0000 d9 f0 68 69 21
		`
		encodeTest(t, v, e, true)
	}

	// 16 bit length
	{
		s := strings.Repeat(base, 800)
		me := msgpack.NewEncoder()
		me.PutString(s)
		v := me.Bytes()
		e := `
			|2403 bytes
            |    0000 da 09 60 68 69 21
		`
		encodeTest(t, v, e, true)
	}

	// 32 bit length
	{
		s := strings.Repeat(base, 80000)
		me := msgpack.NewEncoder()
		me.PutString(s)
		v := me.Bytes()
		e := `
			|240005 bytes
            |    0000 db 00 03 a9 80 68 69 21
		`
		encodeTest(t, v, e, true)
	}

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
		encodeTest(t, v, e, false)
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
		encodeTest(t, v, e, false)
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
		encodeTest(t, v, e, false)
	}
}

func TestEncodeUints(t *testing.T) {
	// fixint
	{
		i1 := uint64(69)
		me := msgpack.NewEncoder()
		me.PutUint(i1)
		v := me.Bytes()
		e := `
			|1 bytes
			|    0000 45
		`
		encodeTest(t, v, e, false)
	}

	// 8 bit
	{
		i1 := uint64(130)
		me := msgpack.NewEncoder()
		me.PutUint(i1)
		v := me.Bytes()
		e := `
			|2 bytes
			|    0000 cc 82
		`
		encodeTest(t, v, e, false)
	}

	// 16 bit
	{
		i1 := uint64(259)
		me := msgpack.NewEncoder()
		me.PutUint(i1)
		v := me.Bytes()
		e := `
			|3 bytes
			|    0000 cd 01 03
		`
		encodeTest(t, v, e, false)
	}

	// 32 bit
	{
		i1 := uint64(65538)
		me := msgpack.NewEncoder()
		me.PutUint(i1)
		v := me.Bytes()
		e := `
			|5 bytes
            |    0000 ce 00 01 00 02
		`
		encodeTest(t, v, e, false)
	}

	// 64 bit
	{
		i1 := uint64(4294967299)
		me := msgpack.NewEncoder()
		me.PutUint(i1)
		v := me.Bytes()
		e := `
		    |9 bytes
            |    0000 cf 00 00 00 01 00 00 00 03
		`
		encodeTest(t, v, e, false)
	}
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
