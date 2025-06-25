package test

import (
	"fmt"
	"testing"
	"time"

	"github.com/ab36245/go-msgpack"
)

func TestDecodeInts(t *testing.T) {
	me := msgpack.NewEncoder()
	me.PutInt(0x45)
	me.PutInt(0x1A4)

	md := msgpack.NewDecoder(me.Bytes())

	v1, _ := md.GetInt()
	e1 := int64(69)
	decodeTest(t, v1, e1)

	v2, _ := md.GetInt()
	e2 := int64(420)
	decodeTest(t, v2, e2)
}

func TestDecodeNil(t *testing.T) {
	en := int64(69)
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
	decodeTest(t, vn, en)

	vb, _ = md.IfNil()
	eb = true
	decodeTest(t, vb, eb)

	vs, _ := md.GetString()
	decodeTest(t, vs, es)
}

func TestDecodeString(t *testing.T) {
	mp := msgpack.NewEncoder()
	mp.PutString("malaka")
	v := mp.Bytes()
	e := `
		|7 bytes
    	|    0000 a6 6d 61 6c 61 6b 61
	`
	encodeTest(t, v, e)
}

func TestDecodeTime(t *testing.T) {
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
