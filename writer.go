package msgpack

import (
	"encoding/binary"
	"math"
)

type writer struct {
	bytes []byte
}

func (w *writer) clear() {
	w.bytes = nil
}

func (w *writer) writeByte(v byte) {
	w.bytes = append(w.bytes, v)
}

func (w *writer) writeBytes(v []byte) {
	w.bytes = append(w.bytes, v...)
}

func (w *writer) writeFloat32(v float32) {
	w.put32(math.Float32bits(v))
}

func (w *writer) writeFloat64(v float64) {
	w.put64(math.Float64bits(v))
}

func (w *writer) writeInt8(v int8) {
	w.put8(uint8(v))
}

func (w *writer) writeInt16(v int16) {
	w.put16(uint16(v))
}

func (w *writer) writeInt32(v int32) {
	w.put32(uint32(v))
}

func (w *writer) writeInt64(v int64) {
	w.put64(uint64(v))
}

func (w *writer) writeUint8(v uint8) {
	w.put8(v)
}

func (w *writer) writeUint16(v uint16) {
	w.put16(v)
}

func (w *writer) writeUint32(v uint32) {
	w.put32(v)
}

func (w *writer) writeUint64(v uint64) {
	w.put64(v)
}

func (w *writer) put8(v uint8) {
	w.bytes = append(w.bytes, v)
}

func (w *writer) put16(v uint16) {
	b := make([]byte, 2)
	binary.BigEndian.PutUint16(b, v)
	w.bytes = append(w.bytes, b...)
}

func (w *writer) put32(v uint32) {
	b := make([]byte, 4)
	binary.BigEndian.PutUint32(b, v)
	w.bytes = append(w.bytes, b...)
}

func (w *writer) put64(v uint64) {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, v)
	w.bytes = append(w.bytes, b...)
}
