package unittest

import (
	"bytes"
	"sync"
	"testing"
)

import (
	"math"
	. "server/common"
	. "server/msg/protocol"
	. "server/test/common"
)

const (
	constClientCount = 10
	constLoopCount   = 10
	constContentSize = 10240
)

func TestClientEcho(t *testing.T) {
	wg := new(sync.WaitGroup)
	wg.Add(constClientCount)
	ch := make(chan int, constClientCount)
	go func() {
		for i := 0; i < constClientCount; i++ {
			ch <- i
		}
	}()
	for i := 0; i < constClientCount; i++ {
		go func() {
			k := <-ch
			str := FormatString("string_%d", k)
			b := bytes.NewBufferString(FormatString("bytes_%d", k)).Bytes()
			c := NewClient()
			for j := 0; j < constLoopCount; j++ {
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
