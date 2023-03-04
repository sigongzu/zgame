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

// 创建连接之后执行的钩子函数
func DoConnectionBegin(conn ziface.IConnection) {
	fmt.Println("DoConnectionBegin is Called ...")
	if err := conn.SendMsg(202, []byte("DoConnection BEGIN")); err != nil {
		fmt.Println(err)
	}

	// 设置链接属性
	fmt.Println("Set Conn Property...")
	conn.SetProperty("Name", "任我行")
	conn.SetProperty("Address", "日月神教")

}

// 销毁连接之前执行的钩子函数
func DoConnectionLost(conn ziface.IConnection) {
	fmt.Println("DoConnectionLost is Called ...")
	if err := conn.SendMsg(203, []byte("DoConnection Lost")); err != nil {
		fmt.Println(err)
	}

	// 获取链接属性
	fmt.Println("Get Conn Property...")
	if name, err := conn.GetProperty("Name"); err == nil {
		fmt.Println("Name = ", name)
	}

	if address, err := conn.GetProperty("Address"); err == nil {
		fmt.Println("Address = ", address)
	}
}

func main() {
	// 1. 创建一个server句柄，使用zinx的api
	s := znet.NewServer("zinx v0.5")

	// 2. 注册链接hook函数
	s.SetOnConnStart(DoConnectionBegin)
	s.SetOnConnStop(DoConnectionLost)

	// 3. 注册路由
	s.AddRouter(0, &PingRouter{})
	s.AddRouter(1, &HelloZinxRouter{})

	// 4. 启动server
	s.Serve()
}
