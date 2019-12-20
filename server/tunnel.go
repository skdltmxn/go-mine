package server

import "github.com/skdltmxn/go-mine/net"

var tunnel = make(chan *DataTunnel)

type DataTunnel struct {
	sess *net.Session
	name string
	eid  int32
}

func getTunnelReceiver() <-chan *DataTunnel {
	return tunnel
}

func getTunnelSender() chan<- *DataTunnel {
	return tunnel
}

func newDataTunnel(sess *net.Session, name string, eid int32) *DataTunnel {
	return &DataTunnel{sess, name, eid}
}
