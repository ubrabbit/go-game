package timer

import (
	"server/leaf/chanrpc"
	"server/timer/wheel"
	"sync"
)

import (
	. "server/common"
)

const (
	timerQueueDefaultLen = 4096
)

type TimerModule struct {
	Name    string
	event   string
	wheel   *wheel.TimingWheel
	chanRPC *chanrpc.Server
}

type TimerItem struct {
	Key      string
	timer    *wheel.Timer
	callback Functor
	event    string
	chanRPC  *chanrpc.Server
}

type TimerObject struct {
	sync.Mutex
	ID        int
	timerList map[string][]*TimerItem
}

var g_Lock sync.Mutex
var g_Module map[string]*TimerModule
var g_Container map[int]*TimerObject
