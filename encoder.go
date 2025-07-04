package msgpack

import (
	"fmt"
	"time"
)

func NewEncoder() *Encoder {
	return &Encoder{
		writer: &writer{},
	}
}

type Encoder struct {
	writer *writer
}

func (e *Encoder) Bytes() []byte {
	return e.writer.bytes
}

func (e *Encoder) Clear() {
	e.writer.clear()
}

func (e *Encoder) PutArrayLength(v uint32) {
	if v <= mask4 {
		e.writer.writeByte(byte(0x90 | v))
	} else if v <= mask16 {
		e.writer.writeByte(0xdc)
		e.writer.writeUint16(uint16(v))
	} else {
		e.writer.writeByte(0xdd)
		e.writer.writeUint32(uint32(v))
	}
}

func (e *Encoder) PutBool(v bool) {
	if !v {
		e.writer.writeByte(0xc2)
	} else {
		e.writer.writeByte(0xc3)
	}
}

func (e *Encoder) PutBytes(v []byte) error {
	n := len(v)
	if n <= mask8 {
		e.writer.writeByte(0xc4)
		e.writer.writeUint8(uint8(n))
	} else if n <= mask16 {
		e.writer.writeByte(0xc5)
		e.writer.writeUint16(uint16(n))
	} else if n <= mask32 {
		e.writer.writeByte(0xc6)
		e.writer.writeUint32(uint32(n))
	} else {
		return fmt.Errorf("byte slice (%d bytes) is too long to encode", n)
	}
	e.writer.writeBytes(v)
	return nil
}

func (e *Encoder) PutFloat(v float64) {
	// Try to work out if and when encoding a float32 is acceptable!
	// TODO
	e.PutFloat64(v)
}

func (e *Encoder) PutFloat32(v float32) {
	e.writer.writeByte(0xca)
	e.writer.writeFloat32(v)
}

func (e *Encoder) PutFloat64(v float64) {
	e.writer.writeByte(0xcb)
	e.writer.writeFloat64(v)
}

func (e *Encoder) PutInt(v int64) {
	if v >= intFixMin && v <= intFixMax {
		e.writer.writeInt8(int8(v))
	} else if v >= int8Min && v <= int8Max {
		e.writer.writeByte(0xd0)
		e.writer.writeInt8(int8(v))
	} else if v >= int16Min && v <= int16Max {
		e.writer.writeByte(0xd1)
		e.writer.writeInt16(int16(v))
	} else if v >= int32Min && v <= int32Max {
		e.writer.writeByte(0xd2)
		e.writer.writeInt32(int32(v))
	} else {
		e.writer.writeByte(0xd3)
		e.writer.writeInt64(v)
	}
}

func (e *Encoder) PutMapLength(v uint32) {
	if v <= mask4 {
		e.writer.writeByte(byte(0x80 | v))
	} else if v <= mask16 {
		e.writer.writeByte(0xde)
		e.writer.writeUint16(uint16(v))
	} else {
		e.writer.writeByte(0xdf)
		e.writer.writeUint32(uint32(v))
	}
}

func (e *Encoder) PutNil() {
	e.writer.writeByte(0xc0)
}

func (e *Encoder) PutString(v string) error {
	u := []byte(v)
	n := len(u)
	if n <= mask5 {
		e.writer.writeByte(byte(0xa0 | n))
	} else if n <= mask8 {
		e.writer.writeByte(0xd9)
		e.writer.writeUint8(uint8(n))
	} else if n <= mask16 {
		e.writer.writeByte(0xda)
		e.writer.writeUint16(uint16(n))
	} else if n <= mask32 {
		e.writer.writeByte(0xdb)
		e.writer.writeUint32(uint32(n))
	} else {
		return fmt.Errorf("string (%d bytes) is too long to encode", n)
	}
	e.writer.writeBytes(u)
	return nil
}

func (e *Encoder) PutTime(v time.Time) {
	sec := v.Unix()
	nsec := v.Nanosecond()
	if sec < 0 || sec > mask34 {
		//timestamp 96
		e.writer.writeByte(0xc7)
		e.writer.writeByte(12)
		e.writer.writeByte(0xff)
		e.writer.writeUint32(uint32(nsec))
		e.writer.writeInt64(sec)
	} else if sec > mask32 || nsec > 0 {
		//timestamp 64
		e.writer.writeByte(0xd7)
		e.writer.writeByte(0xff)
		data64 := uint64(nsec)<<34 | uint64(sec)
		e.writer.writeUint64(data64)
	} else {
		// timestamp 32
		e.writer.writeByte(0xd6)
		e.writer.writeByte(0xff)
		e.writer.writeInt32(int32(sec))
	}
}

func (e *Encoder) PutUint(v uint64) {
	if v <= mask7 {
		e.writer.writeByte(byte(v))
	} else if v <= mask8 {
		e.writer.writeByte(0xcc)
		e.writer.writeUint8(uint8(v))
	} else if v <= mask16 {
		e.writer.writeByte(0xcd)
		e.writer.writeUint16(uint16(v))
	} else if v <= mask32 {
		e.writer.writeByte(0xce)
		e.writer.writeUint32(uint32(v))
	} else {
		e.writer.writeByte(0xcf)
		e.writer.writeUint64(v)
	}
}
