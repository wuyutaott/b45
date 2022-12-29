package utils

import (
	"encoding/json"
	"fmt"
	"os"
)

type GlobalObj struct {
	// server
	Host    string
	TCPPort int
	Name    string

	// b45
	Version          string
	MaxPacketSize    uint32
	MaxConn          int
	WorkerPoolSize   uint32
	MaxWorkerTaskLen uint32
	MaxMsgChanLen    uint32

	// log
	LogDir        string
	LogFile       string
	LogDebugClose bool

	// 配置文件路径
	ConfFilePath string
}

var GlobalObject *GlobalObj

func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (g *GlobalObj) Reload() {
	if exist, _ := PathExists(g.ConfFilePath); exist != true {
		fmt.Printf("Config File %s is not exist!! \n", g.ConfFilePath)
		return
	}
	data, err := os.ReadFile(g.ConfFilePath)
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(data, g)
	fmt.Println("globalObj:", g)
	if err != nil {
		panic(err)
	}
}
