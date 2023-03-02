package znet

import (
	"fmt"
	"net"
	"zgame/zinx/ziface"
)

// 链接模块
type Connection struct {
	// 当前链接的socket TCP套接字
	Conn *net.TCPConn
	// 当前链接的ID，也可以称为sessionID，全局唯一
	ConnID uint32
	// 当前链接的状态
	isClosed bool
	// 当前链接所绑定的处理业务方法API
	handleAPI ziface.HandleFunc
	// 告知当前链接已经退出/停止的channel
	ExitChan chan bool
}

// 初始化链接模块的方法
func NewConnection(conn *net.TCPConn, connID uint32, callbackAPI ziface.HandleFunc) *Connection {
	c := &Connection{
		Conn:      conn,
		ConnID:    connID,
		handleAPI: callbackAPI,
		isClosed:  false,
		ExitChan:  make(chan bool, 1),
	}
	return c
}

// 连接的读业务方法
func (c *Connection) StartReader() {
	fmt.Println("Reader Goroutine is running...")
	defer fmt.Println("connID = ", c.ConnID, " Reader is exit, remote addr is ", c.RemoteAddr().String())
	defer c.Stop()
	for {
		// 读取客户端的数据到buf中，最大512字节
		buf := make([]byte, 512)
		cnt, err := c.Conn.Read(buf) // 从conn中读取数据
		if err != nil {
			fmt.Println("recv buf err", err)
			continue
		}
		// 调用当前链接所绑定的HandleAPI
		if err := c.handleAPI(c.Conn, buf, cnt); err != nil {
			fmt.Println("ConnID", c.ConnID, "handle is error", err)
			break
		}
	}
}

// 启动连接，让当前的连接开始工作
func (c *Connection) Start() {
	fmt.Println("Conn Start()... ConnID = ", c.ConnID)
	// 启动从当前链接的读数据业务
	go c.StartReader()
	// 启动从当前链接写数据业务
	// go c.StartWriter()
}

// 停止连接，结束当前连接状态M
func (c *Connection) Stop() {
	fmt.Println("Conn Stop()... ConnID = ", c.ConnID)
	// 如果当前链接已经关闭
	if c.isClosed {
		return
	}
	c.isClosed = true
	// 关闭socket链接
	c.Conn.Close()
	// 告知Writer关闭
	c.ExitChan <- true
	// 关闭该链接全部管道
	close(c.ExitChan)
}

// 获取当前连接绑定的socket conn
func (c *Connection) GetTCPConnection() *net.TCPConn {
	return c.Conn
}

// 获取当前连接模块的连接ID
func (c *Connection) GetConnID() uint32 {
	return c.ConnID
}

// 获取远程客户端的TCP状态 IP port
func (c *Connection) RemoteAddr() net.Addr {
	return c.Conn.RemoteAddr()
}

// 发送数据，将数据发送给远程的客户端
func (c *Connection) Send(data []byte) error {
	return nil
}
