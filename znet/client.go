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
	go c.Conn.Start()
}

func (c *Client) AddRouter(msgID uint32, router ziface.IRouter) {
	c.NetPort.AddRouter(msgID, router)
}

func (c *Client) GetConn() ziface.IConnection {
	return c.Conn
}

func (c *Client) StartAsClient() {
	c.exitChan = make(chan struct{})

	go func() {

		select {
		case <-c.exitChan:
			err := c.GetConn().GetTCPConnection().Close()
			if err != nil {
				fmt.Println("client close err ", err)
			}
		}
	}()
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
