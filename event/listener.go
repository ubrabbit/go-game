package event

import (
	. "server/common"
)

func CreateListener(module string) *Listener {
	l := new(Listener)
	l.listenList = make(map[int]map[int]*Functor, 0)
	l.module = module
	return l
}

func (l *Listener) AddListen(e Event, id int, f *Functor) bool {
	defer l.Unlock()
	l.Lock()

	eid := e.ID()
	_, ok := l.listenList[eid]
	if !ok {
		l.listenList[eid] = make(map[int]*Functor, 0)
	}
	l.listenList[eid][id] = f
	return true
}

func (l *Listener) RemoveListen(e Event, id int) {
	defer l.Unlock()
	l.Lock()

	list, ok := l.listenList[e.ID()]
	if ok {
		_, ok = list[id]
		if ok {
			delete(list, id)
		}
	}
}

/*
Event和Timer有一样的问题，在触发event并进入channel等待执行时，Event被所属的goroutine删除了，此时这个触发的event仍然会在RemoveEvent后执行。
不过Event和Timer的不同是，Event回调时本身就该根据Event的参数做足条件判断，而Timer属于定时任务，不一定每个任务都需要条件判断。
如果有必要，以后可以给被删除的Event打一个IsAlive标记，但现在不需要。 --- lpx 2019-06-22
*/
func (l *Listener) TriggerEvent(e Event) {
	var listFuncs []*Functor
	l.Lock()
	list, ok := l.listenList[e.ID()]
	if ok {
		for _, f := range list {
			listFuncs = append(listFuncs, f)
		}
	}
	l.Unlock()

	if ok {
		for _, f := range listFuncs {
			ch := getModule(l.module)
			/*
				加锁不能包含ch.Go。考虑一种情况： ch.Go所推送的channel满了，阻塞中。此时监听这个channel的goroutine正试图往Listener添加一个事件。
				因为Listener加锁了，此时这个channel所属的goroutine也会被阻塞。这样就导致了死锁！
			*/
			ch.Go(l.module, f, e)
		}
	}
}
