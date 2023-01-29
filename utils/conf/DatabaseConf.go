package conf

type DatabaseConf struct {
	Host     string //mongodb uri
	Port     int    //端口
	Username string //用户名
	Password string //密码
	Database string //数据库名称
}
