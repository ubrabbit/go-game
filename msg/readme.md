# msg

#### 介绍
协议的注册需要在好几个地方注册。暂时不是很方便，后期考虑优化可能。

#### 注册一个协议时需要做的事情
1）在msg/protocol目录：
	a) 在protocol.go填写协议编号，协议值不能大于uint8(255)。我们约定GS2C和C2GS使用同一个值，方便管理。
	b) 创建该协议的结构体。结构体需要实现 Protocol、PacketData、UnpackData三个函数。
	c) 在msg/msg.go注册协议。
		举例：
		func init() {
			Processor.Register(&TestEcho{})
		}

2）在gate/router.go文件注册协议，用于告知这个协议由哪个模块(独立的goroutine)处理。
	例如： msg.Processor.SetRouter(&protocol.TestEcho{}, game.ChanRPC) 。
	表示网关收到这个协议时，会转发给运行game这个模块的goroutine。

3） 在模块的internal子目录注册协议处理函数。这里以game为例，就是game/internal子目录。
	在handler.go里面的init注册协议处理函数。
	举例：
		func init() {
			setHandler(&C2GSLoadRole{}, command.HandleC2GSLoadRole)
		}
