package main

import (
	"fmt"
	"io"
	"net"
	"time"
)

func main() {
	// 1. 链接远程服务器，得到一个conn链接
	conn, err := net.Dial("tcp", "127.0.0.1:8999")
	if err != nil {
		fmt.Println("client start err, exit!")
		return
	}

	for {
		// 2. 连接调用Write写数据
		_, err := conn.Write([]byte("hello zinx v0.3"))
		if err != nil {
			fmt.Println("write conn err", err)
			return
		}
		// 3. 服务器端回复的消息
		buf := make([]byte, 512)
		cnt, err := conn.Read(buf)
		if err != nil {
			if err == io.EOF {
				fmt.Println("server close")
				return
			}
			fmt.Println("read buf err", err)
			return
		}
		fmt.Printf("server call back: %s, cnt = %d\n", buf, cnt)
		// cpu阻塞
		time.Sleep(1 * time.Second)
	}

}
