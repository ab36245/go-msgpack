package msgpack

import (
	"encoding/binary"
	"fmt"
	"math"
)

type reader struct {
	bytes []byte
}

func (r *reader) peekByte() (byte, error) {
	return peek(r.bytes, 1, func(bytes []byte) byte {
		return bytes[0]
	})
}

func (r *reader) readByte() (byte, error) {
	return read(&r.bytes, 1, func(bytes []byte) byte {
		return bytes[0]
	})
}

func (r *reader) readBytes(n int) ([]byte, error) {
	return read(&r.bytes, n, func(bytes []byte) []byte {
		return bytes
	})
}

func (r *reader) readFloat32() (float64, error) {
	return read(&r.bytes, 4, func(bytes []byte) float64 {
		u := binary.BigEndian.Uint32(bytes)
		return float64(math.Float32frombits(u))
	})
}

func (r *reader) readFloat64() (float64, error) {
	return read(&r.bytes, 8, func(bytes []byte) float64 {
		u := binary.BigEndian.Uint64(bytes)
		return math.Float64frombits(u)
	})
}

func (r *reader) readInt8() (int64, error) {
	return read(&r.bytes, 1, func(bytes []byte) int64 {
		return int64(bytes[0])
	})
}

func (r *reader) readInt16() (int64, error) {
	return read(&r.bytes, 2, func(bytes []byte) int64 {
		return int64(binary.BigEndian.Uint16(bytes))
	})
}

func (r *reader) readInt32() (int64, error) {
	return read(&r.bytes, 4, func(bytes []byte) int64 {
		return int64(binary.BigEndian.Uint32(bytes))
	})
}

func (r *reader) readInt64() (int64, error) {
	return read(&r.bytes, 8, func(bytes []byte) int64 {
		return int64(binary.BigEndian.Uint64(bytes))
	})
}

func (r *reader) readUint8() (uint64, error) {
	return read(&r.bytes, 1, func(bytes []byte) uint64 {
		return uint64(bytes[0])
	})
}

func (r *reader) readUint16() (uint64, error) {
	return read(&r.bytes, 2, func(bytes []byte) uint64 {
		return uint64(binary.BigEndian.Uint16(bytes))
	})
}

func (r *reader) readUint32() (uint64, error) {
	return read(&r.bytes, 4, func(bytes []byte) uint64 {
		return uint64(binary.BigEndian.Uint32(bytes))
	})
}

func (r *reader) readUint64() (uint64, error) {
	return read(&r.bytes, 8, func(bytes []byte) uint64 {
		return binary.BigEndian.Uint64(bytes)
	})
}

func peek[T any](bytes []byte, size int, f func([]byte) T) (T, error) {
	excess := size - len(bytes)
	if excess > 0 {
		return *new(T), fmt.Errorf("trying to read %d bytes beyond end of buffer (%d bytes)", excess, len(bytes))
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
