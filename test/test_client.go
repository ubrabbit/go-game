package main

import (
	. "server/test/common"
)

func PacketSendNoProtocol() {
	c := NewClient()
	c.PacketSendNoProtocol(1024)
}

func PacketSendLengh0() {
	c := NewClient()
	c.PacketSendLengh0(1024)
}

func PacketSendContent0() {
	c := NewClient()
	c.PacketSendContent0(1024)
}

func PacketSendContentLonger() {
	c := NewClient()
	c.PacketSendContentLonger(1024)
}

func PacketSendContentTooLarge() {
	c := NewClient()
	c.PacketSendContentTooLarge()
}

func PacketSendOnlyProtocol() {
	c := NewClient()
	c.PacketSendOnlyProtocol()
}

func main() {
	PacketSendNoProtocol()
	PacketSendLengh0()
	PacketSendContent0()
	PacketSendContentLonger()
	PacketSendContentTooLarge()
	PacketSendOnlyProtocol()
	PacketOnlyConnect(5, 100)
}
