package main

import "zgame/zinx/znet"

func main() {
	s := znet.NewServer("zinx v0.1")
	s.Serve()
}
