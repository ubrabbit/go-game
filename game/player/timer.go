package player

import (
	"server/timer"
)

import (
	. "server/common"
)

const (
	heartInternal = 10
)

func (p *Player) ID() int {
	return p.Pid
}

func (p *Player) Create() {
	timer.AddObject(p.ID())
}

func (p *Player) Delete() {
	timer.RemoveObject(p.ID())
}

func (p *Player) StartTimer(key string, timeout int, f Functor) {
	timer.StartTimer(timer.TIMER_MODULE_GAME, p.ID(), key, timeout, f)
}

func (p *Player) StartTimerMs(key string, timeout int, f Functor) {
	timer.StartTimerMs(timer.TIMER_MODULE_GAME, p.ID(), key, timeout, f)
}

func (p *Player) RemoveTimer(key string) bool {
	return timer.RemoveTimer(timer.TIMER_MODULE_GAME, p.ID(), key)
}

func (p *Player) heartbeat(args ...interface{}) {
	LogInfo("%s heartbeat", p.Repr())
	p.RemoveTimer("heartbeat")
	f := NewFunctor("heartbeat", p.heartbeat)
	p.StartTimer("heartbeat", heartInternal, f)
}
