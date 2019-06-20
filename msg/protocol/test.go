package protocol

import (
	. "server/leaf/network"
)

type TestEcho struct {
	Int1 int
	Int2 int
	Int3 int
	Int4 int
	Str  string
	Byte []byte
}

func (p *TestEcho) Protocol() uint8 {
	return uint8(TEST_ECHO)
}

func (p *TestEcho) PacketData() (uint8, []byte) {
	data := make([]byte, 0)
	data = PacketInt(data, p.Int1, 1)
	data = PacketInt(data, p.Int2, 2)
	data = PacketInt(data, p.Int3, 4)
	data = PacketInt(data, p.Int4, 8)
	nlen1, nlen2 := len(p.Str), len(p.Byte)
	data = PacketInt(data, nlen1, 2)
	data = PacketString(data, p.Str, nlen1)
	data = PacketInt(data, nlen2, 2)
	data = PacketBytes(data, p.Byte, nlen2)
	return p.Protocol(), data
}

func (p *TestEcho) UnpackData(from []byte) {
	int1, from := UnpackInt(from, 1)
	int2, from := UnpackInt(from, 2)
	int3, from := UnpackInt(from, 4)
	int4, from := UnpackInt(from, 8)
	nlen1, from := UnpackInt(from, 2)
	str, from := UnpackString(from, nlen1)
	nlen2, from := UnpackInt(from, 2)
	b, from := UnpackBytes(from, nlen2)

	p.Int1 = int1
	p.Int2 = int2
	p.Int3 = int3
	p.Int4 = int4
	p.Str = str
	p.Byte = b
}
