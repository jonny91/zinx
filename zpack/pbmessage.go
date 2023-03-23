package zpack

type PBMessage struct {
	Message
}

// GetData 获取消息内容
func (msg *PBMessage) GetData() []byte {
	return msg.Message.GetData()
}
