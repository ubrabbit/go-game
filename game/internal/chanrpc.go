package internal

import (
	"server/game/net"
)

import (
	. "server/common"
)

func init() {
	LogInfo("init game rpc")
	skeleton.RegisterChanRPC("NewAgent", net.RpcNewAgent)
	skeleton.RegisterChanRPC("CloseAgent", net.RpcCloseAgent)
}
