package packet

import (
	"encoding/binary"
)

type Writer struct {
	p *Packet
}

func NewWriter(p *Packet) *Writer {
	return &Writer{p}
}

func (w *Writer) Bytes() []byte {
	return w.p.data.Bytes()
}

func (w *Writer) Len() int {
	return w.p.data.Len()
}

func (w *Writer) Write(p []byte) (n int, err error) {
	return w.p.data.Write(p)
}

func (w *Writer) WriteBool(v bool) error {
	return binary.Write(w, binary.BigEndian, &v)
}

func (w *Writer) WriteByte(v int8) error {
	return binary.Write(w, binary.BigEndian, &v)
}

func (w *Writer) WriteUbyte(v uint8) error {
	return binary.Write(w, binary.BigEndian, &v)
}

func (w *Writer) WriteShort(v int16) error {
	return binary.Write(w, binary.BigEndian, &v)
}

func (w *Writer) WriteUshort(v uint16) error {
	return binary.Write(w, binary.BigEndian, &v)
}

func (w *Writer) WriteInt(v int32) error {
	return binary.Write(w, binary.BigEndian, &v)
}

func (w *Writer) WriteLong(v int64) error {
	return binary.Write(w, binary.BigEndian, &v)
}

func (w *Writer) WriteFloat(v float32) error {
	return binary.Write(w, binary.BigEndian, &v)
}

func (w *Writer) WriteDouble(v float64) error {
	return binary.Write(w, binary.BigEndian, &v)
}

func (w *Writer) WriteVarint(v int) error {
	b := make([]byte, 5)
	n := binary.PutUvarint(b, uint64(v))
	_, err := w.p.data.Write(b[:n])
	return err
}

func (w *Writer) WriteVarlong(v int64) error {
	b := make([]byte, 10)
	n := binary.PutUvarint(b, uint64(v))
	_, err := w.Write(b[:n])
	return err
}

func (w *Writer) WriteString(v string) error {
	if err := w.WriteVarint(len(v)); err != nil {
		return err
	}

	_, err := w.Write([]byte(v))
	return err
}
