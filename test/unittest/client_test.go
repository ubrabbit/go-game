package unittest

import (
	"testing"
)

import (
	. "server/test/common"
)

const (
	constClientCount = 10
	constLoopCount   = 10
	constContentSize = 10240
)

func TestClientEcho(t *testing.T) {
	c := NewClient()
	c.TestEcho(constClientCount, constLoopCount)
}

func TestPacketSendNoProtocol(t *testing.T) {
	c := NewClient()
	c.PacketSendNoProtocol(constContentSize)
}

func TestPacketSendLengh0(t *testing.T) {
	c := NewClient()
	c.PacketSendLengh0(constContentSize)
}

func TestPacketSendContent0(t *testing.T) {
	c := NewClient()
	c.PacketSendContent0(constContentSize)
}

func TestPacketSendContentLonger(t *testing.T) {
	c := NewClient()
	c.PacketSendContentLonger(constContentSize)
}

func TestPacketSendContentTooLarge(t *testing.T) {
	c := NewClient()
	c.PacketSendContentTooLarge()
}

func TestPacketSendOnlyProtocol(t *testing.T) {
	c := NewClient()
	c.PacketSendOnlyProtocol()
}

func TestPacketOnlyConnect(t *testing.T) {
	PacketOnlyConnect(1, 1000)
}
