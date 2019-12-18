package net

import (
	"github.com/skdltmxn/go-mine/net/packet"
)

type Dispatcher interface {
	Dispatch(sess *Session, p *packet.Packet) bool
}
