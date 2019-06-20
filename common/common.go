package common

import (
	"server/conf"
)

func CheckFatal(err error) {
	if err != nil {
		LogFatal("%v", err)
	}
}

func CheckPanic(err error) {
	if err != nil {
		LogPanic(err.Error())
	}
}

func GetServerNum() int {
	return conf.Server.ServerNum
}

func NewObjectID() int {
	v := <-g_ObjectIDChan
	return v
}

func IsPlayerID(id int) bool {
	return id >= PlayerIDMin && id <= PlayerIDMax
}

func init() {
	g_ObjectIDChan = make(chan int, 32)
	go func() {
		objectID := ObjectIDMin
		for {
			objectID++
			g_ObjectIDChan <- objectID
		}
	}()
}
