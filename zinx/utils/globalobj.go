package utils

import (
	"encoding/json"
	"os"
	"zgame/zinx/ziface"
)

type GlobalObj struct {
	// Server
	TcpServer ziface.IServer // 当前Zinx的全局Server对象
	Host      string         // 当前服务器主机监听的IP
	TcpPort   int            // 当前服务器主机监听的端口号
	Name      string         // 当前服务器的名称

	// Zinx
	Version          string // 当前Zinx的版本号
	MaxPackageSize   uint32 // 当前框架数据包的最大值
	MaxConn          int    // 当前服务器主机允许的最大连接数
	WorkerPoolSize   uint32 // 业务工作Worker池的worker数量
	MaxWorkerTaskLen uint32 // 业务工作Worker对应负责的任务队列的任务数量最大值
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
		Name:             "ZinxServerApp",
		Version:          "V0.5",
		TcpPort:          7777,
		Host:             "0.0.0.0",
		MaxPackageSize:   512,
		MaxConn:          120,
		WorkerPoolSize:   10,
		MaxWorkerTaskLen: 1024,
	}

	// 应该尝试从conf/zinx.json去加载一些用户自定义的参数
	GlobalObject.Reload()
}
