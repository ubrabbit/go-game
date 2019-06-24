package main

import (
	. "server/common"
)

func test_alloc() []string {
	sliceNum := 10 * 1024 * 1024
	b := make([]string, 0)
	for i := 0; i < sliceNum; i++ {
		b = append(b, "A")
	}
	return b[10:20]
}

func main() {
	GetGoroutineID()
	TraceMemory()
	TraceStack()

	test_alloc()
	GC()
}
