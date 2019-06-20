package msg

import (
	"server/leaf/network/packet"
	. "server/msg/protocol"
)

// 使用默认的 JSON 消息处理器（默认还提供了 protobuf 消息处理器）
var Processor = packet.NewProcessor()

// 一个结构体定义了一个 JSON 消息的格式
// 消息名为 HelloJson
//type HelloJson struct {
//	Name string
//}

func init() {
	// 这里我们注册了一个 JSON 消息 HelloJson
	//.Register(&HelloJson{})
	Processor.Register(&C2GSHello{})
	Processor.Register(&C2GSIdentity{})
	Processor.Register(&C2GSLogin{})
	Processor.Register(&C2GSRoleID{})
	Processor.Register(&C2GSLoadRole{})
	Processor.Register(&C2GSLoginFinished{})
	Processor.Register(&TestEcho{})
}
