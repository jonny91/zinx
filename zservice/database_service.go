package zservice

import (
	"context"
	"github.com/jonny91/zinx/ziface"
)

type DatabaseDefine string

type DatabaseService struct {
	dbImpl map[DatabaseDefine]ziface.IDatabase
}

func (d *DatabaseService) Init(ctx context.Context) (bool, error) {
	return true, nil
}

func (d *DatabaseService) GetName() string {
	return "database"
}
