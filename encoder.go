package msgpack

import (
	"encoding/binary"
	"fmt"
	"math"
	"time"
)

func NewEncoder() *Encoder {
	return &Encoder{}
}

type Encoder struct {
	bytes []byte
}

func (e *Encoder) Bytes() []byte {
	return e.bytes
}

func (e *Encoder) AsString(maxLength int) string {
	s := ""
	for i, b := range e.bytes {
		if maxLength >= 0 && i >= maxLength {
			break
		}
		if s != "" {
			s += " "
		}
		s += fmt.Sprintf("%02x", b)
	}
	return s
}

func (e *Encoder) Clear() {
	e.bytes = nil
}

func (e *Encoder) PutArrayLength(v uint32) {
	if v <= mask4 {
		e.writeByte(byte(0x90 | v))
	} else if v <= mask16 {
		e.writeByte(0xdc)
		e.writeUint16(uint16(v))
	} else {
		e.writeByte(0xdd)
		e.writeUint32(uint32(v))
	}
}

func (e *Encoder) PutBinary(v []byte) error {
	n := len(v)
	if n <= mask8 {
		e.writeByte(0xc4)
		e.writeUint8(uint8(n))
	} else if n <= mask16 {
		e.writeByte(0xc5)
		e.writeUint16(uint16(n))
	} else if n <= mask32 {
		e.writeByte(0xc6)
		e.writeUint32(uint32(n))
	} else {
		return fmt.Errorf("byte slice (%d bytes) is too long to encode", n)
	}
	e.writeBytes(v)
	return nil
}

func (e *Encoder) PutBool(v bool) {
	if !v {
		e.writeByte(0xc2)
	} else {
		e.writeByte(0xc3)
	}
}

func (e *Encoder) PutBytes(v []byte) {
	e.writeBytes(v)
}

func (e *Encoder) PutExtUint(typ uint8, v uint64) {
	if v <= mask8 {
		e.writeByte(0xd4)
		e.writeByte(typ)
		e.writeUint8(uint8(v))
	} else if v <= mask16 {
		e.writeByte(0xd5)
		e.writeByte(typ)
		e.writeUint16(uint16(v))
	} else if v <= mask32 {
		e.writeByte(0xd6)
		e.writeByte(typ)
		e.writeUint32(uint32(v))
	} else {
		e.writeByte(0xd7)
		e.writeByte(typ)
		e.writeUint64(v)
	}
}

func (e *Encoder) PutFloat(v float64) {
	// Try to work out if encoding a single-precision IEEE754 value
	// is acceptable
	ieee754 := math.Float64bits(v)

	// ...check if the exponent is inside the range for single-precision
	biasedExp := (ieee754 >> 52) & 0x7ff
	unbiasedExp := int(biasedExp - 1023)
	if unbiasedExp < -128 || unbiasedExp > 127 {
		// nope!
		e.PutFloat64(v)
		return
	}

	// ...check if the least significant bits of the mantissa are zero
	leastSignficantBits := ieee754 & ((1 << 29) - 1)
	if leastSignficantBits != 0 {
		// nope!
		e.PutFloat64(v)
		return
	}

	// ...it looks ok to only encode a single-precision (32 bit) version
	e.PutFloat32(float32(v))
}

func (e *Encoder) PutFloat32(v float32) {
	e.writeByte(0xca)
	e.writeFloat32(v)
}

func (e *Encoder) PutFloat64(v float64) {
	e.writeByte(0xcb)
	e.writeFloat64(v)
}

func (e *Encoder) PutInt(v int64) {
	if v >= intFixMin && v <= intFixMax {
		e.writeInt8(int8(v))
	} else if v >= int8Min && v <= int8Max {
		e.writeByte(0xd0)
		e.writeInt8(int8(v))
	} else if v >= int16Min && v <= int16Max {
		e.writeByte(0xd1)
		e.writeInt16(int16(v))
	} else if v >= int32Min && v <= int32Max {
		e.writeByte(0xd2)
		e.writeInt32(int32(v))
	} else {
		e.writeByte(0xd3)
		e.writeInt64(v)
	}
}

