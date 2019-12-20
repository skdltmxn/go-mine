package server

import (
	"log"

	"github.com/skdltmxn/go-mine/net"
	"github.com/skdltmxn/go-mine/net/packet"
)

type HandshakeServer struct {
}

func NewHandshakeServer() *HandshakeServer {
	return &HandshakeServer{}
}

func (d *HandshakeServer) Dispatch(sess *net.Session, p *packet.Packet) bool {
	if sess.State() != net.SessionStateStatus {
		return false
	}

	//log.Printf("HandshakeServer got packet %d / %+v", p.Id(), hex.EncodeToString(p.Data()))

	if p.Id() == 0 {
		r := packet.NewReader(p)

		protocolVersion, _ := r.ReadVarint()
		serverAddress, _ := r.ReadString()
		serverPort, _ := r.ReadUshort()
		nextState, _ := r.ReadVarint()

		if protocolVersion < 575 {
			log.Printf("Protocol version incompatible: %d", protocolVersion)
			sess.Close()
			return true
		}

		log.Printf("protocol: %d server: %s port: %d state: %d", protocolVersion, serverAddress, serverPort, nextState)
		sess.SetState(net.SessionStateLogin)
	} else {
		sess.Close()
	}

	return true
}
