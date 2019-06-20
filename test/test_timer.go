package main

import (
	. "server/common"
	. "server/test/common"
)

var testSecondList = []int{1, 1, 2, 3, 3, 5}

func main() {
	TestTimerWithChannel(testSecondList)
	TestTimerNoChannel(testSecondList)
	LogInfo("test success!")
}
