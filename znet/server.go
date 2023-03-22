package znet

import (
	"errors"
	"fmt"
	"github.com/jonny91/zinx/utils"
	"github.com/jonny91/zinx/ziface"
	"github.com/jonny91/zinx/zpack"
	"net"
)

var zinxLogo = `                                        
              ██                        
              ▀▀                        
 ████████   ████     ██▄████▄  ▀██  ██▀ 
     ▄█▀      ██     ██▀   ██    ████   
   ▄█▀        ██     ██    ██    ▄██▄   
 ▄██▄▄▄▄▄  ▄▄▄██▄▄▄  ██    ██   ▄█▀▀█▄  
 ▀▀▀▀▀▀▀▀  ▀▀▀▀▀▀▀▀  ▀▀    ▀▀  ▀▀▀  ▀▀▀ 
                                        `
var topLine = `┌──────────────────────────────────────────────────────┐`
var borderLine = `│`
var bottomLine = `└──────────────────────────────────────────────────────┘`

// Server 接口实现，定义一个Server服务类
type Server struct {
	NetPort
	//tcp4 or other
	IPVersion string
	//服务绑定的IP地址
	IP string
	//服务绑定的端口
	Port int
	//当前Server的链接管理器
	ConnMgr ziface.IConnManager
}

// NewServer 创建一个服务器句柄
func NewServer(opts ...Option) ziface.IServer {
	printLogo()

	s := &Server{
		IPVersion: "tcp4",
		IP:        utils.GlobalObject.Server.Host,
		Port:      utils.GlobalObject.Server.TCPPort,
		ConnMgr:   NewConnManager(),
		NetPort: NetPort{
			Name:       utils.GlobalObject.Server.Name,
			msgHandler: NewMsgHandle(),
			packet:     zpack.Factory().NewPack(ziface.ProtobufDataPack),
			exitChan:   nil,
		},
	}

	for _, opt := range opts {
		opt(s)
	}

	return s
}

// NewUserConfServer 创建一个服务器句柄
func NewUserConfServer(config *utils.Config, opts ...Option) ziface.IServer {
	//打印logo
	printLogo()

	s := &Server{
		IPVersion: config.TcpVersion,
		IP:        config.Host,
		Port:      config.TcpPort,
		ConnMgr:   NewConnManager(),
		NetPort: NetPort{
			Name:       utils.GlobalObject.Server.Name,
			msgHandler: NewMsgHandle(),
			packet:     zpack.Factory().NewPack(ziface.ZinxDataPack),
			exitChan:   nil,
		},
	}
	//更替打包方式
	for _, opt := range opts {
		opt(s)
	}
	//刷新用户配置到全局配置变量
	utils.UserConfToGlobal(config)

	return s
}

//============== 实现 ziface.IServer 里的全部接口方法 ========

// Start 开启网络服务
func (s *Server) Start() {
	fmt.Printf("[START] Server name: %s,listenner at IP: %s, Port %d is starting\n", s.Name, s.IP, s.Port)
	s.exitChan = make(chan struct{})

	//开启一个go去做服务端Listener业务
	go func() {
		//0 启动worker工作池机制
		s.msgHandler.StartWorkerPool()

		//1 获取一个TCP的Addr
		addr, err := net.ResolveTCPAddr(s.IPVersion, fmt.Sprintf("%s:%d", s.IP, s.Port))
		if err != nil {
			fmt.Println("resolve tcp addr err: ", err)
			return
		}

		//2 监听服务器地址
		listener, err := net.ListenTCP(s.IPVersion, addr)
		if err != nil {
			panic(err)
		}

		//已经监听成功
		fmt.Println("start Zinx server  ", s.Name, " succ, now listenning...")

		//TODO server.go 应该有一个自动生成ID的方法
		var cID uint32
		cID = 0

		go func() {
			//3 启动server网络连接业务
			for {
				//3.1 设置服务器最大连接控制,如果超过最大连接，则等待
				if s.ConnMgr.Len() >= utils.GlobalObject.Zinx.MaxConn {
					fmt.Println("Exceeded the maxConnNum:", utils.GlobalObject.Zinx.MaxConn, ", Wait:", AcceptDelay.duration)
					AcceptDelay.Delay()
					continue
				}

				//3.2 阻塞等待客户端建立连接请求
				conn, err := listener.AcceptTCP()
				if err != nil {
					//Go 1.16+
					if errors.Is(err, net.ErrClosed) {
						fmt.Println("Listener closed")
						return
					}
					fmt.Println("Accept err ", err)
					AcceptDelay.Delay()
					continue
				}

				AcceptDelay.Reset()

				//3.3 处理该新连接请求的 业务 方法， 此时应该有 handler 和 conn是绑定的
				dealConn := NewServerConnection(s, conn, cID, s.msgHandler)
				cID++

				//3.4 启动当前链接的处理业务
				go dealConn.Start()
			}
		}()

		select {
		case <-s.exitChan:
			err := listener.Close()
			if err != nil {
				fmt.Println("Listener close err ", err)
			}
		}
	}()
}

// Stop 停止服务
func (s *Server) Stop() {
	fmt.Println("[STOP] Zinx server , name ", s.Name)

	//将其他需要清理的连接信息或者其他信息 也要一并停止或者清理
	s.ConnMgr.ClearConn()
	s.exitChan <- struct{}{}
	close(s.exitChan)
}

// Serve 运行服务
func (s *Server) Serve() {
	s.Start()

	//TODO Server.Serve() 是否在启动服务的时候 还要处理其他的事情呢 可以在这里添加

	//阻塞,否则主Go退出， listenner的go将会退出
	select {}
}

// StartAsServer 以为服务器方式运行
func (s *Server) StartAsServer() {
	s.Serve()
}

// AddRouter 路由功能：给当前服务注册一个路由业务方法，供客户端链接处理使用
func (s *Server) AddRouter(msgID uint32, router ziface.IRouter) {
	s.NetPort.AddRouter(msgID, router)
}

// GetConnMgr 得到链接管理
func (s *Server) GetConnMgr() ziface.IConnManager {
	return s.ConnMgr
}

func printLogo() {
	fmt.Println(zinxLogo)
	fmt.Println(topLine)
	fmt.Println(fmt.Sprintf("%s [Github] https://github.com/jonny91/zinx             %s", borderLine, borderLine))
	fmt.Println(bottomLine)
	fmt.Printf("[Zinx] Version: %s, MaxConn: %d, MaxPacketSize: %d\n",
		utils.GlobalObject.Zinx.Version,
		utils.GlobalObject.Zinx.MaxConn,
		utils.GlobalObject.Zinx.MaxPacketSize)
}

func init() {
}
