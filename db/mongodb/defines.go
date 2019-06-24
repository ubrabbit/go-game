package mongodb

import (
	"gopkg.in/mgo.v2"
	db "server/leaf/db/mongodb"
)

/*
mongodb的文档定义说明
*/

const (
	sessionNum  = 128
	minPlayerID = 10000
	/*
		BSON大小限制。这个暂时也不知道怎么用，先写在代码里用于给自己提醒。
		BSON Documents. The maximum BSON document size is 16 megabytes.
		The maximum document size helps ensure that a single document cannot use excessive amount of RAM or,
		during transmission, excessive amount of bandwidth.
		To store documents larger than the maximum size, MongoDB provides the GridFS API.
	*/
	maxBsonSize = 16 * 1024 * 1024
)

type MongoSession struct {
	Context *db.DialContext
}

//这几个错误变量仅用于减少其他模块对mgo的import
var ErrNotFound = mgo.ErrNotFound
var ErrCursor = mgo.ErrCursor

//数据库名称，一个服务器使用一个
var DatabaseName string
var g_MongoSession *MongoSession
