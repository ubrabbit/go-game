package internal

import (
	"server/timer"
)

import (
	. "server/common"
)

func LoginTimerCallback(args []interface{}) {
	LogDebug("LoginTimerCallback")
	i := args[0].(*timer.TimerItem)
	i.Execute()
}

func init() {
	LogInfo("init login timer")
	timer.RegistModule(timer.TIMER_MODULE_LOGIN, ChanRPC)
	skeleton.RegisterChanRPC(timer.TIMER_MODULE_LOGIN, LoginTimerCallback)
}
