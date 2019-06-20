package main

import (
	. "server/common"
	. "server/test/common"
)

var loginUsers = map[string]string{
	"demo_1": "123456",
	"demo_2": "123456",
	"demo_3": "123456",
	"xyc":    "xyc",
	"lpx":    "lpx",
}

func TestLogin() {
	c := NewLoginClient()
	for user, password := range loginUsers {
		pid := c.Login(user, password)
		LogInfo("%s(%s) login success! pid=%d", user, password, pid)
	}
}

func main() {
	TestLogin()
}
