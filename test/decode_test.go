package test

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/ab36245/go-msgpack"
)

func TestDecodeArrayLengths(t *testing.T) {
	// fixarray
	{
		l1 := uint32(10)
		me := msgpack.NewEncoder()
		me.PutArrayLength(l1)
		md := msgpack.NewDecoder(me.Bytes())
		v1, _ := md.GetArrayLength()
		decodeTest(t, v1, l1)
	}

	// 16 bit
	{
		l1 := uint32(30000)
		me := msgpack.NewEncoder()
		me.PutArrayLength(l1)
		md := msgpack.NewDecoder(me.Bytes())
		v1, _ := md.GetArrayLength()
		decodeTest(t, v1, l1)
	}

	// 32bit
	{
		l1 := uint32(80000)
		me := msgpack.NewEncoder()
		me.PutArrayLength(l1)
		md := msgpack.NewDecoder(me.Bytes())
		v1, _ := md.GetArrayLength()
		decodeTest(t, v1, l1)
	}
}

func TestDecodeBools(t *testing.T) {
	b1 := true
	b2 := false
	me := msgpack.NewEncoder()
	me.PutBool(b1)
	me.PutBool(b2)
	md := msgpack.NewDecoder(me.Bytes())
	v1, _ := md.GetBool()
	decodeTest(t, v1, b1)
	v2, _ := md.GetBool()
	decodeTest(t, v2, b2)
}

func TestDecodeBytes(t *testing.T) {
	// 8 bit length
	{
		b := make([]byte, 240)
		me := msgpack.NewEncoder()
		me.PutBytes(b)
		md := msgpack.NewDecoder(me.Bytes())
		v, _ := md.GetBytes()
		decodeTest(t, len(v), len(b))
	}

	// 16 bit length
	{
		b := make([]byte, 2400)
		me := msgpack.NewEncoder()
		me.PutBytes(b)
		md := msgpack.NewDecoder(me.Bytes())
		v, _ := md.GetBytes()
		decodeTest(t, len(v), len(b))
	}

	// 32 bit length
	{
		b := make([]byte, 240000)
		me := msgpack.NewEncoder()
		me.PutBytes(b)
		md := msgpack.NewDecoder(me.Bytes())
		v, _ := md.GetBytes()
		decodeTest(t, len(v), len(b))
	}
}

func TestDecodeFloats(t *testing.T) {
	// float 32
	{
		i1 := float32(85.125)
		i2 := float32(85.3)
		i3 := float32(0.00085125)
		me := msgpack.NewEncoder()
		me.PutFloat32(i1)
		me.PutFloat32(i2)
		me.PutFloat32(i3)
		md := msgpack.NewDecoder(me.Bytes())
		v1, _ := md.GetFloat()
		decodeTest(t, v1, float64(i1))
		v2, _ := md.GetFloat()
		decodeTest(t, v2, float64(i2))
		v3, _ := md.GetFloat()
		decodeTest(t, v3, float64(i3))
	}

	// float 64
	{
		i1 := float64(85.125)
		i2 := float64(85.3)
		i3 := float64(0.00085125)
		me := msgpack.NewEncoder()
		me.PutFloat64(i1)
		me.PutFloat64(i2)
		me.PutFloat64(i3)
		md := msgpack.NewDecoder(me.Bytes())
		v1, _ := md.GetFloat()
		decodeTest(t, v1, i1)
		v2, _ := md.GetFloat()
		decodeTest(t, v2, i2)
		v3, _ := md.GetFloat()
		decodeTest(t, v3, i3)
	}
}

func TestDecodeInts(t *testing.T) {
	// fixint
	{
		i1 := int64(69)
		i2 := int64(-11)
		me := msgpack.NewEncoder()
		me.PutInt(i1)
		me.PutInt(i2)
		md := msgpack.NewDecoder(me.Bytes())
		v1, _ := md.GetInt()
		decodeTest(t, v1, i1)
		v2, _ := md.GetInt()
		decodeTest(t, v2, i2)
	}

	// 8 bit
	{
		i1 := int64(-42)
		me := msgpack.NewEncoder()
		me.PutInt(i1)
		md := msgpack.NewDecoder(me.Bytes())
		v1, _ := md.GetInt()
		decodeTest(t, v1, i1)
	}

	// 16 bit
	{
		i1 := int64(259)
		i2 := int64(-259)
		me := msgpack.NewEncoder()
		me.PutInt(i1)
		me.PutInt(i2)
		md := msgpack.NewDecoder(me.Bytes())
		v1, _ := md.GetInt()
		decodeTest(t, v1, i1)
		v2, _ := md.GetInt()
		decodeTest(t, v2, i2)
	}

	// 32 bit
	{
		i1 := int64(65538)
		i2 := int64(-65538)
		me := msgpack.NewEncoder()
		me.PutInt(i1)
		me.PutInt(i2)
		md := msgpack.NewDecoder(me.Bytes())
		v1, _ := md.GetInt()
		decodeTest(t, v1, i1)
		v2, _ := md.GetInt()
		decodeTest(t, v2, i2)
	}

	// 64 bit
	{
		i1 := int64(4294967299)
		i2 := int64(-4294967299)
		me := msgpack.NewEncoder()
		me.PutInt(i1)
		me.PutInt(i2)
		md := msgpack.NewDecoder(me.Bytes())
		v1, _ := md.GetInt()
		decodeTest(t, v1, i1)
		v2, _ := md.GetInt()
		decodeTest(t, v2, i2)
	}
}

