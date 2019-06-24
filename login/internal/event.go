package internal

import (
	"server/event"
)
import (
	. "server/common"
)

func LoginEventCallback(args []interface{}) {
	defer func() {
		r := recover()
		if r != nil {
			LogError("Callback Error: %v %v", args, r)
		}
	}()
	//LogDebug("LoginEventCallback")
	f := args[0].(*Functor)
	f.Call(args[1])
}

func init() {
	LogInfo("init login event")

	event.RegistEventModule(event.EVENT_MODULE_LOGIN, ChanRPC)
	skeleton.RegisterChanRPC(event.EVENT_MODULE_LOGIN, LoginEventCallback)
}
