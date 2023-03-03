package main

import (
	"fmt"
	"io"
	"net"
	"time"
	"zgame/zinx/znet"
)

func main() {
	// 1. 链接远程服务器，得到一个conn链接
	conn, err := net.Dial("tcp", "127.0.0.1:7777")
	if err != nil {
		fmt.Println("client start err, exit!")
		return
	}

	for {
		// 发送封包的msg消息, 0号消息
		dp := znet.NewDataPack()
		binaryMsg, err := dp.Pack(znet.NewMsgPackage(0, []byte("zinx0.5 client test message")))
		if err != nil {
			fmt.Println("Pack error msg id = ", 0)
			return
		}
		if _, err := conn.Write(binaryMsg); err != nil {
			fmt.Println("write error err", err)
			return
		}

		// 服务器应该先回复一个MsgID:0的消息
		// 1. 先读出流中的head部分，得到ID和dataLen
		headData := make([]byte, dp.GetHeadLen())
		if _, err := io.ReadFull(conn, headData); err != nil {
			fmt.Println("read head error")
			break
		}
		// 将headData字节流 拆包到msg中
		msgHead, err := dp.Unpack(headData)
		if err != nil {
			fmt.Println("client unpack err", err)
			return
		}
		if msgHead.GetDataLen() > 0 {
			// msg 是有data数据的，需要再次读取data数据
			msg := msgHead.(*znet.Message)
			msg.Data = make([]byte, msg.GetDataLen())

			// 根据dataLen从io中读取字节流
			if _, err := io.ReadFull(conn, msg.Data); err != nil {
				fmt.Println("read msg data error", err)
				return
			}
			fmt.Println("---> Recv Server Msg: ID=", msg.Id, ", len=", msg.DataLen, ", data=", string(msg.Data))
		}

		// cpu阻塞
		time.Sleep(1 * time.Second)
	}

}
