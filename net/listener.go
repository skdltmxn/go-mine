package net

import (
	"log"
	"net"
	"strconv"
	"time"
)

type Listener struct {
	dispatchers []Dispatcher
}

func NewListener() *Listener {
	return &Listener{}
}

func (l *Listener) Run(port int) {
	listener, err := net.Listen("tcp", ":"+strconv.Itoa(port))
	if err != nil {
		log.Fatal("Failed to listen")
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal("failed to accept client")
		}

		go l.handleClient(newSession(conn))
	}
}

func (l *Listener) RegisterDispatcher(dispatcher Dispatcher) {
	l.dispatchers = append(l.dispatchers, dispatcher)
}

func (l *Listener) handleClient(sess *Session) {
	log.Println("new client")
	defer sess.Close()

	go sess.receiveData()

	for !sess.eof {
		packet := sess.getNextPacket()

		if packet == nil {
			time.Sleep(time.Millisecond)
			continue
		}

		for _, d := range l.dispatchers {
			if d.Dispatch(sess, packet) {
				break
			}
		}
	}

	log.Println("client disconnected")
}
