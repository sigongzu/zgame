package ziface

import "net"

// IConnection 定义连接模块的抽象层
type IConnection interface {
	// 启动连接，让当前的连接开始工作
	Start()
	// 停止连接，结束当前连接状态M
	Stop()
	// 获取当前连接绑定的socket conn
	GetTCPConnection() *net.TCPConn
	// 获取当前连接模块的连接ID
	GetConnID() uint32
	// 获取远程客户端的TCP状态 IP port
	RemoteAddr() net.Addr
	// 发送数据，将数据发送给远程的客户端
	SendMsg(msgId uint32, data []byte) error
	// 设置链接属性
	SetProperty(key string, value interface{})
	// 获取链接属性
	GetProperty(key string) (interface{}, error)
	// 移除链接属性
	RemoveProperty(key string)
}

// HandleFunc 定义一个处理链接业务的方法
type HandleFunc func(*net.TCPConn, []byte, int) error // 传入的参数：conn, data, len
