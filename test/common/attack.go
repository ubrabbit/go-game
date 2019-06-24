package common

import (
	"time"
)

import (
	"server/test/common/synflood"
)

//Synflood Attack!
func AttackSynflood(timeout int) {
	for i := 0; i < constSyncfloodCount; i++ {
		go func() {
			synflood.Attack("127.0.0.1", 38320)
		}()
	}
	time.Sleep(time.Duration(timeout) * time.Second)
}
