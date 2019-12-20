package server

import "sync/atomic"

var currentEntityId int32 = -1

func getNextEntityId() int32 {
	return atomic.AddInt32(&currentEntityId, 1)
}
