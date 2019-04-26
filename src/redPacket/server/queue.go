package server

import (
	"math/rand"
	"time"
)

type task struct {
	id       uint32
	callback chan uint
}

var taskNum int = 16

var tasks = make([]chan task, taskNum)

func fetchPacketService(chTasks chan task) {
	for {
		t := <-chTasks
		id := t.id
		list, ok := packetList.Load(id)
		list1 := list.([]uint)
		l := len(list1)
		if ok && l > 0 {
			r := rand.New(rand.NewSource(time.Now().UnixNano()))
			i := r.Intn(l)
			money := list1[i]
			if l > 1 {
				if i == l-1 {
					packetList.Store(id, list1[:i])
				} else if i == 0 {
					packetList.Store(id, list1[1:])
				} else {
					packetList.Store(id, append(list1[0:i], list1[i+1:]...))
				}
			} else {
				packetList.Store(id, list1[0:0])
			}
			t.callback <- money
		} else {
			t.callback <- 0
		}
	}
}
