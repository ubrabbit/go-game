package unittest

import (
	"testing"
)

import (
	. "server/test/common"
)

var testSecondList = []int{1, 1, 2, 3}

func TestTimer(t *testing.T) {
	TestTimerWithChannel(testSecondList)
	TestTimerNoChannel(testSecondList)
	TestTimerGoroutineSafe()
}
