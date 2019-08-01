package main

import (
	"server/conf"
	"server/db"
	"server/game"
	"server/gate"
	"server/leaf"
	lconf "server/leaf/conf"
	"server/login"
	"server/scene"
)

import (
	"fmt"
	"net/http"
	_ "net/http/pprof"
)

import (
	. "server/common"
)

func main() {
	lconf.LogLevel = conf.Server.LogLevel
	lconf.LogPath = conf.Server.LogPath
	lconf.LogFlag = conf.LogFlag
	lconf.ConsolePort = conf.Server.ConsolePort
	lconf.ProfilePath = conf.Server.ProfilePath

	go func() {
		addr := fmt.Sprintf("localhost:%d", conf.Server.PprofPort)
		err := http.ListenAndServe(addr, nil)
		if err != nil {
			LogError("pprof listen error: %v", err)
		}
	}()

	leaf.Run(
		db.Module,
		gate.Module,
		game.Module,
		scene.Module,
		login.Module,
	)
}
