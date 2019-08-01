package main

import (
	. "server/common"
	. "server/test/common"
	"time"
)

var loginUsers = map[string]string{
	"demo_1": "123456",
	"demo_2": "123456",
	"demo_3": "123456",
	"xyc":    "xyc",
	"lpx":    "lpx",
}

func TestLogin() {
	for user, password := range loginUsers {
		c := NewLoginClient()
		pid := c.Login(user, password)
		LogInfo("%s(%s) login success! pid=%d", user, password, pid)
	}
	time.Sleep(1 * time.Second)
}

func TestLoginFail1() {
	for user, password := range loginUsers {
		c := NewLoginClient()
		pid := c.LoginNoHello(user, password)
		if pid != 0 {
			LogPanic("%s(%s) TestLoginFail1 fail! pid=%d", user, password, pid)
		} else {
			LogInfo("%s(%s) TestLoginFail1 success! pid=%d", user, password, pid)
		}
	}
	time.Sleep(1 * time.Second)
}

func TestLoginFail2() {
	for user, password := range loginUsers {
		c := NewLoginClient()
		pid := c.LoginNoIdentity(user, password)
		if pid != 0 {
			LogPanic("%s(%s) TestLoginFail2 fail! pid=%d", user, password, pid)
		} else {
			LogInfo("%s(%s) TestLoginFail2 success! pid=%d", user, password, pid)
		}
	}
	time.Sleep(1 * time.Second)
}

func TestLoginNoVerify() {
	for user, password := range loginUsers {
		c := NewLoginClient()
		c.LoginNoVerify(user, password)
	}
	time.Sleep(1 * time.Second)
}

func TestLoginNoAuth() {
	for user, password := range loginUsers {
		c := NewLoginClient()
		c.LoginNoAuth(user, password)
	}
	time.Sleep(1 * time.Second)
}

func main() {
	TestLogin()
	//TestLoginFail1()
	//TestLoginFail2()
	//TestLoginNoVerify()
	//TestLoginNoAuth()
}
