package msgpack

import (
	"bytes"
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
	mp *msgpack.Decoder
}

func (d *Decoder) GetArrayLength() (int, error) {
	return d.mp.DecodeArrayLen()
}

func (d *Decoder) GetBool() (bool, error) {
	return d.mp.DecodeBool()
}

func (d *Decoder) GetBytes() ([]byte, error) {
	return d.mp.DecodeBytes()
}

func (d *Decoder) GetFloat() (float64, error) {
	return d.mp.DecodeFloat64()
}

func (d *Decoder) GetInt() (int64, error) {
	return d.mp.DecodeInt64()
}

func (d *Decoder) GetMapLength() (int, error) {
	return d.mp.DecodeMapLen()
}

func (d *Decoder) GetString() (string, error) {
	return d.mp.DecodeString()
}

func (d *Decoder) GetTime() (time.Time, error) {
	return d.mp.DecodeTime()
}

func (d *Decoder) GetUint() (uint64, error) {
	return d.mp.DecodeUint64()
}
