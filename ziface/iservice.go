package ziface

import "context"

// IService 基本服务接口
type IService interface {
	Init(ctx context.Context) (bool, error)
	GetName() string
}
