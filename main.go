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

func main() {
	lconf.LogLevel = conf.Server.LogLevel
	lconf.LogPath = conf.Server.LogPath
	lconf.LogFlag = conf.LogFlag
	lconf.ConsolePort = conf.Server.ConsolePort
	lconf.ProfilePath = conf.Server.ProfilePath

	leaf.Run(
		db.Module,
		gate.Module,
		game.Module,
		scene.Module,
		login.Module,
	)
}
