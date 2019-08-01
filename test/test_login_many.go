package main

/*
模拟大量客户端登陆并发包
*/

import (
	"math"
	. "server/common"
	db "server/db/mongodb"
	. "server/msg/protocol"
	. "server/test/common"
	"time"
)

const (
	constLoginNum = 1000
)

func TestLoginMany() {
	db.InitDB()
	ch := make(chan int, constLoginNum)
	for i := 0; i < constLoginNum; i++ {
		ch <- i
	}
	for {
		i := <-ch
		go func(i int) {
			user := FormatString("User_%d", i)
			password := FormatString("Password_%d", i)
			defer func() {
				r := recover()
				if r != nil {
					LogError("%s(%s) error: %v", user, password, r)
				}
				ch <- i
			}()
			c := NewLoginClient()
			pid := c.Login(user, password)
			LogInfo("%s(%s) Login Success: %d", user, password, pid)
			client := c.Client
			responseList := make([]int, 0)
			lastStat := GetMsSecond()
			for {
				v1, v2, v3, v4, v5, v6 := i, i+1, i*2, i*4, user, []byte{'a', 'b', 'c', 'b', 'a'}
				v1 = v1 % math.MaxUint8
				v2 = v2 % math.MaxUint16

				s1 := GetMsSecond()
				client.C2GSEcho(v1, v2, v3, v4, v5, v6)
				p0 := client.Protocol.(*TestEcho)
				client.GS2CEcho()
				e1 := GetMsSecond()
				p := client.Protocol.(*TestEcho)
				if (p.Int1 != v1) || (p.Int2 != v2) || (p.Int3 != v3) || (p.Int4 != v4) || (p.Str != v5) || (string(p.Byte) != string(v6)) {
					LogPanic("Player %d ClientEcho Fail! p: %v p0: %v", pid, p, p0)
				}
				responseList = append(responseList, e1-s1)
				if len(responseList) >= 100 {
					sum, max, min := 0, 0, 0
					for _, v := range responseList {
						sum += v
						if v > max {
							max = v
						}
						if min == 0 || v < min {
							min = v
						}
					}
					avg := sum / len(responseList)
					costTime := GetMsSecond() - lastStat
					lastStat = GetMsSecond()
					LogInfo("%s(%s) %d response total cost: %d MS , avg: %d MS, max: %d MS, min: %d MS", user, password, pid, costTime, avg, max, min)
					responseList = make([]int, 0)
				}
			}
		}(i)
	}

	time.Sleep(30 * time.Second)
}

func main() {
	TestLoginMany()
}
