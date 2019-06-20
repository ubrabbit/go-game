package common

import (
	"net"
	"sync"
)

import (
	. "server/common"
	. "server/leaf/network"
	"server/leaf/network/packet"
	. "server/msg/protocol"
)

type Client struct {
	sync.Mutex
	Conn     net.Conn
	SendData []byte
	Protocol packet.Packet
}

func NewClient() *Client {
	client := new(Client)
	conn, err := net.Dial("tcp", "127.0.0.1:38320")
	if err != nil {
		LogPanic("NewClient: %v", err)
	}
	client.Conn = conn
	client.SendData = make([]byte, 0)
	return client
}

func (c *Client) UnpackProto() (uint8, []byte, error) {
	proto, msgData, err := UnpackProto(c.Conn)
	return proto, msgData, err
}

func (c *Client) SendProto(proto uint8) {
	conn := c.Conn
	msg := PacketProto(proto, c.SendData)
	//LogInfo("SendProto: %v", msg)
	conn.Write(msg)
}

func (c *Client) C2GSEcho(v1, v2, v3, v4 int, s string, b []byte) {
	p := TestEcho{
		Int1: v1,
		Int2: v2,
		Int3: v3,
		Int4: v4,
		Str:  s,
		Byte: b,
	}
	c.Protocol = &p
	_, c.SendData = p.PacketData()
	//LogInfo("C2GSEcho: %v", c.SendData)
	c.SendProto(p.Protocol())
}

func (c *Client) GS2CEcho() {
	p := TestEcho{}
	c.Protocol = &p
	proto, msg, err := c.UnpackProto()
	p.UnpackData(msg)
	if err != nil {
		LogInfo("GS2CEcho err: %d : %v %v\n", proto, p, err)
	}
}
