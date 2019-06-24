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
	sync.Mutex
	ID       int
	Key      string
	event    string
	timer    *wheel.Timer
	callback *Functor
	chanRPC  *chanrpc.Server
	stopped  bool
}

type TimerObject struct {
	sync.Mutex
	ID int
	/*
		我认为并没有必要给同一个key的定时器支持同时挂多个定时器的设计。这种设计在写业务逻辑时是大坑。
		如果真的需要同一个key可以挂多个定时器， 数据结构是：
		timerList map[string]map[int]*TimerItem
		同时，对最里面那层map的数据清除，在startTimer里面添加相同key时做一次即可。
		不要在TimerItem的回调里删除，否则代码逻辑就太乱了。
		-- lpx  2019-06-22
	*/
	timerList map[string]*TimerItem
}

var g_Lock sync.Mutex
var g_Module map[string]*TimerModule
var g_Container map[int]*TimerObject
