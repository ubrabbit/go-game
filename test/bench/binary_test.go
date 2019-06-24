package main

//go test -bench=. -benchmem -benchtime="1s"

import (
	"bytes"
	"encoding/binary"
	"math"
	"testing"
)

import (
	. "server/leaf/network"
)

func BenchmarkReadBinaryUint8(b *testing.B) {
	var err error
	array := make([]byte, 1)
	array[0] = uint8(math.MaxUint8)
	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		buf := bytes.NewBuffer(array)
		var value uint8
		if LittleEndian {
			err = binary.Read(buf, binary.LittleEndian, &value)
		} else {
			err = binary.Read(buf, binary.BigEndian, &value)
		}
		if err != nil || value != math.MaxUint8 {
			panic("BenchmarkBinaryUint8 fail!")
		}
	}
}

func BenchmarkReadBinaryUint16(b *testing.B) {
	var err error
	array := make([]byte, 2)
	if LittleEndian {
		binary.LittleEndian.PutUint16(array, math.MaxUint16)
	} else {
		binary.BigEndian.PutUint16(array, math.MaxUint16)
	}
	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		buf := bytes.NewBuffer(array)
		var value uint16
		if LittleEndian {
			err = binary.Read(buf, binary.LittleEndian, &value)
		} else {
			err = binary.Read(buf, binary.BigEndian, &value)
		}
		if err != nil || value != math.MaxUint16 {
			panic("BenchmarkBinaryUint16 fail!")
		}
	}
}

func BenchmarkReadBinaryUint32(b *testing.B) {
	var err error
	array := make([]byte, 4)
	if LittleEndian {
		binary.LittleEndian.PutUint32(array, math.MaxUint32)
	} else {
		binary.BigEndian.PutUint32(array, math.MaxUint32)
	}
	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		buf := bytes.NewBuffer(array)
		var value uint32
		if LittleEndian {
			err = binary.Read(buf, binary.LittleEndian, &value)
		} else {
			err = binary.Read(buf, binary.BigEndian, &value)
		}
		if err != nil || value != math.MaxUint32 {
			panic("BenchmarkBinaryUint32 fail!")
		}
	}
}

func BenchmarkReadBinaryUint64(b *testing.B) {
	var err error
	array := make([]byte, 8)
	if LittleEndian {
		binary.LittleEndian.PutUint64(array, math.MaxUint64)
	} else {
		binary.BigEndian.PutUint64(array, math.MaxUint64)
	}
	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		buf := bytes.NewBuffer(array)
		var value uint64
		if LittleEndian {
			err = binary.Read(buf, binary.LittleEndian, &value)
		} else {
			err = binary.Read(buf, binary.BigEndian, &value)
		}
		if err != nil || value != math.MaxUint64 {
			panic("BenchmarkBinaryUint64 fail!")
		}
	}
}
