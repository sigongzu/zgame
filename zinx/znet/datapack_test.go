package znet

import (
	"fmt"
	"io"
	"net"
	"testing"
)

// 负责测试datapack拆包，封包的单元测试模块
func TestDataPack(t *testing.T) {
	// 模拟服务器
	// 创建socketTCP
	listener, err := net.Listen("tcp4", "127.0.0.1:7777")
	if err != nil {
		t.Error("server listen err:", err)
		return
	}
	// 创建一个go承载 负责从客户端处理业务
	go func() {
		// 从客户端读取数据，拆包处理
		for {
			conn, err := listener.Accept()
			if err != nil {
				t.Error("server accept err:", err)
				return
			}
			go func(conn net.Conn) {
				// 处理客户端的请求
				// 定义一个拆包的对象
				dp := NewDataPack()
				for {
					// 1 先读出流中的head部分
					headData := make([]byte, dp.GetHeadLen())
					_, err := io.ReadFull(conn, headData) // ReadFull 会把msg填充满为止
					if err != nil {
						t.Error("read head error")
						return
					}
					// 将headData字节流 拆包到msg中
					msgHead, err := dp.Unpack(headData)
					if err != nil {
						t.Error("server unpack err:", err)
						return
					}
					if msgHead.GetDataLen() > 0 {
						// msg 是有data数据的，需要再次读取data数据
						// 2 再根据dataLen从io流中读取字节流，放在msg.Data中
						msg := msgHead.(*Message)
						msg.Data = make([]byte, msg.GetDataLen())
						// 根据dataLen从io中读取字节流
						_, err := io.ReadFull(conn, msg.Data)
						if err != nil {
							t.Error("server unpack data err:", err)
							return
						}
						// 完整的一个消息已经读取完毕
						fmt.Println("==> Recv Msg: ID=", msg.Id, ", len=", msg.DataLen, ", data=", string(msg.Data))
					}
				}
			}(conn)
		}
	}()

	// 模拟客户端
	conn, err := net.Dial("tcp4", "127.0.0.1:7777")
	if err != nil {
		t.Error("client dial err:", err)
		return
	}
	// 创建一个封包对象 dp
	dp := NewDataPack()
	// 模拟粘包过程，封装两个msg一同发送
	// 封装第一个msg1包
	msg1 := &Message{
		Id:      0,
		DataLen: 5,
		Data:    []byte{'h', 'e', 'l', 'l', 'o'},
	}
	sendData1, err := dp.Pack(msg1)
	if err != nil {
		t.Error("client pack msg1 error")
		return
	}
	// 封装第二个msg2包
	msg2 := &Message{
		Id:      1,
		DataLen: 7,
		Data:    []byte{'z', 'i', 'n', 'x', ' ', 'v', '0'},
	}
	sendData2, err := dp.Pack(msg2)
	if err != nil {
		t.Error("client pack msg2 error")
		return
	}
	// 将两个包粘在一起
	sendData1 = append(sendData1, sendData2...)
	// 一次性发送给服务端
	conn.Write(sendData1)

	// 客户端阻塞
	select {}
}
