package msgpack

import (
	"fmt"
	"time"
)

func NewDecoder(bytes []byte) *Decoder {
	return &Decoder{
		reader: &reader{
			bytes: bytes,
		},
	}
}

type Decoder struct {
	reader *reader
}

func (d *Decoder) GetArrayLength() (uint64, error) {
	b, err := d.reader.readByte()
	if err != nil {
		return 0, err
	}
	if b&0xf0 == 0x90 {
		return uint64(b & 0x0f), nil
	}
	switch b {
	case 0xdc:
		return d.reader.readUint16()
	case 0xdd:
		return d.reader.readUint32()
	default:
		return badByte[uint64](b, "array length")
	}
}

func (d *Decoder) GetBool() (bool, error) {
	b, err := d.reader.readByte()
	if err != nil {
		return false, err
	}
	switch b {
	case 0xc2:
		return false, nil
	case 0xc3:
		return true, nil
	default:
		return badByte[bool](b, "bool")
	}
}

func (d *Decoder) GetBytes() ([]byte, error) {
	b, err := d.reader.readByte()
	if err != nil {
		return nil, err
	}
	var n uint64
	switch b {
	case 0xc4:
		n, err = d.reader.readUint8()
	case 0xc5:
		n, err = d.reader.readUint16()
	case 0xc6:
		n, err = d.reader.readUint32()
	default:
		return badByte[[]byte](b, "bytes")
	}
	if err != nil {
		return nil, err
	}
	return d.reader.readBytes(int(n))
}

func (d *Decoder) GetFloat() (float64, error) {
	b, err := d.reader.readByte()
	if err != nil {
		return 0, err
	}
	switch b {
	case 0xca:
		return d.reader.readFloat32()
	case 0xcb:
		return d.reader.readFloat64()
	default:
		return badByte[float64](b, "float")
	}
}

func (d *Decoder) GetInt() (int64, error) {
	b, err := d.reader.readByte()
	if err != nil {
		return 0, err
	}
	if b&0x80 == 0 {
		// positive
		return int64(b), nil
	}
	if b&0xe0 == 0xe0 {
		// negative
		return int64(b) - 256, nil
	}
	switch b {
	case 0xd0:
		return d.reader.readInt8()
	case 0xd1:
		return d.reader.readInt16()
	case 0xd2:
		return d.reader.readInt32()
	case 0xd3:
		return d.reader.readInt64()
	default:
		return badByte[int64](b, "int")
	}
}

func (d *Decoder) GetMapLength() (uint64, error) {
	b, err := d.reader.readByte()
	if err != nil {
		return 0, err
	}
	if b&0xf0 == 0x80 {
		return uint64(b & 0x0f), nil
	}
	switch b {
	case 0xde:
		return d.reader.readUint16()
	case 0xdf:
		return d.reader.readUint32()
	default:
		return badByte[uint64](b, "map length")
	}
}

func (d *Decoder) GetString() (string, error) {
	b, err := d.reader.readByte()
	if err != nil {
		return "", err
	}
	var n uint64
	if b&0xe0 == 0xa0 {
		n = uint64(b & 0x1f)
	} else {
		switch b {
		case 0xd9:
			n, err = d.reader.readUint8()
		case 0xda:
			n, err = d.reader.readUint16()
		case 0xdb:
			n, err = d.reader.readUint32()
		default:
			return badByte[string](b, "string")
		}
	}
	if err != nil {
		return "", nil
	}
	u, err := d.reader.readBytes(int(n))
	if err != nil {
		return "", err
	}
	return string(u), nil
}

func (d *Decoder) GetTime() (time.Time, error) {
	fail := func(err error) (time.Time, error) {
		return time.Time{}, err
	}

	mesg := func(mesg string, args ...any) (time.Time, error) {
		return fail(fmt.Errorf(mesg, args...))
	}

	b, err := d.reader.readByte()
	if err != nil {
		return fail(err)
	}
	nsec := uint64(0)
	sec := int64(0)
	switch b {
	case 0xd6:
		// timestamp 32
		t, err := d.reader.readByte()
		if err != nil {
			return fail(err)
		}
		if t != 255 {
			return mesg("invalid type for timestamp 32 extension (%#02x)", t)
		}
		nsec = 0
		sec, err = d.reader.readInt32()
		if err != nil {
			return fail(err)
		}

	case 0xd7:
		// timestamp 64
		t, err := d.reader.readByte()
		if err != nil {
			return fail(err)
		}
		if t != 255 {
			return mesg("invalid type for timestamp 64 extension (%#02x)", t)
		}
		data64, err := d.reader.readUint64()
		if err != nil {
			return fail(err)
		}
		nsec = data64 >> 34
		sec = int64(data64 & 0x7fffffff)

	case 0xc7:
		// timestamp 96
		n, err := d.reader.readByte()
		if err != nil {
			return fail(err)
		}
		if n != 12 {
			return mesg("invalid length for timestamp 96 extension (%#02x)", n)
		}
		t, err := d.reader.readByte()
		if err != nil {
			return fail(err)
		}
		if t != 255 {
			return mesg("invalid type for timestamp 96 extension (%#02x)", t)
		}
		nsec, err = d.reader.readUint32()
		if err != nil {
			return fail(err)
		}
		sec, err = d.reader.readInt64()
		if err != nil {
			return fail(err)
		}

	default:
		badByte[time.Time](b, "timestamp extension")
	}
	usec := sec*1000*1000 + int64(nsec%1000)
	return time.UnixMicro(usec), nil
}

func (d *Decoder) GetUint() (uint64, error) {
	b, err := d.reader.readByte()
	if err != nil {
		return 0, err
	}
	if b&0x80 == 0 {
		return uint64(b), nil
	}
	switch b {
	case 0xcc:
		return d.reader.readUint8()
	case 0xcd:
		return d.reader.readUint16()
	case 0xce:
		return d.reader.readUint32()
	case 0xcf:
		return d.reader.readUint64()
	default:
		return badByte[uint64](b, "uint")
	}
}

func (d *Decoder) IfNil() (bool, error) {
	isNil, err := d.IsNil()
	if err != nil {
		return false, err
	}
	if isNil {
		d.reader.readByte()
	}
	return isNil, err
}

func (d *Decoder) IsNil() (bool, error) {
	b, err := d.reader.peekByte()
	if err != nil {
		return false, err
	}
	return b == 0xc0, nil
}

func badByte[T any](b byte, what string) (T, error) {
	return *new(T), fmt.Errorf("invalid byte for %s (%#02x)", what, b)
}
