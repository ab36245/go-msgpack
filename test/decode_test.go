package test

import (
	"fmt"
	"testing"
	"time"

	"github.com/ab36245/go-msgpack"
)

func TestDecodeInts(t *testing.T) {
	// fixint
	{
		me := msgpack.NewEncoder()
		me.PutInt(69)
		me.PutInt(-11)
		md := msgpack.NewDecoder(me.Bytes())
		v, _ := md.GetInt()
		decodeTest(t, v, 69)
		v, _ = md.GetInt()
		decodeTest(t, v, -11)
	}
}

func TestDecodeNil(t *testing.T) {
	en := 69
	es := "follows nil"
	me := msgpack.NewEncoder()
	me.PutInt(en)
	me.PutNil()
	me.PutString(es)

	md := msgpack.NewDecoder(me.Bytes())

	vb, _ := md.IfNil()
	eb := false
	decodeTest(t, vb, eb)

	vn, _ := md.GetInt()
	decodeTest(t, vn, int64(en))

	vb, _ = md.IfNil()
	eb = true
	decodeTest(t, vb, eb)

	vs, _ := md.GetString()
	decodeTest(t, vs, es)
}

func TestDecodeString(t *testing.T) {
	e := "malaka was here"
	me := msgpack.NewEncoder()
	me.PutString(e)
	md := msgpack.NewDecoder(me.Bytes())

	v, _ := md.GetString()
	decodeTest(t, v, e)
}

func TestDecodeTime(t *testing.T) {
	e := "1997-08-28"
	me := msgpack.NewEncoder()
	d, _ := time.Parse(time.DateOnly, "1997-08-28")
	me.PutTime(d)
	md := msgpack.NewDecoder(me.Bytes())

	d, _ = md.GetTime()
	v := d.Format(time.DateOnly)
	decodeTest(t, v, e)
}

func decodeTest(t *testing.T, v any, e any) {
	if v == e {
		return
	}
	es := fmt.Sprintf("%-8s %v", "expected", e)
	vs := fmt.Sprintf("%-8s %v", "actual", v)
	t.Fatalf("\n%s\n%s", es, vs)
}
