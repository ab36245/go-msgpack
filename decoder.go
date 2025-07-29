package msgpack

import (
	"encoding/binary"
	"fmt"
	"math"
	"time"
)

func NewDecoder(bytes []byte) *Decoder {
	return &Decoder{bytes}
}

type Decoder struct {
	bytes []byte
}

func (d *Decoder) Bytes() []byte {
	return d.bytes
}

func (d *Decoder) IsEmpty() bool {
	return len(d.bytes) == 0
}

func (d *Decoder) Length() int {
	return len(d.bytes)
}

func (d *Decoder) GetArrayLength() (uint32, error) {
	b, err := d.readByte()
	if err != nil {
		return 0, err
	}
	if b&0xf0 == 0x90 {
		return uint32(b & 0x0f), nil
	}
	switch b {
	case 0xdc:
		n, err := d.readUint16()
		if err != nil {
			return 0, err
		}
		return uint32(n), nil
	case 0xdd:
		n, err := d.readUint32()
		if err != nil {
			return 0, err
		}
		return n, nil
	default:
		return invalid[uint32]("array length", b)
	}
}

func (d *Decoder) GetBinary() ([]byte, error) {
	b, err := d.readByte()
	if err != nil {
		return nil, err
	}
	var size int
	switch b {
	case 0xc4:
		n, err := d.readUint8()
		if err != nil {
			return nil, err
		}
		size = int(n)
	case 0xc5:
		n, err := d.readUint16()
		if err != nil {
			return nil, err
		}
		size = int(n)
	case 0xc6:
		n, err := d.readUint32()
		if err != nil {
			return nil, err
		}
		size = int(n)
	default:
		return invalid[[]byte]("binary", b)
	}
	return d.readBytes(size)
}

func (d *Decoder) GetBool() (bool, error) {
	b, err := d.readByte()
	if err != nil {
		return false, err
	}
	switch b {
	case 0xc2:
		return false, nil
	case 0xc3:
		return true, nil
	default:
		return invalid[bool]("bool", b)
	}
}

func (d *Decoder) GetExtUint() (byte, uint64, error) {
	b, err := d.readByte()
	if err != nil {
		return 0, 0, err
	}
	typ, err := d.readByte()
	if err != nil {
		return 0, 0, err
	}
	switch b {
	case 0xd4:
		n, err := d.readUint8()
		if err != nil {
			return 0, 0, err
		}
		return typ, uint64(n), nil
	case 0xd5:
		n, err := d.readUint16()
		if err != nil {
			return 0, 0, err
		}
		return typ, uint64(n), nil
	case 0xd6:
		n, err := d.readUint32()
		if err != nil {
			return 0, 0, err
		}
		return typ, uint64(n), nil
	case 0xd7:
		n, err := d.readUint64()
		if err != nil {
			return 0, 0, err
		}
		return typ, n, nil
	default:
		return invalid2[byte, uint64]("ext uint", b)
	}
}

func (d *Decoder) GetFloat() (float64, error) {
	b, err := d.readByte()
	if err != nil {
		return 0, err
	}
	switch b {
	case 0xca:
		n, err := d.readFloat32()
		if err != nil {
			return 0, err
		}
		return float64(n), nil
	case 0xcb:
		n, err := d.readFloat64()
		if err != nil {
			return 0, err
		}
		return n, nil
	default:
		return invalid[float64]("float", b)
	}
}

func (d *Decoder) GetInt() (int64, error) {
	b, err := d.readByte()
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
		n, err := d.readInt8()
		if err != nil {
			return 0, err
		}
		return int64(n), nil
	case 0xd1:
		n, err := d.readInt16()
		if err != nil {
			return 0, err
		}
		return int64(n), nil
	case 0xd2:
		n, err := d.readInt32()
		if err != nil {
			return 0, err
		}
		return int64(n), nil
	case 0xd3:
		n, err := d.readInt64()
		if err != nil {
			return 0, err
		}
		return n, nil
	default:
		return invalid[int64]("int", b)
	}
}

func (d *Decoder) GetMapLength() (uint32, error) {
	b, err := d.readByte()
	if err != nil {
		return 0, err
	}
	if b&0xf0 == 0x80 {
		return uint32(b & 0x0f), nil
	}
	switch b {
	case 0xde:
		n, err := d.readUint16()
		if err != nil {
			return 0, err
		}
		return uint32(n), nil
	case 0xdf:
		n, err := d.readUint32()
		if err != nil {
			return 0, err
		}
		return n, nil
	default:
		return invalid[uint32]("map length", b)
	}
}

