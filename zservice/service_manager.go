package zservice

import (
	"context"
	"github.com/jonny91/zinx/ziface"
)

type IServiceManager[T ziface.IService] interface {
	Register(service ziface.IService)
	Get(serviceName string) ziface.IService
}

type ServiceManager struct {
	serviceCenter map[string]ziface.IService
}

func (s *ServiceManager) Register(service ziface.IService) {
	s.serviceCenter[service.GetName()] = service
}

func (s *ServiceManager) Get(serviceName string) ziface.IService {
	return s.serviceCenter[serviceName]
}

func NewServiceManager() *ServiceManager {
	sm := &ServiceManager{
		serviceCenter: map[string]ziface.IService{},
	}
	return sm
}

func (s *ServiceManager) Start(ctx context.Context) {
	for _, service := range s.serviceCenter {
		ok, err := service.Init(ctx)
		if !ok {
			panic(err)
		}
	}
}

func (s *ServiceManager) GetName() string {
	return "ServiceManager"
}
