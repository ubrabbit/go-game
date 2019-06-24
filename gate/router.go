package gate

import (
	"server/game"
	"server/login"
	"server/msg"
	"server/msg/protocol"
)

func init() {
	// 这里指定消息 HelloJson 路由到 game 模块
	// 模块间使用 ChanRPC 通讯，消息路由也不例外
	//msg.Processor.SetRouter(&msg.HelloJson{}, game.ChanRPC)

	msg.Processor.SetRouter(&protocol.C2GSHello{}, game.ChanRPC)
	msg.Processor.SetRouter(&protocol.C2GSIdentity{}, game.ChanRPC)

	msg.Processor.SetRouter(&protocol.C2GSLogin{}, login.ChanRPC)
	msg.Processor.SetRouter(&protocol.C2GSRoleID{}, login.ChanRPC)

	msg.Processor.SetRouter(&protocol.C2GSLoadRole{}, game.ChanRPC)
	msg.Processor.SetRouter(&protocol.C2GSLoginFinished{}, game.ChanRPC)

	msg.Processor.SetRouter(&protocol.TestEcho{}, game.ChanRPC)
}
