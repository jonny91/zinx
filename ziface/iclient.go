package ziface

type IClient interface {
	INetPort
	Dail(network string, ip string, port int)

	GetConn() IConnection
	StartAsClient()
}
