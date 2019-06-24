package internal

import (
	"server/event"
	"server/game/net"
)
import (
	. "server/common"
)

func GameEventCallback(args []interface{}) {
	defer func() {
		r := recover()
		if r != nil {
			LogError("Callback Error: %v %v", args, r)
		}
	}()
	LogDebug("GameEventCallback")
	f := args[0].(*Functor)
	f.Call(args[1])
}

func init() {
	LogInfo("init game event")

	skeleton.RegisterChanRPC("NewAgent", net.RpcNewAgent)
	skeleton.RegisterChanRPC("CloseAgent", net.RpcCloseAgent)

	event.RegistEventModule(event.EVENT_MODULE_GAME, ChanRPC)
	skeleton.RegisterChanRPC(event.EVENT_MODULE_GAME, GameEventCallback)
}
