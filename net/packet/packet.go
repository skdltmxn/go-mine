package packet

import (
	"bytes"
	"encoding/binary"
)

type Packet struct {
	id   int
	data bytes.Buffer
}

func NewPacket(id int) *Packet {
	return &Packet{id: id}
}

func ParsePacket(rawData []byte) (*Packet, int) {
	length, lengthBytes := binary.Uvarint(rawData)

	if lengthBytes == 0 {
		// too short packet
		return nil, 0
	} else if lengthBytes < 0 {
		// invalid packet
		return nil, -1
	}

	// too short
	if len(rawData[lengthBytes:]) < int(length) {
		return nil, 0
	}

	payload := rawData[lengthBytes : lengthBytes+int(length)]
	id, idBytes := binary.Uvarint(payload)

	if idBytes == 0 {
		// too short packet
		return nil, 0
	} else if idBytes < 0 {
		// invalid packet
		return nil, -1
	}

	p := &Packet{id: int(id)}
	p.data.Write(payload[idBytes:])

	return p, int(length) + lengthBytes
}

func (p *Packet) Id() int {
	return p.id
}

func (p *Packet) Data() []byte {
	return p.data.Bytes()
}

func (p *Packet) Raw() []byte {
	id := make([]byte, 5)
	idN := binary.PutUvarint(id, uint64(p.id))

	length := make([]byte, 5)
	lengthN := binary.PutUvarint(length, uint64(p.data.Len()+idN))

	raw := make([]byte, lengthN+idN+p.data.Len())
	i := copy(raw, length[:lengthN])
	i += copy(raw[i:], id[:idN])
	i += copy(raw[i:], p.data.Bytes())

	return raw
}