func (d *Decoder) GetString() (string, error) {
	b, err := d.readByte()
	if err != nil {
		return "", err
	}
	var size int
	if b&0xe0 == 0xa0 {
		size = int(b & 0x1f)
	} else {
		switch b {
		case 0xd9:
			n, err := d.readUint8()
			if err != nil {
				return "", err
			}
			size = int(n)
		case 0xda:
			n, err := d.readUint16()
			if err != nil {
				return "", err
			}
			size = int(n)
		case 0xdb:
			n, err := d.readUint32()
			if err != nil {
				return "", err
			}
			size = int(n)
		default:
			return invalid[string]("string", b)
		}
	}
	bytes, err := d.readBytes(size)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func (d *Decoder) GetTime() (time.Time, error) {
	b, err := d.readByte()
	if err != nil {
		return time.Time{}, err
	}
	var sec int64
	var nsec int64
	switch b {
	case 0xd6:
		// timestamp 32
		t, err := d.readByte()
		if err != nil {
			return time.Time{}, err
		}
		if t != 255 {
			return invalid[time.Time]("timestamp32 type", t)
		}
		n, err := d.readInt32()
		if err != nil {
			return time.Time{}, err
		}
		sec = int64(n)
		nsec = 0

	case 0xd7:
		// timestamp 64
		t, err := d.readByte()
		if err != nil {
			return time.Time{}, err
		}
		if t != 255 {
			return invalid[time.Time]("timestamp64 type", t)
		}
		n, err := d.readUint64()
		if err != nil {
			return time.Time{}, err
		}
		sec = int64(n & 0x7fffffff)
		nsec = int64(n >> 34)

	case 0xc7:
		// timestamp 96
		n, err := d.readByte()
		if err != nil {
			return time.Time{}, err
		}
		if n != 12 {
			return invalid[time.Time]("timestamp95 length", n)
		}
		t, err := d.readByte()
		if err != nil {
			return time.Time{}, err
		}
		if t != 255 {
			return invalid[time.Time]("timestamp95 type", t)
		}
		n1, err := d.readUint32()
		if err != nil {
			return time.Time{}, err
		}
		n2, err := d.readInt64()
		if err != nil {
			return time.Time{}, err
		}
		sec = n2
		nsec = int64(n1)

	default:
		return invalid[time.Time]("timestamp extension", b)
	}
	return time.Unix(sec, nsec).UTC(), nil
}

func (d *Decoder) GetUint() (uint64, error) {
	b, err := d.readByte()
	if err != nil {
		return 0, err
	}
	if b&0x80 == 0 {
		return uint64(b), nil
	}
	switch b {
	case 0xcc:
		n, err := d.readUint8()
		if err != nil {
			return 0, err
		}
		return uint64(n), nil
	case 0xcd:
		n, err := d.readUint16()
		if err != nil {
			return 0, err
		}
		return uint64(n), nil
	case 0xce:
		n, err := d.readUint32()
		if err != nil {
			return 0, err
		}
		return uint64(n), nil
	case 0xcf:
		n, err := d.readUint64()
		if err != nil {
			return 0, err
		}
		return n, nil
	default:
		return invalid[uint64]("uint", b)
	}
}

func (d *Decoder) IfNil() (bool, error) {
	isNil, err := d.IsNil()
	if err != nil {
		return false, err
	}
	if isNil {
		d.readByte()
	}
	return isNil, err
}

func (d *Decoder) IsNil() (bool, error) {
	b, err := d.peekByte()
	if err != nil {
		return false, err
	}
	return b == 0xc0, nil
}

func (d *Decoder) peekByte() (byte, error) {
	return peek(d.bytes, 1, func(bytes []byte) byte {
		return bytes[0]
	})
}

func (d *Decoder) readByte() (byte, error) {
	return read(&d.bytes, 1, func(bytes []byte) byte {
		return bytes[0]
	})
}

func (d *Decoder) readBytes(n int) ([]byte, error) {
	return read(&d.bytes, n, func(bytes []byte) []byte {
		return bytes[0:n]
	})
}

func (d *Decoder) readFloat32() (float32, error) {
	return read(&d.bytes, 4, func(bytes []byte) float32 {
		u := binary.BigEndian.Uint32(bytes)
		return math.Float32frombits(u)
	})
}

func (d *Decoder) readFloat64() (float64, error) {
	return read(&d.bytes, 8, func(bytes []byte) float64 {
		u := binary.BigEndian.Uint64(bytes)
		return math.Float64frombits(u)
	})
}

func (d *Decoder) readInt8() (int8, error) {
	return read(&d.bytes, 1, func(bytes []byte) int8 {
		return int8(bytes[0])
	})
}

func (d *Decoder) readInt16() (int16, error) {
	return read(&d.bytes, 2, func(bytes []byte) int16 {
		return int16(binary.BigEndian.Uint16(bytes))
	})
}

func (d *Decoder) readInt32() (int32, error) {
	return read(&d.bytes, 4, func(bytes []byte) int32 {
		return int32(binary.BigEndian.Uint32(bytes))
	})
}

func (d *Decoder) readInt64() (int64, error) {
	return read(&d.bytes, 8, func(bytes []byte) int64 {
		return int64(binary.BigEndian.Uint64(bytes))
	})
}

func (d *Decoder) readUint8() (uint8, error) {
	return read(&d.bytes, 1, func(bytes []byte) uint8 {
		return bytes[0]
	})
}

func (d *Decoder) readUint16() (uint16, error) {
	return read(&d.bytes, 2, func(bytes []byte) uint16 {
		return binary.BigEndian.Uint16(bytes)
	})
}

func (d *Decoder) readUint32() (uint32, error) {
	return read(&d.bytes, 4, func(bytes []byte) uint32 {
		return binary.BigEndian.Uint32(bytes)
	})
}

func (d *Decoder) readUint64() (uint64, error) {
	return read(&d.bytes, 8, func(bytes []byte) uint64 {
		return binary.BigEndian.Uint64(bytes)
	})
}

func peek[T any](bytes []byte, size int, f func([]byte) T) (T, error) {
	excess := size - len(bytes)
	if excess > 0 {
		return fail[T]("trying to read %d bytes beyond end of buffer (%d bytes)", excess, len(bytes))
	}
	return f(bytes), nil
}

func read[T any](bytes *[]byte, size int, f func([]byte) T) (T, error) {
	value, err := peek(*bytes, size, f)
	if err == nil {
		*bytes = (*bytes)[size:]
	}
	return value, err
}

func fail[T any](mesg string, args ...any) (T, error) {
	return *new(T), fmt.Errorf(mesg, args...)
}

func invalid[T any](what string, b byte) (T, error) {
	return fail[T]("invalid byte for %s (%#02x)", what, b)
}

func invalid2[T, U any](what string, b byte) (T, U, error) {
	return *new(T), *new(U), fmt.Errorf("invalid byte for %s (%#02x)", what, b)
}
