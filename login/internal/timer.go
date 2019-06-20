package internal

import (
	"server/timer"
)

import (
	. "server/common"
)

func LoginTimerCallback(args []interface{}) {
	LogDebug("LoginTimerCallback")
	f := args[0].(Functor)
	f.Call()
}

func init() {
	LogInfo("init login timer")
	timer.RegistModule(timer.TIMER_MODULE_LOGIN, timer.TIMER_EVENT_LOGIN, ChanRPC)
	skeleton.RegisterChanRPC(timer.TIMER_EVENT_LOGIN, LoginTimerCallback)
}
