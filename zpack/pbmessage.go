package zpack

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

type PBAsyncMessage struct {
	Cmd  uint32
	Body []byte
}

type PBMessage struct {
	Message
}

// GetData 获取消息内容
func (msg *PBMessage) GetData() []byte {
	//byte[]内容反序列化 成 PBAsyncMessage对象
	dataBuff := bytes.NewReader(msg.Data)
	pbMsg := &PBAsyncMessage{}
	err := binary.Read(dataBuff, binary.BigEndian, &pbMsg)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return pbMsg.Body
}
