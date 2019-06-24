package world

const (
	heartbeatInternal = 3600 * 12 //服务器心跳间隔
)

type World struct {
	id int
}

var g_World *World
