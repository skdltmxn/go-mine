package net

import (
	"bufio"
	"bytes"
	"net"
	"sync"

	"github.com/skdltmxn/go-mine/net/packet"
)

const (
	SessionStateStatus = 1
	SessionStateLogin  = 2
)

type Session struct {
	conn   net.Conn
	buffer bytes.Buffer
	eof    bool
	m      sync.Mutex
	state  int
}

func newSession(conn net.Conn) *Session {
	return &Session{
		conn:  conn,
		eof:   false,
		state: SessionStateStatus,
	}
}

func (sess *Session) State() int {
	return sess.state
}

func (sess *Session) SetState(newState int) {
	if newState > SessionStateLogin {
		newState = SessionStateStatus
	}
	sess.state = newState
}

func (sess *Session) Close() {
	sess.conn.Close()
}

func (sess *Session) SendData(data []byte) (int, error) {
	return sess.conn.Write(data)
}

func (sess *Session) receiveData() {
	t := make([]byte, 4096)
	reader := bufio.NewReader(sess.conn)

	for !sess.eof {
		n, err := reader.Read(t)
		if err != nil {
			sess.eof = true
			break
		}

		sess.m.Lock()
		sess.buffer.Write(t[:n])
		sess.m.Unlock()
	}
}

func (sess *Session) getNextPacket() *packet.Packet {
	defer sess.m.Unlock()

	sess.m.Lock()
	data := sess.buffer.Bytes()

	p, n := packet.ParsePacket(data)
	if n == 0 {
		return nil
	} else if n < 0 {
		sess.eof = true
		return nil
	}

	sess.buffer.Next(n)

	return p
}
