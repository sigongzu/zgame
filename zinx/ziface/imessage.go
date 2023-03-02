package ziface

// 将请求的一个消息封装到一个message中，定义抽象层接口
type IMessage interface {
	GetDataLen() uint32 // 获取消息的长度
	GetMsgID() uint32   // 获取消息的ID
	GetData() []byte    // 获取消息的内容
	SetMsgID(uint32)    // 设置消息的ID
	SetData([]byte)     // 设置消息的内容
	SetDataLen(uint32)  // 设置消息的长度
}