func (e *Encoder) PutMapLength(v uint32) {
	if v <= mask4 {
		e.writeByte(byte(0x80 | v))
	} else if v <= mask16 {
		e.writeByte(0xde)
		e.writeUint16(uint16(v))
	} else {
		e.writeByte(0xdf)
		e.writeUint32(uint32(v))
	}
}

func (e *Encoder) PutNil() {
	e.writeByte(0xc0)
}

func (e *Encoder) PutString(v string) error {
	u := []byte(v)
	n := len(u)
	if n <= mask5 {
		e.writeByte(byte(0xa0 | n))
	} else if n <= mask8 {
		e.writeByte(0xd9)
		e.writeUint8(uint8(n))
	} else if n <= mask16 {
		e.writeByte(0xda)
		e.writeUint16(uint16(n))
	} else if n <= mask32 {
		e.writeByte(0xdb)
		e.writeUint32(uint32(n))
	} else {
		return fmt.Errorf("string (%d bytes) is too long to encode", n)
	}
	e.writeBytes(u)
	return nil
}

func (e *Encoder) PutTime(v time.Time) {
	sec := v.Unix()
	nsec := v.Nanosecond()
	if sec < 0 || sec > mask34 {
		//timestamp 96
		e.writeByte(0xc7)
		e.writeByte(12)
		e.writeByte(0xff)
		e.writeUint32(uint32(nsec))
		e.writeInt64(sec)
	} else if sec > mask32 || nsec > 0 {
		//timestamp 64
		e.writeByte(0xd7)
		e.writeByte(0xff)
		data64 := uint64(nsec)<<34 | uint64(sec)
		e.writeUint64(data64)
	} else {
		// timestamp 32
		e.writeByte(0xd6)
		e.writeByte(0xff)
		e.writeInt32(int32(sec))
	}
}

func (e *Encoder) PutUint(v uint64) {
	if v <= mask7 {
		e.writeByte(byte(v))
	} else if v <= mask8 {
		e.writeByte(0xcc)
		e.writeUint8(uint8(v))
	} else if v <= mask16 {
		e.writeByte(0xcd)
		e.writeUint16(uint16(v))
	} else if v <= mask32 {
		e.writeByte(0xce)
		e.writeUint32(uint32(v))
	} else {
		e.writeByte(0xcf)
		e.writeUint64(v)
	}
}

func (e *Encoder) writeByte(v byte) {
	e.bytes = append(e.bytes, v)
}

func (e *Encoder) writeBytes(v []byte) {
	e.bytes = append(e.bytes, v...)
}

func (e *Encoder) writeFloat32(v float32) {
	e.put32(math.Float32bits(v))
}

func (e *Encoder) writeFloat64(v float64) {
	e.put64(math.Float64bits(v))
}

func (e *Encoder) writeInt8(v int8) {
	e.put8(uint8(v))
}

func (e *Encoder) writeInt16(v int16) {
	e.put16(uint16(v))
}

func (e *Encoder) writeInt32(v int32) {
	e.put32(uint32(v))
}

func (e *Encoder) writeInt64(v int64) {
	e.put64(uint64(v))
}

func (e *Encoder) writeUint8(v uint8) {
	e.put8(v)
}

func (e *Encoder) writeUint16(v uint16) {
	e.put16(v)
}

func (e *Encoder) writeUint32(v uint32) {
	e.put32(v)
}

func (e *Encoder) writeUint64(v uint64) {
	e.put64(v)
}

func (e *Encoder) put8(v uint8) {
	e.bytes = append(e.bytes, v)
}

func (e *Encoder) put16(v uint16) {
	b := make([]byte, 2)
	binary.BigEndian.PutUint16(b, v)
	e.bytes = append(e.bytes, b...)
}

func (e *Encoder) put32(v uint32) {
	b := make([]byte, 4)
	binary.BigEndian.PutUint32(b, v)
	e.bytes = append(e.bytes, b...)
}

func (e *Encoder) put64(v uint64) {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, v)
	e.bytes = append(e.bytes, b...)
}
