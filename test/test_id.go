package main

import (
	"sync/atomic"
	"time"
)

import (
	. "server/common"
	db "server/db/mongodb"
)

var g_Finished bool
var g_Count int64

func test_objectid() {
	g_Finished = false
	atomic.StoreInt64(&g_Count, 0)
	for {
		if g_Finished {
			break
		}
		NewObjectID()
		atomic.AddInt64(&g_Count, 1)
	}
}

func test_playerid() {
	db.InitDB()
	g_Finished = false
	atomic.StoreInt64(&g_Count, 0)
	for {
		if g_Finished {
			break
		}
		db.NewPlayerID()
		atomic.AddInt64(&g_Count, 1)
	}
}

func main() {
	//for i := 0; i < 10; i++ {
	//	go test_objectid()
	//}
	go test_objectid()
	time.Sleep(1 * time.Second)
	g_Finished = true
	LogInfo("test_objectid: %d", g_Count)

	go test_playerid()
	time.Sleep(1 * time.Second)
	g_Finished = true
	LogInfo("test_playerid: %d", g_Count)
}
