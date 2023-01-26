package znet

import "github.com/jonny91/zinx/ziface"

// NetPort 网络上的一个节点 可能是Client 也可能是Server
type NetPort struct {
	//服务器的名称
	Name     string
	exitChan chan struct{}
	packet   ziface.IDataPack
	//当前Server的消息管理模块，用来绑定MsgID和对应的处理方法
	msgHandler ziface.IMsgHandle
}

// AddRouter 路由功能：给当前服务注册一个路由业务方法，供客户端链接处理使用
func (p *NetPort) AddRouter(msgID uint32, router ziface.IRouter) {
	p.msgHandler.AddRouter(msgID, router)
}
