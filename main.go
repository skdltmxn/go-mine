package main

import (
	"github.com/skdltmxn/go-mine/net"
	"github.com/skdltmxn/go-mine/server"
)

func main() {
	listener := net.NewListener()
	listener.RegisterDispatcher(server.NewHandshakeServer())
	listener.RegisterDispatcher(server.NewLoginServer())

	listener.Run(25565)
}
