package common

import (
	"runtime"
	debug "server/common/debug"
)

func GetGoroutineID() string {
	idStr := debug.GoroutineID()
	return idStr
}

func TraceMemory() *runtime.MemStats {
	var ms runtime.MemStats
	runtime.ReadMemStats(&ms)
	LogDebug("[TraceMemory] Alloc:%d(bytes) HeapIdle:%d(bytes) HeapReleased:%d(bytes)", ms.Alloc, ms.HeapIdle, ms.HeapReleased)
	return &ms
}

func TraceStack() string {
	msg := debug.StackTrace(0).String("    ")
	LogDebug("[TraceStack]:\n%s", msg)
	return msg
}

func GC() {
	LogDebug("Before GC:")
	TraceMemory()
	runtime.GC() // 调用强制gc函数
	LogDebug("After GC:")
	TraceMemory()
}
