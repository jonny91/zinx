package ziface

type IClient interface {
	Dail(network string, ip string, port int)
	AddRouter(msgID uint32, router IRouter)
	GetConn() IConnection
}
