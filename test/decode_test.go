package test

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/ab36245/go-msgpack"
)

func TestDecodeArrayLengths(t *testing.T) {
	run := func(n int) {
		l := uint32(n)
		me := msgpack.NewEncoder()
		me.PutArrayLength(l)
		md := msgpack.NewDecoder(me.Bytes())
		v, _ := md.GetArrayLength()
		decodeTest(t, v, l)
	}

	run(10)    // fixarray
	run(30000) // 16 bit
	run(80000) // 32bit
}

func TestDecodeBools(t *testing.T) {
	run := func(b bool) {
		me := msgpack.NewEncoder()
		me.PutBool(b)
		md := msgpack.NewDecoder(me.Bytes())
		v, _ := md.GetBool()
		decodeTest(t, v, b)
	}

	run(true)
	run(false)
}

func TestDecodeBytes(t *testing.T) {
	run := func(n int) {
		b := make([]byte, n)
		me := msgpack.NewEncoder()
		me.PutBytes(b)
		md := msgpack.NewDecoder(me.Bytes())
		v, _ := md.GetBytes()
		decodeTest(t, len(v), len(b))
	}

	run(240)    // 8 bit length
	run(2400)   // 16 bit length
	run(240000) // 32 bit length
}

func TestDecodeFloat32s(t *testing.T) {
	run := func(f float32) {
		me := msgpack.NewEncoder()
		me.PutFloat32(f)
		md := msgpack.NewDecoder(me.Bytes())
		v, _ := md.GetFloat()
		decodeTest(t, v, float64(f))
	}

	run(85.125)
	run(85.3)
	run(0.00085125)
}

func TestDecodeFloat64s(t *testing.T) {
	run := func(f float64) {
		me := msgpack.NewEncoder()
		me.PutFloat64(f)
		md := msgpack.NewDecoder(me.Bytes())
		v, _ := md.GetFloat()
		decodeTest(t, v, f)
	}

	run(85.125)
	run(85.3)
	run(0.00085125)
}

func TestDecodeInts(t *testing.T) {
	run := func(i int64) {
		me := msgpack.NewEncoder()
		me.PutInt(i)
		md := msgpack.NewDecoder(me.Bytes())
		v, _ := md.GetInt()
		decodeTest(t, v, i)
	}

	run(69)          // fixint
	run(-11)         // fixint
	run(-42)         // 8 bit
	run(259)         // 16 bit
	run(-259)        // 16 bit
	run(65538)       // 32 bit
	run(-65538)      // 32 bit
	run(4294967299)  // 64 bit
	run(-4294967299) // 64 bit
}

func TestDecodeMapLengths(t *testing.T) {
	run := func(n int) {
		l := uint32(n)
		me := msgpack.NewEncoder()
		me.PutMapLength(l)
		md := msgpack.NewDecoder(me.Bytes())
		v, _ := md.GetMapLength()
		decodeTest(t, v, l)

	}
	run(10)    // fixmap
	run(30000) // 16 bit
	run(80000) // 32bit
}

func TestDecodeNil(t *testing.T) {
	{
		me := msgpack.NewEncoder()
		me.PutNil()
		md := msgpack.NewDecoder(me.Bytes())
		v1, _ := md.IsNil()
		decodeTest(t, v1, true)
	}
	{
		me := msgpack.NewEncoder()
		me.PutInt(42)
		md := msgpack.NewDecoder(me.Bytes())
		v1, _ := md.IsNil()
		decodeTest(t, v1, false)
	}
	{
		me := msgpack.NewEncoder()
		me.PutNil()
		me.PutInt(42)
		md := msgpack.NewDecoder(me.Bytes())
		v1, _ := md.IfNil()
		decodeTest(t, v1, true)
		v2, _ := md.GetInt()
		decodeTest(t, v2, int64(42))
	}
	{
		me := msgpack.NewEncoder()
		me.PutInt(42)
		md := msgpack.NewDecoder(me.Bytes())
		v1, _ := md.IfNil()
		decodeTest(t, v1, false)
		v2, _ := md.GetInt()
		decodeTest(t, v2, int64(42))
	}
}

func TestDecodeStrings(t *testing.T) {
	base := "hi!"
	run := func(n int) {
		s := strings.Repeat(base, n)
		me := msgpack.NewEncoder()
		me.PutString(s)
		md := msgpack.NewDecoder(me.Bytes())
		v, _ := md.GetString()
		decodeTest(t, v, s)
	}

	run(10)    // fixstr
	run(80)    // 8 bit length
	run(800)   // 16 bit length
	run(80000) // 32 bit length
}

func TestDecodeTime(t *testing.T) {
	loc, _ := time.LoadLocation("UTC")
	run := func(d time.Time) {
		me := msgpack.NewEncoder()
		me.PutTime(d)
		md := msgpack.NewDecoder(me.Bytes())
		v, _ := md.GetTime()
		decodeTest(t, v, d)
	}

	run(time.Date(1997, 8, 28, 0, 0, 0, 0, loc))    // timestamp32
	run(time.Date(1995, 9, 12, 0, 0, 0, 420, loc))  // timestamp64
	run(time.Date(1961, 10, 19, 0, 0, 0, 420, loc)) // timestamp96
}

func TestDecodeUints(t *testing.T) {
	run := func(u uint64) {
		me := msgpack.NewEncoder()
		me.PutUint(u)
		md := msgpack.NewDecoder(me.Bytes())
		v, _ := md.GetUint()
		decodeTest(t, v, u)
	}

	run(69)         // fixint
	run(130)        // 8 bit
	run(259)        // 16 bit
	run(65538)      // 32 bit
	run(4294967299) // 64 bit
}

func decodeTest(t *testing.T, v any, e any) {
	if v == e {
		return
	}
	es := fmt.Sprintf("%-8s %v", "expected", e)
	vs := fmt.Sprintf("%-8s %v", "actual", v)
	t.Fatalf("\n%s\n%s", es, vs)
}
