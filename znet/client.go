package znet

import (
	"fmt"
	"github.com/jonny91/zinx/ziface"
	"github.com/jonny91/zinx/zpack"
	"net"
)

type Client struct {
	NetPort
	Conn *Connection
}

func (c *Client) Dail(network string, ip string, port int) {
	addr := fmt.Sprintf("%s:%d", ip, port)

	address, _ := net.ResolveTCPAddr(network, addr)
	var conn, err = net.DialTCP(network, nil, address)
	if err != nil {
		fmt.Println("client start err, exit!")
		return
	}
	c.Conn = NewClientConnection(c, conn, 0, c.msgHandler)
}

func (c *Client) AddRouter(msgID uint32, router ziface.IRouter) {
	c.NetPort.AddRouter(msgID, router)
}

func (c *Client) GetConn() ziface.IConnection {
	return c.Conn
}

func NewClient(name string) ziface.IClient {
	client := &Client{
		NetPort: NetPort{
			Name:       name,
			msgHandler: NewMsgHandle(),
			exitChan:   nil,
			packet:     zpack.Factory().NewPack(ziface.ZinxMessage),
		},
	}
	return client
}
