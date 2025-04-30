package msgpack

import (
	"errors"
	"fmt"
	"time"

	"github.com/vmihailenco/msgpack/v5"
)

func NewEncoder() *Encoder {
	e := &Encoder{}
	e.mp = msgpack.NewEncoder(e)
	return e
}

type Encoder struct {
	buf  []byte
	errs []error
	mp   *msgpack.Encoder
}

func (e *Encoder) Bytes() []byte {
	return e.buf
}

func (e *Encoder) Err() error {
	if len(e.errs) == 0 {
		return nil
	}
	return fmt.Errorf("msgpack encoding: %w", errors.Join(e.errs...))
}

func (e *Encoder) Errs() []error {
	return e.errs
}

func (e *Encoder) PutArrayLength(value int) {
	e.error(e.mp.EncodeArrayLen(value))
}

func (e *Encoder) PutBool(value bool) {
	e.error(e.mp.EncodeBool(value))
}

func (e *Encoder) PutBytes(value []byte) {
	e.error(e.mp.EncodeBytes(value))
}

func (e *Encoder) PutFloat(value float64) {
	e.error(e.mp.EncodeFloat64(value))
}

func (e *Encoder) PutInt(value int64) {
	e.error(e.mp.EncodeInt(value))
}

func (e *Encoder) PutMapLength(value int) {
	e.error(e.mp.EncodeMapLen(value))
}

func (e *Encoder) PutString(value string) {
	e.error(e.mp.EncodeString(value))
}

func (e *Encoder) PutTime(value time.Time) {
	e.error(e.mp.EncodeTime(value))
}

func (e *Encoder) PutUint(value uint64) {
	e.error(e.mp.EncodeUint(value))
}

func (e *Encoder) Result() ([]byte, error) {
	if len(e.errs) > 0 {
		return nil, e.Err()
	}
	return e.Bytes(), nil
}

func (e *Encoder) Write(p []byte) (int, error) {
	e.buf = append(e.buf, p...)
	return len(p), nil
}

func (e *Encoder) error(err error) {
	if err != nil {
		e.errs = append(e.errs, err)
	}
}
