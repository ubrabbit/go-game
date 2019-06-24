package world

import (
	"runtime"
	"server/conf"
	"server/timer"
)

import (
	. "server/common"
)

func (w *World) ID() int {
	return w.id
}

func (w *World) Create() {
	timer.AddObject(w.ID())
}

func (w *World) Delete() {
	timer.RemoveObject(w.ID())
}

func (w *World) ServerNum() int {
	return conf.Server.ServerNum
}

func (w *World) StartTimer(key string, timeout int, f *Functor) {
	timer.StartTimer(timer.TIMER_MODULE_GAME, w.ID(), key, timeout, f)
}

func (w *World) StartTimerMs(key string, timeout int, f *Functor) {
	timer.StartTimerMs(timer.TIMER_MODULE_GAME, w.ID(), key, timeout, f)
}

func (w *World) RemoveTimer(key string) bool {
	return timer.RemoveTimer(timer.TIMER_MODULE_GAME, w.ID(), key)
}

func (w *World) Heartbeat(args ...interface{}) {
	w.RemoveTimer("Heartbeat")
	w.StartTimer("Heartbeat", heartbeatInternal, NewFunctor("Heartbeat", w.Heartbeat))
	w.heartbeatExecute()
}

func (w *World) heartbeatExecute() {
	defer func() {
		r := recover()
		if r != nil {
			LogInfo("world heartbeatExecute error: %v", r)
		}
	}()
	LogInfo("world gc")
	runtime.GC()
}

func InitWorld() {
	g_World.Heartbeat()
}

func init() {
	LogInfo("init world")
	g_World = &World{id: NewObjectID()}
	g_World.Create()
}
