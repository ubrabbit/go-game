package internal

import (
	"server/event"
)
import (
	. "server/common"
)

func DBEventCallback(args []interface{}) {
	defer func() {
		r := recover()
		if r != nil {
			LogError("Callback Error: %v %v", args, r)
		}
	}()
	//LogDebug("DBEventCallback")
	f := args[0].(*Functor)
	f.Call(args[1])
}

func init() {
	LogInfo("init db event")
	event.RegistEventModule(event.EVENT_MODULE_DB, ChanRPC)
	skeleton.RegisterChanRPC(event.EVENT_MODULE_DB, DBEventCallback)
}
