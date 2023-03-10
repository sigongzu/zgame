package znet

import (
	"fmt"
	"net"
	"zgame/zinx/utils"
	"zgame/zinx/ziface"
)

// IServer的接口实现，定义一个Server的服务器模块
type Server struct {
	// 服务器的名称
	Name string
	// 服务器绑定的IP版本
	IPVersion string
	// 服务器监听的IP
	IP string
	// 服务器监听的端口
	Port int
	// 当前Server的消息管理模块，用来绑定MsgID和对应的处理业务API关系
	MsgHandler ziface.IMsgHandle
	// 该Server的连接管理器
	ConnMgr ziface.IConnManager
	// 该Server的连接创建时Hook函数
	OnConnStart func(conn ziface.IConnection)
	// 该Server的连接断开时的Hook函数
	OnConnStop func(conn ziface.IConnection)
}

// NewServer 初始化Server模块的方法
func NewServer(name string) ziface.IServer {
	s := &Server{
		Name:       utils.GlobalObject.Name,
		IPVersion:  "tcp4",
		IP:         utils.GlobalObject.Host,
		Port:       utils.GlobalObject.TcpPort,
		MsgHandler: NewMsgHandle(),
		ConnMgr:    NewConnManager(),
	}
	return s
}

// SetOnConnStart 注册OnConnStart钩子函数的方法
func (s *Server) SetOnConnStart(hookFunc func(connection ziface.IConnection)) {
	s.OnConnStart = hookFunc
}

// SetOnConnStop 注册OnConnStop钩子函数的方法
func (s *Server) SetOnConnStop(hookFunc func(connection ziface.IConnection)) {
	s.OnConnStop = hookFunc
}

// CallOnConnStart 调用OnConnStart钩子函数的方法
func (s *Server) CallOnConnStart(conn ziface.IConnection) {
	if s.OnConnStart != nil {
		fmt.Println("Call OnConnStart()...")
		s.OnConnStart(conn)
	}
}

// CallOnConnStop 调用OnConnStop钩子函数的方法
func (s *Server) CallOnConnStop(conn ziface.IConnection) {
	if s.OnConnStop != nil {
		fmt.Println("Call OnConnStop()...")
		s.OnConnStop(conn)
	}
}

// AddRouter 给当前的Server注册一个路由方法，供客户端的链接处理使用
func (s *Server) AddRouter(msgID uint32, router ziface.IRouter) {
	s.MsgHandler.AddRouter(msgID, router)
	fmt.Println("Add Router Succ!")
}

// GetConnMgr 得到链接管理
func (s *Server) GetConnMgr() ziface.IConnManager {
	return s.ConnMgr
}

// Start 启动服务器
func (s *Server) Start() {
	fmt.Printf("Server Name: %s, listenner at IP: %s, Port: %d is starting\n",
		utils.GlobalObject.Name, utils.GlobalObject.Host, utils.GlobalObject.TcpPort)
	fmt.Printf("Version: %s, MaxConn: %d, MaxPackageSize: %d\n",
		utils.GlobalObject.Version, utils.GlobalObject.MaxConn, utils.GlobalObject.MaxPackageSize)

	go func() {
		// 0 开启消息队列及Worker工作池
		s.MsgHandler.StartWorkerPool()
		// 1 获取一个TCP的Addr
		addr, err := net.ResolveTCPAddr(s.IPVersion, fmt.Sprintf("%s:%d", s.IP, s.Port))
		if err != nil {
			fmt.Println("resolve tcp addr error:", err)
			return
		}
		// 2 监听服务器的地址
		listener, err := net.ListenTCP(s.IPVersion, addr)
		if err != nil {
			fmt.Println("listen", s.IPVersion, "err", err)
			return
		}
		fmt.Println("start Zinx server succ, ", s.Name, " succ, now listening...")

		var cid uint32 = 0
		// 3 阻塞的等待客户端链接，处理客户端链接业务（读写）
		for {
			// 如果有客户端链接过来，阻塞会返回
			conn, err := listener.AcceptTCP()
			if err != nil {
				fmt.Println("Accept err", err)
				continue
			}

			// 设置最大连接个数的判断，如果超过最大连接，那么则关闭此新的连接
			if s.ConnMgr.Len() >= utils.GlobalObject.MaxConn {
				// TODO 给客户端响应一个超出最大连接的错误包
				fmt.Println("Too Many Connections MaxConn ====================== ", utils.GlobalObject.MaxConn)
				conn.Close()
				continue
			}

			// 已经与客户端建立链接
			dealConn := NewConnection(s, conn, cid, s.MsgHandler)
			cid++

			// 启动当前链接的处理业务
			go dealConn.Start()
		}
	}()
}

// Stop 停止服务器
func (s *Server) Stop() {
	// TODO 将一些服务器的资源、状态或者已经开辟的链接信息 进行停止或者回收
	fmt.Println("[STOP] Zinx server name ", s.Name)
	s.ConnMgr.ClearConn()
}

// Serve 运行服务器
func (s *Server) Serve() {
	// 启动server的服务功能
	s.Start()

	// TODO 做一些启动服务器之后的额外业务

	// 阻塞状态
	select {}
}
