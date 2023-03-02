package utils

import (
	"encoding/json"
	"os"
	"zgame/zinx/ziface"
)

type GlobalObj struct {
	// Server
	TcpServer ziface.IServer
	Host      string
	TcpPort   int
	Name      string
	// Version
	Version        string
	MaxPackageSize uint32
	MaxConn        int
}

// 定义一个全局的对外GlobalObj
var GlobalObject *GlobalObj

// 加载用户配置的文件 zinx.json
func (g *GlobalObj) Reload() {
	data, err := os.ReadFile("conf/zinx.json")
	if err != nil {
		panic(err)
	}
	// 将json数据解析到struct中
	err = json.Unmarshal(data, &GlobalObject)
	if err != nil {
		panic(err)
	}
}

// 提供一个init方法，初始化当前的GlobalObject变量
func init() {
	// 如果配置文件没有加载，默认的值
	GlobalObject = &GlobalObj{
		Name:           "ZinxServerApp",
		Version:        "V0.4",
		TcpPort:        8999,
		Host:           "0.0.0.0",
		MaxPackageSize: 512,
		MaxConn:        120,
	}

	// 应该尝试从conf/zinx.json去加载一些用户自定义的参数
	GlobalObject.Reload()
}
