// Package utils 提供zinx相关工具类函数
// 包括:
//
//	全局配置
//	配置文件加载
//
// 当前文件描述:
// @Title  globalobj.go
// @Description  相关配置文件定义及加载方式
// @Author  Aceld - Thu Mar 11 10:32:29 CST 2019
package utils

import (
	"bytes"
	"fmt"
	"github.com/jonny91/zinx/utils/conf"
	"github.com/jonny91/zinx/ziface"
	"github.com/jonny91/zinx/zservice"
	"github.com/spf13/viper"
	"os"

	"github.com/jonny91/zinx/utils/commandline/args"
	"github.com/jonny91/zinx/utils/commandline/uflag"
	"github.com/jonny91/zinx/zlog"
)

/*
存储一切有关Zinx框架的全局参数，供其他模块使用
一些参数也可以通过 用户根据 config.toml来配置
*/
type GlobalObj struct {
	Server   *conf.ServerConf
	Database *conf.DatabaseConf
	Zinx     *conf.ZinxConf

	/*
		config file path
	*/
	ConfFilePath string

	/*
		logger
	*/
	LogDir        string //日志所在文件夹 默认"./log"
	LogFile       string //日志文件名称   默认""  --如果没有设置日志文件，打印信息将打印至stderr
	LogDebugClose bool   //是否关闭Debug日志级别调试信息 默认false  -- 默认打开debug信息
}

/*
定义一个全局的对象
*/
var (
	GlobalObject  *GlobalObj
	DB            ziface.IDatabase
	ServiceCenter zservice.IServiceManager
)

// PathExists 判断一个文件是否存在
func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

// Reload 读取用户的配置文件
func (g *GlobalObj) Reload() {

	if confFileExists, _ := PathExists(g.ConfFilePath); confFileExists != true {
		fmt.Println("Config File ", g.ConfFilePath, " is not exist!!")
		return
	}

	viper.SetConfigType("toml")
	data, err := os.ReadFile(g.ConfFilePath)

	if err = viper.ReadConfig(bytes.NewBuffer(data)); err != nil {
		return
	}

	//将数据解析到struct中
	err = viper.Unmarshal(&g)
	//err = viper.UnmarshalKey("server", &g)
	fmt.Println("globalObj:", g)
	if err != nil {
		panic(err)
	}

	//Logger 设置
	if g.LogFile != "" {
		zlog.SetLogFile(g.LogDir, g.LogFile)
	}
	if g.LogDebugClose {
		zlog.CloseDebug()
	}
}

/*
提供init方法，默认加载
*/
func init() {
	pwd, err := os.Getwd()
	if err != nil {
		pwd = "."
	}

	// 初始化配置模块flag
	args.InitConfigFlag(pwd+"/conf/config.toml", "配置文件，如果没有设置，则默认为<exeDir>/conf/config.toml")
	// 初始化日志模块flag TODO
	// 解析
	uflag.Parse()
	// 解析之后的操作
	args.FlagHandle()

	//初始化GlobalObject变量，设置一些默认值
	GlobalObject = &GlobalObj{
		Server: &conf.ServerConf{
			Name:    "ZinxServerApp",
			TCPPort: 8999,
			Host:    "0.0.0.0",
		},
		Zinx: &conf.ZinxConf{
			Version:          "V1.0",
			MaxConn:          12000,
			MaxPacketSize:    4096,
			WorkerPoolSize:   10,
			MaxWorkerTaskLen: 1024,
			MaxMsgChanLen:    1024,
		},

		ConfFilePath: args.Args.ConfigFile,

		LogDir:        pwd + "/log",
		LogFile:       "",
		LogDebugClose: false,
	}
	//NOTE: 从配置文件中加载一些用户配置的参数
	GlobalObject.Reload()
}
