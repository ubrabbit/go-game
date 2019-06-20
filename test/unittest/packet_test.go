package unittest

import (
	"math"
	"strings"
	"testing"
)

import (
	. "server/common"
	. "server/leaf/network"
)

func TestPacketInt(t *testing.T) {
	LogInfo("math.MaxUint8 == %d", math.MaxUint8)
	LogInfo("math.MaxUint16 == %d", math.MaxUint16)
	LogInfo("math.MaxUint32 == %d", math.MaxUint32)

	data := make([]byte, 0)
	data = PacketInt(data, math.MaxUint8, 1)
	data = PacketInt(data, math.MaxUint16, 2)
	data = PacketInt(data, math.MaxUint32, 4)

	v, data := UnpackInt(data, 1)
	if v != math.MaxUint8 {
		t.Error("TestPacketInt Fail !")
	}
	v, data = UnpackInt(data, 2)
	if v != math.MaxUint16 {
		t.Error("TestPacketInt Fail !")
	}
	v, data = UnpackInt(data, 4)
	if v != math.MaxUint32 {
		t.Error("TestPacketInt Fail !")
	}
	if len(data) != 0 {
		t.Error("TestPacketInt Fail !")
	}

	value_list := []int{}
	for i := 0; i <= math.MaxInt8; i++ {
		value_list = append(value_list, i)
		data = PacketInt(data, i, 1)
	}
	for _, i := range value_list {
		v, data = UnpackInt(data, 1)
		if v != i {
			t.Error("PacketInt(1) Fail !")
		}
	}
	if len(data) != 0 {
		t.Error("PacketInt(1) Fail !")
	}

	int16_list := []int{0, 65534, 65535}
	for _, i := range int16_list {
		value_list = append(value_list, i)
		data = PacketInt(data, i, 2)
		v, data = UnpackInt(data, 2)
		if v != i {
			t.Error("PacketInt(2) Fail !")
		}
		if len(data) != 0 {
			t.Error("PacketInt(2) Fail !")
		}
	}

	data = PacketInt(data, math.MaxUint8+1, 1)
	data = PacketInt(data, math.MaxUint16+2, 2)
	data = PacketInt(data, math.MaxUint32+3, 4)
	v, data = UnpackInt(data, 1)
	if v != math.MaxUint8 {
		t.Error("PacketInt(overflow) Fail !")
	}
	v, data = UnpackInt(data, 2)
	if v != math.MaxUint16 {
		t.Error("PacketInt(overflow) Fail !")
	}
	v, data = UnpackInt(data, 4)
	if v != math.MaxUint32 {
		t.Error("PacketInt(overflow) Fail !")
	}

	int32_list := []int{0, 127, 128, 255, 256, 65534, 65535, 65536, 2147483647, 2147483646, 2147483648}
	for _, i := range int32_list {
		value_list = append(value_list, i)
		data = PacketInt(data, i, 4)

	}
	for _, i := range int32_list {
		v, data = UnpackInt(data, 4)
		if v != i {
			t.Error("PacketInt(4) Fail !")
		}
	}
	if len(data) != 0 {
		t.Error("PacketInt(4) Fail !")
	}
}

func TestPacketString(t *testing.T) {
	data := make([]byte, 0)
	sizeList := []int{1, 2, 4, 8, 16, 32, 64, 128, 256, 1024, 65535}
	stringList := []string{}
	for _, l := range sizeList {
		str := strings.Repeat("A", l)
		data = PacketString(data, str, len(str))
		stringList = append(stringList, str)
	}
	for _, str := range stringList {
		l := len(str)
		v := ""
		v, data = UnpackString(data, l)
		if v != str {
			t.Error("PacketString() Fail !")
		}
	}
	if len(data) != 0 {
		LogInfo("len(data) == %d", len(data))
		t.Error("PacketString() Fail !")
	}

	//400M
	max := 400 * 1024 * 1024
	str := strings.Repeat("B", max)
	data = PacketString(data, str, max)
	str2, data := UnpackString(data, 0)
	if str != str2 {
		t.Error("PacketString(max) Fail !")
	}
}

func TestPacket(t *testing.T) {
	data := make([]byte, 0)
	proto := 127
	strLen1, strLen2 := 32, 40960
	str := strings.Repeat("A", strLen1)
	str2 := strings.Repeat("B", strLen2)

	data = PacketInt(data, proto, 1)
	data = PacketInt(data, 0, 1)
	data = PacketInt(data, 1, 1)
	data = PacketInt(data, math.MaxUint8, 1)
	data = PacketInt(data, math.MaxUint16, 2)
	data = PacketString(data, str, strLen1)
	data = PacketInt(data, math.MaxUint32, 4)
	data = PacketString(data, str2, strLen2)

	i, data := UnpackInt(data, 1)
	if i != proto {
		t.Error("UnpackInt(proto) Fail !")
	}
	i, data = UnpackInt(data, 1)
	if i != 0 {
		t.Error("UnpackInt(proto) Fail !")
	}
	i, data = UnpackInt(data, 1)
	if i != 1 {
		t.Error("UnpackInt(proto) Fail !")
	}
	i, data = UnpackInt(data, 1)
	if i != math.MaxUint8 {
		t.Error("UnpackInt(1) Fail !")
	}
	i, data = UnpackInt(data, 2)
	if i != math.MaxUint16 {
		t.Error("UnpackInt(2) Fail !")
	}

	str3, data := UnpackString(data, strLen1)
	if str != str3 {
		t.Error("UnpackString() Fail !")
	}
	i, data = UnpackInt(data, 4)
	if i != math.MaxUint32 {
		t.Error("UnpackInt(4) Fail !")
	}
	str4, data := UnpackString(data, strLen2)
	if str2 != str4 {
		t.Error("UnpackString() Fail !")
	}
	if len(data) != 0 {
		t.Error("TestPacket Fail !")
	}
}
