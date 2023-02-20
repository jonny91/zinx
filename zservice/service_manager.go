package zservice

import (
	"github.com/jonny91/zinx/ziface"
)

type IServiceManager interface {
	Register(service ziface.IService)
	Get(serviceName string) *ziface.IService
}

type ServiceManager struct {
	serviceCenter map[string]*ziface.IService
}

func (s *ServiceManager) Init() {
	s.serviceCenter = make(map[string]*ziface.IService)
}

func (s *ServiceManager) Register(service ziface.IService) {
	s.serviceCenter[service.GetName()] = &service
}

func (s *ServiceManager) Get(serviceName string) *ziface.IService {
	return s.serviceCenter[serviceName]
}

func (s *ServiceManager) GetName() string {
	return "ServiceManager"
}
