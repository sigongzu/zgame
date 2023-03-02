package ziface

// IRequest 定义抽象层接口
type IRequest interface {
	// 得到当前连接
	GetConnection() IConnection
	// 得到请求的消息数据
	GetData() []byte
	// 得到请求的消息的ID
	// GetMsgID() uint32
}
