package main

import (
	. "server/test/common"
)

const (
	constClientCount = 100
	constLoopCount   = 100
)

func main() {
	c := NewClient()
	c.TestEcho(constClientCount, constLoopCount)
}
