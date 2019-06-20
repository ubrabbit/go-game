package common

import (
	"bytes"
	"encoding/binary"
	"strconv"
)

func StringToInt(str string) int {
	i, err := strconv.Atoi(str)
	CheckPanic(err)
	return i
}

func StringToInt64(str string) int64 {
	/*
	   参数1 数字的字符串形式
	   参数2 数字字符串的进制 比如二进制 八进制 十进制 十六进制
	   参数3 返回结果的bit大小 也就是int8 int16 int32 int64
	*/
	i, err := strconv.ParseInt(str, 0, 64)
	CheckPanic(err)
	return i
}

func StringToFloat64(str string) float64 {
	f, err := strconv.ParseFloat(str, 64)
	CheckPanic(err)
	return f
}

func IntToString(v int) string {
	return strconv.Itoa(v)
}

func BytesToInt(buf []byte) int {
	data := int(binary.BigEndian.Uint32(buf))
	return data
}

func IntToBytes(n int) []byte {
	x := uint32(n)
	//创建一个内容是[]byte的slice的缓冲器
	//与bytes.NewBufferString("")等效
	bytesBuffer := bytes.NewBuffer([]byte{})
	binary.Write(bytesBuffer, binary.BigEndian, x)
	return bytesBuffer.Bytes()
}