func TestDecodeMapLengths(t *testing.T) {
	// fixmap
	{
		l1 := uint32(10)
		me := msgpack.NewEncoder()
		me.PutMapLength(l1)
		md := msgpack.NewDecoder(me.Bytes())
		v1, _ := md.GetMapLength()
		decodeTest(t, v1, l1)
	}

	// 16 bit
	{
		l1 := uint32(30000)
		me := msgpack.NewEncoder()
		me.PutMapLength(l1)
		md := msgpack.NewDecoder(me.Bytes())
		v1, _ := md.GetMapLength()
		decodeTest(t, v1, l1)
	}

	// 32bit
	{
		l1 := uint32(80000)
		me := msgpack.NewEncoder()
		me.PutMapLength(l1)
		md := msgpack.NewDecoder(me.Bytes())
		v1, _ := md.GetMapLength()
		decodeTest(t, v1, l1)
	}
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
	// fixstr
	{
		s := strings.Repeat(base, 10)
		me := msgpack.NewEncoder()
		me.PutString(s)
		md := msgpack.NewDecoder(me.Bytes())
		v, _ := md.GetString()
		decodeTest(t, v, s)
	}

	// 8 bit length
	{
		s := strings.Repeat(base, 80)
		me := msgpack.NewEncoder()
		me.PutString(s)
		md := msgpack.NewDecoder(me.Bytes())
		v, _ := md.GetString()
		decodeTest(t, v, s)
	}

	// 16 bit length
	{
		s := strings.Repeat(base, 800)
		me := msgpack.NewEncoder()
		me.PutString(s)
		md := msgpack.NewDecoder(me.Bytes())
		v, _ := md.GetString()
		decodeTest(t, v, s)
	}

	// 32 bit length
	{
		s := strings.Repeat(base, 80000)
		me := msgpack.NewEncoder()
		me.PutString(s)
		md := msgpack.NewDecoder(me.Bytes())
		v, _ := md.GetString()
		decodeTest(t, v, s)
	}
}

func TestDecodeTime(t *testing.T) {
	loc, _ := time.LoadLocation("UTC")

	// timestamp32
	{
		d := time.Date(1997, 8, 28, 0, 0, 0, 0, loc)
		me := msgpack.NewEncoder()
		me.PutTime(d)
		md := msgpack.NewDecoder(me.Bytes())
		v, _ := md.GetTime()
		decodeTest(t, v, d)
	}

	// timestamp64
	{
		d := time.Date(1995, 9, 12, 0, 0, 0, 420, loc)
		me := msgpack.NewEncoder()
		me.PutTime(d)
		md := msgpack.NewDecoder(me.Bytes())
		v, _ := md.GetTime()
		decodeTest(t, v, d)
	}

	// timestamp96
	{
		d := time.Date(1961, 10, 19, 0, 0, 0, 420, loc)
		me := msgpack.NewEncoder()
		me.PutTime(d)
		md := msgpack.NewDecoder(me.Bytes())
		v, _ := md.GetTime()
		decodeTest(t, v, d)
	}
}

func TestDecodeUints(t *testing.T) {
	// fixint
	{
		i1 := uint64(69)
		me := msgpack.NewEncoder()
		me.PutUint(i1)
		md := msgpack.NewDecoder(me.Bytes())
		v1, _ := md.GetUint()
		decodeTest(t, v1, i1)
	}

	// 8 bit
	{
		i1 := uint64(130)
		me := msgpack.NewEncoder()
		me.PutUint(i1)
		md := msgpack.NewDecoder(me.Bytes())
		v1, _ := md.GetUint()
		decodeTest(t, v1, i1)
	}

	// 16 bit
	{
		i1 := uint64(259)
		me := msgpack.NewEncoder()
		me.PutUint(i1)
		md := msgpack.NewDecoder(me.Bytes())
		v1, _ := md.GetUint()
		decodeTest(t, v1, i1)
	}

	// 32 bit
	{
		i1 := uint64(65538)
		me := msgpack.NewEncoder()
		me.PutUint(i1)
		md := msgpack.NewDecoder(me.Bytes())
		v1, _ := md.GetUint()
		decodeTest(t, v1, i1)
	}

	// 64 bit
	{
		i1 := uint64(4294967299)
		me := msgpack.NewEncoder()
		me.PutUint(i1)
		md := msgpack.NewDecoder(me.Bytes())
		v1, _ := md.GetUint()
		decodeTest(t, v1, i1)
	}
}

func decodeTest(t *testing.T, v any, e any) {
	if v == e {
		return
	}
	es := fmt.Sprintf("%-8s %v", "expected", e)
	vs := fmt.Sprintf("%-8s %v", "actual", v)
	t.Fatalf("\n%s\n%s", es, vs)
}
