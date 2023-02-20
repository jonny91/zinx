package zservice

import "github.com/jonny91/zinx/ziface"

type DatabaseDefine string

type DatabaseService struct {
	dbImpl map[DatabaseDefine]ziface.IDatabase
}

func (d *DatabaseService) Init() {
}

func (d *DatabaseService) GetName() string {
	return "database"
}
