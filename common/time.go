package common

import (
	"time"
)

func GetTimeString() string {
	var currentTime time.Time

	currentTime = time.Now().Local()
	newFormat := currentTime.Format("2006-01-02 15:04:05")
	return newFormat
}

func GetSecond() int {
	now := time.Now()
	return int(now.Unix())
}

//毫秒
func GetMsSecond() int {
	now := time.Now()
	//把纳秒转换成毫秒
	return int(now.UnixNano() / 1000000)
}

//纳秒
func GetNaSecond() int64 {
	now := time.Now()
	//把纳秒转换成毫秒
	return now.UnixNano()
}

func CreateTimer(ms int) chan bool {
	ch := make(chan bool)
	go func() {
		time.Sleep(time.Millisecond * time.Duration(ms))
		ch <- true
	}()
	return ch
}
