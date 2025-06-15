package msgpack

import (
	"time"

	"github.com/vmihailenco/msgpack/v5"
)

func NewEncoder() *Encoder {
	e := &Encoder{}
	e.mp = msgpack.NewEncoder(e)
	return e
}

type Encoder struct {
	buf []byte
	mp  *msgpack.Encoder
}

func (e *Encoder) Bytes() []byte {
	return e.buf
}

func (e *Encoder) PutArrayLength(value int) error {
	return e.mp.EncodeArrayLen(value)
}

func (e *Encoder) PutBool(value bool) error {
	return e.mp.EncodeBool(value)
}

func (e *Encoder) PutBytes(value []byte) error {
	return e.mp.EncodeBytes(value)
}

func (e *Encoder) PutFloat(value float64) error {
	return e.mp.EncodeFloat64(value)
}

func (e *Encoder) PutInt(value int64) error {
	return e.mp.EncodeInt(value)
}

func (e *Encoder) PutMapLength(value int) error {
	return e.mp.EncodeMapLen(value)
}

func (e *Encoder) PutString(value string) error {
	return e.mp.EncodeString(value)
}

func (e *Encoder) PutTime(value time.Time) error {
	return e.mp.EncodeTime(value)
}

func (e *Encoder) PutUint(value uint64) error {
	return e.mp.EncodeUint(value)
}

func (e *Encoder) Write(p []byte) (int, error) {
	e.buf = append(e.buf, p...)
	return len(p), nil
}
