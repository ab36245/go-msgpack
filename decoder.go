package msgpack

import (
	"bytes"
	"errors"
	"time"

	"github.com/vmihailenco/msgpack/v5"
)

func NewDecoder(data []byte) *Decoder {
	d := &Decoder{
		mp: msgpack.NewDecoder(bytes.NewReader(data)),
	}
	return d
}

type Decoder struct {
	errs []error
	mp   *msgpack.Decoder
}

func (d *Decoder) Err() error {
	if len(d.errs) == 0 {
		return nil
	}
	return DecodeError.Wrap(errors.Join(d.errs...))
}

func (d *Decoder) Errs() []error {
	return d.errs
}

func (d *Decoder) GetArrayLength() int {
	value, err := d.mp.DecodeArrayLen()
	d.error(err)
	return value
}

func (d *Decoder) GetBool() bool {
	value, err := d.mp.DecodeBool()
	d.error(err)
	return value
}

func (d *Decoder) GetBytes() []byte {
	value, err := d.mp.DecodeBytes()
	d.error(err)
	return value
}

func (d *Decoder) GetFloat() float64 {
	value, err := d.mp.DecodeFloat64()
	d.error(err)
	return value
}

func (d *Decoder) GetInt() int64 {
	value, err := d.mp.DecodeInt64()
	d.error(err)
	return value
}

func (d *Decoder) GetMapLength() int {
	value, err := d.mp.DecodeMapLen()
	d.error(err)
	return value
}

func (d *Decoder) GetString() string {
	value, err := d.mp.DecodeString()
	d.error(err)
	return value
}

func (d *Decoder) GetTime() time.Time {
	value, err := d.mp.DecodeTime()
	d.error(err)
	return value
}

func (d *Decoder) GetUint() uint64 {
	value, err := d.mp.DecodeUint64()
	d.error(err)
	return value
}

func (d *Decoder) error(err error) {
	if err != nil {
		d.errs = append(d.errs, err)
	}
}
