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

func (d *Decoder) GetArrayLength() (uint32, error) {
	b, err := d.reader.readByte()
	if err != nil {
		return 0, err
	}
	if b&0xf0 == 0x90 {
		return uint32(b & 0x0f), nil
	}
	switch b {
	case 0xdc:
		n, err := d.reader.readUint16()
		if err != nil {
			return 0, err
		}
		return uint32(n), nil
	case 0xdd:
		n, err := d.reader.readUint32()
		if err != nil {
			return 0, err
		}
		return n, nil
	default:
		return 0, invalid("array length", b)
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
		return false, invalid("bool", b)
	}
}

func (d *Decoder) GetBytes() ([]byte, error) {
	b, err := d.reader.readByte()
	if err != nil {
		return nil, err
	}
	var size int
	switch b {
	case 0xc4:
		n, err := d.reader.readUint8()
		if err != nil {
			return nil, err
		}
		size = int(n)
	case 0xc5:
		n, err := d.reader.readUint16()
		if err != nil {
			return nil, err
		}
		size = int(n)
	case 0xc6:
		n, err := d.reader.readUint32()
		if err != nil {
			return nil, err
		}
		size = int(n)
	default:
		return nil, invalid("byte slice", b)
	}
	if err != nil {
		return nil, err
	}
	return d.reader.readBytes(size)
}

func (d *Decoder) GetFloat() (float64, error) {
	b, err := d.reader.readByte()
	if err != nil {
		return 0, err
	}
	switch b {
	case 0xca:
		n, err := d.reader.readFloat32()
		if err != nil {
			return 0, err
		}
		return float64(n), nil
	case 0xcb:
		n, err := d.reader.readFloat64()
		if err != nil {
			return 0, err
		}
		return n, nil
	default:
		return 0, invalid("float", b)
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
		n, err := d.reader.readInt8()
		if err != nil {
			return 0, err
		}
		return int64(n), nil
	case 0xd1:
		n, err := d.reader.readInt16()
		if err != nil {
			return 0, err
		}
		return int64(n), nil
	case 0xd2:
		n, err := d.reader.readInt32()
		if err != nil {
			return 0, err
		}
		return int64(n), nil
	case 0xd3:
		n, err := d.reader.readInt64()
		if err != nil {
			return 0, err
		}
		return n, nil
	default:
		return 0, invalid("int", b)
	}
}

func (d *Decoder) GetMapLength() (uint32, error) {
	b, err := d.reader.readByte()
	if err != nil {
		return 0, err
	}
	if b&0xf0 == 0x80 {
		return uint32(b & 0x0f), nil
	}
	switch b {
	case 0xde:
		n, err := d.reader.readUint16()
		if err != nil {
			return 0, err
		}
		return uint32(n), nil
	case 0xdf:
		n, err := d.reader.readUint32()
		if err != nil {
			return 0, err
		}
		return n, nil
	default:
		return 0, invalid("map length", b)
	}
}

func (d *Decoder) GetString() (string, error) {
	b, err := d.reader.readByte()
	if err != nil {
		return "", err
	}
	var size int
	if b&0xe0 == 0xa0 {
		size = int(b & 0x1f)
	} else {
		switch b {
		case 0xd9:
			n, err := d.reader.readUint8()
			if err != nil {
				return "", err
			}
			size = int(n)
		case 0xda:
			n, err := d.reader.readUint16()
			if err != nil {
				return "", err
			}
			size = int(n)
		case 0xdb:
			n, err := d.reader.readUint32()
			if err != nil {
				return "", err
			}
			size = int(n)
		default:
			return "", invalid("string", b)
		}
	}
	if err != nil {
		return "", err
	}
	bytes, err := d.reader.readBytes(size)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func (d *Decoder) GetTime() (time.Time, error) {
	b, err := d.reader.readByte()
	if err != nil {
		return time.Time{}, err
	}
	var sec int64
	var nsec int64
	switch b {
	case 0xd6:
		// timestamp 32
		t, err := d.reader.readByte()
		if err != nil {
			return time.Time{}, err
		}
		if t != 255 {
			err := fmt.Errorf("invalid type for timestamp 32 extension (%#02x)", t)
			return time.Time{}, err
		}
		n, err := d.reader.readInt32()
		if err != nil {
			return time.Time{}, err
		}
		sec = int64(n)
		nsec = 0

	case 0xd7:
		// timestamp 64
		t, err := d.reader.readByte()
		if err != nil {
			return time.Time{}, err
		}
		if t != 255 {
			err := fmt.Errorf("invalid type for timestamp 64 extension (%#02x)", t)
			return time.Time{}, err
		}
		n, err := d.reader.readUint64()
		if err != nil {
			return time.Time{}, err
		}
		sec = int64(n & 0x7fffffff)
		nsec = int64(n >> 34)

	case 0xc7:
		// timestamp 96
		n, err := d.reader.readByte()
		if err != nil {
			return time.Time{}, err
		}
		if n != 12 {
			err := fmt.Errorf("invalid length for timestamp 96 extension (%#02x)", n)
			return time.Time{}, err
		}
		t, err := d.reader.readByte()
		if err != nil {
			return time.Time{}, err
		}
		if t != 255 {
			err := fmt.Errorf("invalid type for timestamp 96 extension (%#02x)", t)
			return time.Time{}, err
		}
		n1, err := d.reader.readUint32()
		if err != nil {
			return time.Time{}, err
		}
		n2, err := d.reader.readInt64()
		if err != nil {
			return time.Time{}, err
		}
		sec = n2
		nsec = int64(n1)

	default:
		return time.Time{}, invalid("timestamp extension", b)
	}
	return time.Unix(sec, nsec).UTC(), nil
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
		n, err := d.reader.readUint8()
		if err != nil {
			return 0, err
		}
		return uint64(n), nil
	case 0xcd:
		n, err := d.reader.readUint16()
		if err != nil {
			return 0, err
		}
		return uint64(n), nil
	case 0xce:
		n, err := d.reader.readUint32()
		if err != nil {
			return 0, err
		}
		return uint64(n), nil
	case 0xcf:
		n, err := d.reader.readUint64()
		if err != nil {
			return 0, err
		}
		return n, nil
	default:
		return 0, invalid("uint", b)
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

func invalid(what string, b byte) error {
	return fmt.Errorf("invalid byte for %s (%#02x)", what, b)
}
