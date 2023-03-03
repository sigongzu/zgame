package main

import (
	"fmt"
	"zgame/zinx/ziface"
	"zgame/zinx/znet"
)

// ping test 自定义路由
type PingRouter struct {
	znet.BaseRouter
}

// Test Handle
func (p *PingRouter) Handle(request ziface.IRequest) {
	fmt.Println("Call pingRouter Handle...")
	// 先读取客户端的数据，再ping...ping...ping
	fmt.Println("recv from client: msgID=", request.GetMsgID(), ", data=", string(request.GetData()))

	err := request.GetConnection().SendMsg(200, []byte("ping...ping...ping"))
	if err != nil {
		fmt.Println(err)
	}
}

// hello zinx test 自定义路由
type HelloZinxRouter struct {
	znet.BaseRouter
}

// HelloZinxRouter Handle
func (p *HelloZinxRouter) Handle(request ziface.IRequest) {
	fmt.Println("Call HelloZinxRouter Handle...")
	// 先读取客户端的数据，再ping...ping...ping
	fmt.Println("recv from client: msgID=", request.GetMsgID(), ", data=", string(request.GetData()))

	err := request.GetConnection().SendMsg(201, []byte("HelloZinxRouter..............."))
	if err != nil {
		fmt.Println(err)
	}
}

func main() {
	s := znet.NewServer("zinx v0.5")
	s.AddRouter(0, &PingRouter{})
	s.AddRouter(1, &HelloZinxRouter{})
	s.Serve()
}
