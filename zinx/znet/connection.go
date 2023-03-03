package znet

import (
	"errors"
	"fmt"
	"io"
	"net"
	"zgame/zinx/utils"
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

	// 告知当前链接已经退出/停止的channel,由Reader告知Writer退出
	ExitChan chan bool

	// 无缓冲的管道，用于读、写Goroutine之间的消息通信
	MsgChan chan []byte

	// 消息的管理MsgID和对应的处理业务API关系
	MsgHandler ziface.IMsgHandle
}

// 初始化链接模块的方法
func NewConnection(conn *net.TCPConn, connID uint32, msgHandler ziface.IMsgHandle) *Connection {
	c := &Connection{
		Conn:       conn,
		ConnID:     connID,
		MsgHandler: msgHandler,
		isClosed:   false,
		MsgChan:    make(chan []byte),
		ExitChan:   make(chan bool, 1),
	}
	return c
}

// 连接的读业务方法
func (c *Connection) StartReader() {
	fmt.Println("[Reader Goroutine is running...]")
	defer fmt.Println("connID = ", c.ConnID, " [Reader is exit], remote addr is ", c.RemoteAddr().String())
	defer c.Stop()
	for {
		// 读取客户端的数据到buf中
		// buf := make([]byte, utils.GlobalObject.MaxPackageSize)
		// _, err := c.Conn.Read(buf) // 从conn中读取数据
		// if err != nil {
		// 	fmt.Println("recv buf err", err)
		// 	continue
		// }

		// 创建拆包解包对象
		dp := NewDataPack()
		// 读取客户端的Msg Head 二进制流 8个字节
		headData := make([]byte, dp.GetHeadLen())
		_, err := io.ReadFull(c.GetTCPConnection(), headData) // ReadFull 会把msg填充满为止
		if err != nil {
			fmt.Println("read msg head error", err)
			break
		}

		// 拆包，得到msgID和msgDataLen放在msg消息中
		msg, err := dp.Unpack(headData)
		if err != nil {
			fmt.Println("unpack error", err)
			break
		}

		// 根据dataLen再次读取Data，放在msg.Data中
		var data []byte
		if msg.GetDataLen() > 0 {
			data = make([]byte, msg.GetDataLen())
			_, err := io.ReadFull(c.GetTCPConnection(), data)
			if err != nil {
				fmt.Println("read msg data error", err)
				break
			}
		}
		msg.SetData(data)

		// 得到当前conn数据的Request请求数据
		req := Request{
			conn: c,
			msg:  msg,
		}

		if utils.GlobalObject.WorkerPoolSize > 0 {
			// 已经开启工作池机制，将消息交给Worker工作池处理
			c.MsgHandler.SendMsgToTaskQueue(&req)
		} else {

			// 从路由中，找到注册绑定的Conn对应的router调用
			// 根据绑定好的MsgID找到对应处理api业务 执行
			go c.MsgHandler.DoMsgHandler(&req)
		}

	}
}

// 写消息Goroutine，专门发送给客户端消息的模块
func (c *Connection) StartWriter() {
	fmt.Println("[Writer Goroutine is running]")
	defer fmt.Println(c.RemoteAddr().String(), "[conn Writer exit!]")
	// 不断的阻塞等待channel的消息，进行写给客户端
	for {
		select {
		// 如果channel有数据，就发送给客户端
		case data := <-c.MsgChan:
			if _, err := c.Conn.Write(data); err != nil {
				fmt.Println("Send data error", err)
				return
			}
		// 如果channel被关闭，就退出
		case <-c.ExitChan:
			// 代表Reader已经退出，此时Writer也要退出
			return
		}
	}
}

// 启动连接，让当前的连接开始工作
func (c *Connection) Start() {
	fmt.Println("Conn Start()... ConnID = ", c.ConnID)
	// 启动从当前链接的读数据业务
	go c.StartReader()
	// 启动从当前链接写数据业务
	go c.StartWriter()
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
	close(c.MsgChan)
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
func (c *Connection) SendMsg(msgId uint32, data []byte) error {
	if c.isClosed {
		return errors.New("Connection closed when send msg")
	}

	// 将data进行封包 MsgDataLen | MsgID | Data
	dp := NewDataPack()
	binaryMsg, err := dp.Pack(NewMsgPackage(msgId, data))
	if err != nil {
		fmt.Println("Pack error msg id = ", msgId)
		return errors.New("Pack error msg")
	}

	// 将数据发送给客户端,将数据发送给Writer的channel中,由Writer进行写给客户端
	c.MsgChan <- binaryMsg
	return nil
}
