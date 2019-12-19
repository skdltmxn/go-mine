package packet

import (
	"encoding/binary"
	"io"
)

type Reader struct {
	off int
	buf []byte
}

func NewReader(p *Packet) *Reader {
	return &Reader{
		0,
		p.Data(),
	}
}

func (reader *Reader) Len() int {
	return len(reader.buf) - reader.off
}

func (reader *Reader) Read(p []byte) (n int, err error) {
	if reader.off >= len(reader.buf) {
		return 0, io.EOF
	}

	n = copy(p, reader.buf[reader.off:])
	reader.off += n
	return
}

func (reader *Reader) ReadBoolean() (value bool, err error) {
	err = binary.Read(reader, binary.BigEndian, &value)
	return
}

func (reader *Reader) ReadByte() (value int8, err error) {
	err = binary.Read(reader, binary.BigEndian, &value)
	return
}

func (reader *Reader) ReadUbyte() (value uint8, err error) {
	err = binary.Read(reader, binary.BigEndian, &value)
	return
}

func (reader *Reader) ReadShort() (value int16, err error) {
	err = binary.Read(reader, binary.BigEndian, &value)
	return
}

func (reader *Reader) ReadUshort() (value uint16, err error) {
	err = binary.Read(reader, binary.BigEndian, &value)
	return
}

func (reader *Reader) ReadInt() (value int32, err error) {
	err = binary.Read(reader, binary.BigEndian, &value)
	return
}

func (reader *Reader) ReadLong() (value int64, err error) {
	err = binary.Read(reader, binary.BigEndian, &value)
	return
}

func (reader *Reader) ReadFloat() (value float32, err error) {
	err = binary.Read(reader, binary.BigEndian, &value)
	return
}

func (reader *Reader) ReadDouble() (value float64, err error) {
	err = binary.Read(reader, binary.BigEndian, &value)
	return
}

func (reader *Reader) ReadVarint() (int, error) {
	v, n := binary.Uvarint(reader.buf[reader.off:])
	if n == 0 || n > 5 {
		return 0, io.EOF
	} else if n < 0 {
		return 0, io.ErrShortBuffer
	}

	reader.off += n
	return int(v), nil
}

func (reader *Reader) ReadVarlong() (int64, error) {
	v, n := binary.Uvarint(reader.buf[reader.off:])
	if n == 0 {
		return 0, io.EOF
	} else if n < 0 {
		return 0, io.ErrShortBuffer
	}

	reader.off += n
	return int64(v), nil
}

func (reader *Reader) ReadString() (value string, err error) {
	length, err := reader.ReadVarint()
	if err != nil {
		return "", err
	}

	if length > reader.Len() {
		return "", io.ErrShortBuffer
	}

	value = string(reader.buf[reader.off : reader.off+length])

	reader.off += length
	return
}
