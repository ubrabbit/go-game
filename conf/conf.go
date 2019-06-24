package conf

import (
	"log"
	"time"
)

var (
	// log conf
	LogFlag = log.LstdFlags

	// gate conf
	PendingWriteNum        = 8192
	MaxMsgLen       uint32 = 40960
	HTTPTimeout            = 10 * time.Second
	LenMsgLen              = 2

	// skeleton conf
	GoLen              = 20480
	TimerDispatcherLen = 10000
	AsynCallLen        = 10000
	ChanRPCLen         = 10000
)
