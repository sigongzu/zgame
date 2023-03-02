package znet

import (
	"fmt"
	"net"
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
	// 当前Server添加一个router，server注册的链接对应的处理业务
	Router ziface.IRouter
}

// AddRouter 给当前的Server注册一个路由方法，供客户端的链接处理使用
func (s *Server) AddRouter(router ziface.IRouter) {
	s.Router = router
	fmt.Println("Add Router Succ!")
}

// NewServer 初始化Server模块的方法
func NewServer(name string) ziface.IServer {
	s := &Server{
		Name:      name,
		IPVersion: "tcp4",
		IP:        "0.0.0.0",
		Port:      8999,
		Router:    nil,
	}
	return s
}

// Start 启动服务器
func (s *Server) Start() {
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
		// 已经与客户端建立链接
		dealConn := NewConnection(conn, cid, s.Router)
		cid++

		// 启动当前链接的处理业务
		go dealConn.Start()
	}

}

// Stop 停止服务器
func (s *Server) Stop() {
	// TODO 将一些服务器的资源、状态或者已经开辟的链接信息 进行停止或者回收
}

// Serve 运行服务器
func (s *Server) Serve() {
	// 启动server的服务功能
	s.Start()

	// TODO 做一些启动服务器之后的额外业务

	// 阻塞状态
	select {}
}
