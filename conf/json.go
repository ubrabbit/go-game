package conf

import (
	"encoding/json"
	"io/ioutil"
	"server/leaf/log"
)

var Server struct {
	LogLevel     string
	LogPath      string
	WSAddr       string
	CertFile     string
	KeyFile      string
	TCPAddr      string
	MaxConnNum   int
	ConsolePort  int
	PprofPort    int
	ProfilePath  string
	ServerNum    int
	DatabaseAddr string
}

func init() {
	data, err := ioutil.ReadFile("conf/server.json")
	if err != nil {
		//找不到，就在当前目录查找是否有测试用配置文件
		data, err = ioutil.ReadFile("server_test.json")
		if err != nil {
			log.Fatal("read server.json failure: '%v'", err)
		} else {
			log.Release("use server_test.json")
		}
	}
	err = json.Unmarshal(data, &Server)
	if err != nil {
		log.Fatal("unmarshal server.json failure: '%v'", err)
	}
}
