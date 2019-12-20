package main

import (
	"flag"

	"github.com/skdltmxn/go-mine/net"
	"github.com/skdltmxn/go-mine/server"
)

func main() {
	listener := net.NewListener()
	listener.RegisterDispatcher(server.NewHandshakeServer())
	listener.RegisterDispatcher(server.NewLoginServer())
	listener.RegisterDispatcher(server.NewGameServer())

	portPtr := flag.Int("port", 25565, "Port number for server")
	flag.Parse()

	listener.Run(*portPtr)
}
