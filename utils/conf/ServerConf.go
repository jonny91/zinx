package conf

type ServerConf struct {
	Host    string //当前服务器主机IP
	TCPPort int    //当前服务器主机监听端口号
	Name    string //当前服务器名称
}
