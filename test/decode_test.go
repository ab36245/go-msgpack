package test

import (
	"fmt"
	"testing"

	"github.com/ab36245/go-msgpack"
)

func TestDecodeArrayLengths(t *testing.T) {
	t.Fatal("not implemented")
}

func TestDecodeBools(t *testing.T) {
	t.Fatal("not implemented")
}

func TestDecodeBytes(t *testing.T) {
	t.Fatal("not implemented")
}

func TestDecodeFloats(t *testing.T) {
	t.Fatal("not implemented")
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
	t.Fatal("not implemented")
}

func TestDecodeNil(t *testing.T) {
	t.Fatal("not implemented")
}

func TestDecodeStrings(t *testing.T) {
	t.Fatal("not implemented")
}

func TestDecodeTime(t *testing.T) {
	t.Fatal("not implemented")
}

func TestDecodeUints(t *testing.T) {
	t.Fatal("not implemented")
}

func decodeTest(t *testing.T, v any, e any) {
	if v == e {
		return
	}
	es := fmt.Sprintf("%-8s %v", "expected", e)
	vs := fmt.Sprintf("%-8s %v", "actual", v)
	t.Fatalf("\n%s\n%s", es, vs)
}
