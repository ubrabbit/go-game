package common

import (
	"bytes"
	"math"
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
	conn, err := net.Dial("tcp", SERVER_ADDR)
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

func (c *Client) TestEcho(clientCount int, loopCount int) {
	wg := new(sync.WaitGroup)
	wg.Add(clientCount)
	ch := make(chan int, clientCount)
	go func() {
		for i := 0; i < clientCount; i++ {
			ch <- i
		}
	}()
	for i := 0; i < clientCount; i++ {
		go func() {
			k := <-ch
			str := FormatString("string_%d", k)
			b := bytes.NewBufferString(FormatString("bytes_%d", k)).Bytes()
			c := NewClient()
			for j := 0; j < loopCount; j++ {
				v1, v2, v3, v4, v5, v6 := k, k+1, k*2, k*4, str, b
				v1 = v1 % math.MaxUint8
				v2 = v2 % math.MaxUint16
				c.C2GSEcho(v1, v2, v3, v4, v5, v6)
				p0 := c.Protocol.(*TestEcho)
				c.GS2CEcho()
				p := c.Protocol.(*TestEcho)
				if (p.Int1 != v1) || (p.Int2 != v2) || (p.Int3 != v3) || (p.Int4 != v4) || (p.Str != v5) || (string(p.Byte) != string(v6)) {
					LogPanic("TestClientEcho Fail! p: %v p0: %v", p, p0)
				}
				//LogInfo("%d(%d) response success", k, j)
			}
			wg.Done()
		}()
	}
	wg.Wait()
}
