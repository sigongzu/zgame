package znet

import (
	"bytes"
	"encoding/binary"
	"errors"
	"zgame/zinx/utils"
	"zgame/zinx/ziface"
)

// 封包、拆包 具体模块
type DataPack struct{}

// NewDataPack 初始化方法
func NewDataPack() *DataPack {
	return &DataPack{}
}

// GetHeadLen 获取包头长度方法
func (dp *DataPack) GetHeadLen() uint32 {
	// ID uint32(4字节) + DataLen uint32(4字节)
	return 8
}

// Pack 封包方法, dateLen + msgID + data
func (dp *DataPack) Pack(msg ziface.IMessage) ([]byte, error) {
	// 创建一个存放bytes字节的缓冲
	dataBuff := bytes.NewBuffer([]byte{})

	// 将dataLen写进dataBuff中
	if err := binary.Write(dataBuff, binary.LittleEndian, msg.GetDataLen()); err != nil {
		return nil, err
	}

	// 将msgID写进dataBuff中
	if err := binary.Write(dataBuff, binary.LittleEndian, msg.GetMsgID()); err != nil {
		return nil, err
	}

	// 将data写进dataBuff中
	if err := binary.Write(dataBuff, binary.LittleEndian, msg.GetData()); err != nil {
		return nil, err
	}

	return dataBuff.Bytes(), nil
}

// Unpack 拆包方法, 将包的head信息读出来，然后再根据head的信息里的data的长度，再读取data数据，封装成Message消息
func (dp *DataPack) Unpack(binaryData []byte) (ziface.IMessage, error) {
	// 创建一个从输入二进制数据的ioReader
	dataBuff := bytes.NewReader(binaryData)

	// 只解压head的信息，得到dataLen和msgID
	msg := &Message{}

	// 读dataLen
	if err := binary.Read(dataBuff, binary.LittleEndian, &msg.DataLen); err != nil {
		return nil, err
	}

	// 读msgID
	if err := binary.Read(dataBuff, binary.LittleEndian, &msg.Id); err != nil {
		return nil, err
	}

	// 判断dataLen是否已经超出了我们允许的最大包长度
	if utils.GlobalObject.MaxPackageSize > 0 && msg.DataLen > utils.GlobalObject.MaxPackageSize {
		return nil, errors.New("large msg data received")
	}

	return msg, nil
}
