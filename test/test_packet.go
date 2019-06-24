package main

import (
	"math"
	"strings"
	"sync/atomic"
	"time"
)

import (
	. "server/common"
	. "server/leaf/network"
)

var g_Finished bool
var g_Count int64
var g_String string = strings.Repeat("A", 64)
var g_Bytes []byte = []byte(strings.Repeat("A", 64))

var g_BytesPacketInt []byte
var g_BytesPacketIntLen int
var g_BytesPacket []byte
var g_BytesPacketLen int

func test_PacketInt() {
	g_Finished = false
	atomic.StoreInt64(&g_Count, 0)
	for {
		if g_Finished {
			break
		}
		data := make([]byte, 0)
		data = PacketInt(data, math.MaxUint8, 1)
		data = PacketInt(data, math.MaxUint16, 2)
		data = PacketInt(data, math.MaxUint32, 4)
		data = PacketInt(data, math.MaxInt64, 8)
		atomic.AddInt64(&g_Count, 1)
	}
}

func test_PacketString() {
	g_Finished = false
	atomic.StoreInt64(&g_Count, 0)
	for {
		if g_Finished {
			break
		}
		data := make([]byte, 0)
		data = PacketString(data, g_String, len(g_String))
		atomic.AddInt64(&g_Count, 1)
	}
}

func test_PacketBytes() {
	g_Finished = false
	atomic.StoreInt64(&g_Count, 0)
	for {
		if g_Finished {
			break
		}
		data := make([]byte, 0)
		data = PacketBytes(data, g_Bytes, len(g_Bytes))
		atomic.AddInt64(&g_Count, 1)
	}
}

func test_UnpackInt() {
	g_Finished = false
	atomic.StoreInt64(&g_Count, 0)
	for {
		if g_Finished {
			break
		}
		data := make([]byte, g_BytesPacketIntLen)
		copy(data, g_BytesPacketInt)

		UnpackInt(data, 1)
		UnpackInt(data, 2)
		UnpackInt(data, 4)
		UnpackInt(data, 8)
		atomic.AddInt64(&g_Count, 1)
	}
}

func test_UnpackString() {
	g_Finished = false
	atomic.StoreInt64(&g_Count, 0)
	for {
		if g_Finished {
			break
		}
		UnpackString(g_BytesPacket, g_BytesPacketLen)
		atomic.AddInt64(&g_Count, 1)
	}
}

func test_UnpackBytes() {
	g_Finished = false
	atomic.StoreInt64(&g_Count, 0)
	for {
		if g_Finished {
			break
		}
		UnpackBytes(g_BytesPacket, g_BytesPacketLen)
		atomic.AddInt64(&g_Count, 1)
	}
}

func main() {
	g_BytesPacketInt = make([]byte, 0)
	g_BytesPacketInt = PacketInt(g_BytesPacketInt, math.MaxUint8, 1)
	g_BytesPacketInt = PacketInt(g_BytesPacketInt, math.MaxUint16, 2)
	g_BytesPacketInt = PacketInt(g_BytesPacketInt, math.MaxUint32, 4)
	g_BytesPacketInt = PacketInt(g_BytesPacketInt, math.MaxInt64, 8)
	g_BytesPacketIntLen = len(g_BytesPacketInt)

	g_BytesPacket = make([]byte, 0)
	g_BytesPacket = PacketBytes(g_BytesPacket, g_Bytes, len(g_Bytes))
	g_BytesPacketLen = len(g_Bytes)

	go test_PacketInt()
	time.Sleep(1 * time.Second)
	g_Finished = true
	LogInfo("test_PacketInt: %d", g_Count)

	go test_PacketString()
	time.Sleep(1 * time.Second)
	g_Finished = true
	LogInfo("test_PacketString: %d", g_Count)

	go test_PacketBytes()
	time.Sleep(1 * time.Second)
	g_Finished = true
	LogInfo("test_PacketBytes: %d", g_Count)

	go test_UnpackInt()
	time.Sleep(1 * time.Second)
	g_Finished = true
	LogInfo("test_UnpackInt: %d", g_Count)

	go test_UnpackString()
	time.Sleep(1 * time.Second)
	g_Finished = true
	LogInfo("test_UnpackString: %d", g_Count)

	go test_UnpackBytes()
	time.Sleep(1 * time.Second)
	g_Finished = true
	LogInfo("test_UnpackBytes: %d", g_Count)
}
