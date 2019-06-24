package main

//go test -bench=. -benchmem -benchtime="1s"
/*
2019-06-24:
    goos: linux
    goarch: amd64
    pkg: server/test/bench
    BenchmarkPacketInt-8        10000000           156 ns/op          48 B/op          8 allocs/op
    BenchmarkPacketString-8     20000000            68.2 ns/op       128 B/op          2 allocs/op
    BenchmarkPacketBytes-8      50000000            33.9 ns/op        64 B/op          1 allocs/op
    BenchmarkUnpackInt-8         2000000          1011 ns/op         224 B/op         12 allocs/op
    BenchmarkUnpackString-8     30000000            48.3 ns/op        64 B/op          1 allocs/op
    BenchmarkUnpackBytes-8      200000000            6.15 ns/op        0 B/op          0 allocs/op
    PASS
    ok      server/test/bench   59.055s
*/

import (
	"math"
	"strings"
	"testing"
)

import (
	. "server/leaf/network"
)

func BenchmarkPacketInt(b *testing.B) {
	for n := 0; n < b.N; n++ {
		data := make([]byte, 0)
		data = PacketInt(data, math.MaxUint8, 1)
		data = PacketInt(data, math.MaxUint16, 2)
		data = PacketInt(data, math.MaxUint32, 4)
		data = PacketInt(data, math.MaxInt64, 8)
	}
}

func BenchmarkPacketString(b *testing.B) {
	b.StopTimer()
	str := strings.Repeat("A", 64)
	b.StartTimer()
	for n := 0; n < b.N; n++ {
		data := make([]byte, 0)
		data = PacketString(data, str, len(str))
	}
}

func BenchmarkPacketBytes(b *testing.B) {
	b.StopTimer()
	str := []byte(strings.Repeat("A", 64))
	b.StartTimer()
	for n := 0; n < b.N; n++ {
		data := make([]byte, 0)
		data = PacketBytes(data, str, len(str))
	}
}

func BenchmarkUnpackInt(b *testing.B) {
	for n := 0; n < b.N; n++ {
		b.StopTimer()
		data := make([]byte, 0)
		data = PacketInt(data, math.MaxUint8, 1)
		data = PacketInt(data, math.MaxUint16, 2)
		data = PacketInt(data, math.MaxUint32, 4)
		data = PacketInt(data, math.MaxInt64, 8)
		b.StartTimer()

		v1, data := UnpackInt(data, 1)
		v2, data := UnpackInt(data, 2)
		v3, data := UnpackInt(data, 4)
		v4, data := UnpackInt(data, 8)
		if v1 != math.MaxUint8 || v2 != math.MaxUint16 || v3 != math.MaxUint32 || v4 != math.MaxInt64 {
			panic("BenchmarkUnpackInt fail!")
		}
	}
}

func BenchmarkUnpackString(b *testing.B) {
	str := strings.Repeat("A", 64)
	l := len(str)
	data := make([]byte, 0)
	data = PacketString(data, str, l)
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		str2, _ := UnpackString(data, l)
		if str2 != str {
			panic("BenchmarkUnpackString fail!")
		}
	}
}

func BenchmarkUnpackBytes(b *testing.B) {
	str := []byte(strings.Repeat("A", 64))
	l := len(str)
	data := make([]byte, 0)
	data = PacketBytes(data, str, l)
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		str2, _ := UnpackBytes(data, l)
		if string(str2) != string(str) {
			panic("BenchmarkUnpackString fail!")
		}
	}
}
