package common

import (
	"encoding/binary"
	"math"
	"net"
	"strings"
	"time"
)
import (
	. "server/common"
	. "server/leaf/network"
	. "server/msg/protocol"
)

const (
	packetHeaderSize    = 3
	debugShowPacketSize = 32
	constSyncfloodCount = 10000
)

func debugPrint(s string, b []byte) {
	l := debugShowPacketSize
	if len(b) < debugShowPacketSize {
		l = len(b)
	}
	LogDebug("%s: %x", s, b[:l])
}

//发送不存在的协议
func (c *Client) PacketSendNoProtocol(size int) {
	c.SendData = make([]byte, 0)
	str := strings.Repeat("A", size)
	c.SendData = PacketString(c.SendData, str, len(str))
	debugPrint("PacketSendNoProtocol", c.SendData)
	c.SendProto(0x00)
}

//长度与实际内容不匹配1： 长度为0,有内容
func (c *Client) PacketSendLengh0(size int) {
	c.SendData = make([]byte, 0)
	str := strings.Repeat("A", size)
	msgData := PacketString(c.SendData, str, len(str))

	msg := make([]byte, len(str)+packetHeaderSize)
	msg[0] = byte(TEST_ECHO)
	//长度写入0
	binary.LittleEndian.PutUint16(msg[1:], uint16(0))
	copy(msg[packetHeaderSize:], msgData)
	debugPrint("PacketSendLengh0", msg)
	c.Conn.Write(msg)
}

//长度与实际内容不匹配2： 有长度,内容为0
func (c *Client) PacketSendContent0(size int) {
	c.SendData = make([]byte, 0)
	str := strings.Repeat("A", size)
	//实际写入内容为0
	msgData := make([]byte, 0)
	packetSize := len(str)
	msg := make([]byte, len(str)+packetHeaderSize)
	msg[0] = byte(TEST_ECHO)
	if LittleEndian {
		binary.LittleEndian.PutUint16(msg[1:], uint16(packetSize))
	} else {
		binary.BigEndian.PutUint16(msg[1:], uint16(packetSize))
	}
	copy(msg[packetHeaderSize:], msgData)
	debugPrint("PacketSendContent0", msg)
	c.Conn.Write(msg)
}

//长度与实际内容不匹配3： 有长度,实际内容超出总长度
func (c *Client) PacketSendContentLonger(size int) {
	c.SendData = make([]byte, 0)
	str := strings.Repeat("A", size)
	msgData := PacketString(c.SendData, str, len(str))

	msg := make([]byte, size+packetHeaderSize)
	msg[0] = byte(TEST_ECHO)
	//写入长度超出实际总长度
	fakeSize := len(str) * 2
	if LittleEndian {
		binary.LittleEndian.PutUint16(msg[1:], uint16(fakeSize))
	} else {
		binary.BigEndian.PutUint16(msg[1:], uint16(fakeSize))
	}
	copy(msg[packetHeaderSize:], msgData)
	debugPrint("PacketSendContentLonger", msg)
	c.Conn.Write(msg)
}

//发送一个超长的内容串
func (c *Client) PacketSendContentTooLarge() {
	c.SendData = make([]byte, 0)
	size := 400 * 1024 * 1024 //400M
	str := strings.Repeat("A", size)
	msgData := PacketString(c.SendData, str, len(str))

	msg := make([]byte, size+packetHeaderSize)
	msg[0] = byte(TEST_ECHO)
	//写入长度
	packetSize := len(str)
	if packetSize > math.MaxUint16 {
		packetSize = math.MaxUint16
	}
	if LittleEndian {
		binary.LittleEndian.PutUint16(msg[1:], uint16(packetSize))
	} else {
		binary.BigEndian.PutUint16(msg[1:], uint16(packetSize))
	}
	copy(msg[packetHeaderSize:], msgData)
	debugPrint("PacketSendContentTooLarge", msg)
	c.Conn.Write(msg)
}

//除了协议号什么都没发
func (c *Client) PacketSendOnlyProtocol() {
	msgData := make([]byte, 0)
	msg := make([]byte, packetHeaderSize)
	msg[0] = byte(TEST_ECHO)
	if LittleEndian {
		binary.LittleEndian.PutUint16(msg[1:], uint16(0))
	} else {
		binary.BigEndian.PutUint16(msg[1:], uint16(0))
	}
	copy(msg[packetHeaderSize:], msgData)
	debugPrint("PacketSendOnlyProtocol", msg)
	c.Conn.Write(msg)
}

//除了发起连接什么都不干
func PacketOnlyConnect(timeout int, count int) {
	for i := 0; i < count; i++ {
		go func() (err error) {
			defer func() {
				r := recover()
				if r != nil {
					err = r.(error)
				}
			}()
			conn, err := net.Dial("tcp", SERVER_ADDR)
			if err != nil {
				return err
			}
			time.Sleep(time.Duration(timeout-10) * time.Second)

			msgData := make([]byte, 0)
			msg := make([]byte, packetHeaderSize)
			msg[0] = byte(TEST_ECHO)
			if LittleEndian {
				binary.LittleEndian.PutUint16(msg[1:], uint16(0))
			} else {
				binary.BigEndian.PutUint16(msg[1:], uint16(0))
			}
			copy(msg[packetHeaderSize:], msgData)
			conn.Write(msg)
			return nil
		}()
	}
	time.Sleep(time.Duration(timeout) * time.Second)
}
