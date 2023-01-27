package znet

import (
	"fmt"
	"github.com/jonny91/zinx/ziface"
)

// NetPort 网络上的一个节点 可能是Client 也可能是Server
type NetPort struct {
	ziface.INetPort
	//服务器的名称
	Name        string
	exitChan    chan struct{}
	packet      ziface.IDataPack
	msgHandler  ziface.IMsgHandle             //当前Server的消息管理模块，用来绑定MsgID和对应的处理方法
	OnConnStart func(conn ziface.IConnection) //该Server的连接创建时Hook函数
	OnConnStop  func(conn ziface.IConnection) //该Server的连接断开时的Hook函数
}

// AddRouter 路由功能：给当前服务注册一个路由业务方法，供客户端链接处理使用
func (p *NetPort) AddRouter(msgID uint32, router ziface.IRouter) {
	p.msgHandler.AddRouter(msgID, router)
}

func (p *NetPort) Packet() ziface.IDataPack {
	return p.packet
}

// SetOnConnStart 设置该Server的连接创建时Hook函数
func (p *NetPort) SetOnConnStart(hookFunc func(ziface.IConnection)) {
	p.OnConnStart = hookFunc
}

// SetOnConnStop 设置该Server的连接断开时的Hook函数
func (p *NetPort) SetOnConnStop(hookFunc func(ziface.IConnection)) {
	p.OnConnStop = hookFunc
}

// CallOnConnStart 调用连接OnConnStart Hook函数
func (p *NetPort) CallOnConnStart(conn ziface.IConnection) {
	if p.OnConnStart != nil {
		fmt.Println("---> CallOnConnStart....")
		p.OnConnStart(conn)
	}
}

// CallOnConnStop 调用连接OnConnStop Hook函数
func (p *NetPort) CallOnConnStop(conn ziface.IConnection) {
	if p.OnConnStop != nil {
		fmt.Println("---> CallOnConnStop....")
		p.OnConnStop(conn)
	}
}
