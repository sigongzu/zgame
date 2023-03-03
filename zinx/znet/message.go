package znet

type Message struct {
	Id      uint32 // 消息的ID
	DataLen uint32 // 消息的长度
	Data    []byte // 消息的内容
}

// 创建一个Message消息包
func NewMsgPackage(id uint32, data []byte) *Message {
	return &Message{
		Id:      id,
		DataLen: uint32(len(data)),
		Data:    data,
	}
}

// 获取消息的长度
func (msg *Message) GetDataLen() uint32 {
	return msg.DataLen
}

// 获取消息的ID
func (msg *Message) GetMsgID() uint32 {
	return msg.Id
}

// 获取消息的内容
func (msg *Message) GetData() []byte {
	return msg.Data
}

// 设置消息的ID
func (msg *Message) SetMsgID(id uint32) {
	msg.Id = id
}

// 设置消息的内容
func (msg *Message) SetData(data []byte) {
	msg.Data = data
}

// 设置消息的长度
func (msg *Message) SetDataLen(len uint32) {
	msg.DataLen = len
}
