package microservice

import (
	"go.etcd.io/etcd/clientv3"

	"wailik.com/internal/pkg/log"
)

const (
	_ttl = 10
)

type (
	SrvcRunMode  string
	SrvcPickMode string
)

const (
	SrvcRunModeFair        SrvcRunMode = "fair"
	SrvcRunModeMasterSlave SrvcRunMode = "master-slave"
)

const (
	SrvcPickModeRandom SrvcPickMode = "random"
	SrvcPickModeHash   SrvcPickMode = "hash"
	SrvcPickModeMaster SrvcPickMode = "master"
)

type (
	Registered func() error
	BeenMaster func() error
)

type MicroServiceConfig struct {
	Disabled  bool
	Endpoints []string
}

type ClientConfig clientv3.Config

type microService struct {
	register  *Register
	discovery *Discovery
	manager   *NodeManager
	running   bool
}

type MicroService interface {
	Start()
	Stop()
	Pick(name string, value string) *ServiceNode
}

type MicroServiceHelper interface {
	SetMicroService(ms MicroService)
	GetMicroService() MicroService
}

type MicroServiceObject struct {
	microService MicroService
}

func (s *MicroServiceObject) SetMicroService(ms MicroService) {
	s.microService = ms
}

func (s *MicroServiceObject) GetMicroService() MicroService {
	return s.microService
}

func New(node ServiceNode, conf MicroServiceConfig, registered Registered, beenMaster BeenMaster) (*microService, error) {
	log.Info("new micro service...")

	var clientConf = ClientConfig{
		Endpoints: conf.Endpoints,
	}

	manager := NewNodeManager()

	discovery, err := NewDiscovery(clientConf, manager)
	if err != nil {
		return nil, err
	}

	register, err := NewRegister(&node, clientConf, registered, beenMaster)
	if err != nil {
		return nil, err
	}

	return &microService{
		discovery: discovery,
		register:  register,
		manager:   manager,
		running:   false,
	}, nil
}

func (ms *microService) Start() {
	log.Info("run micro service...")
	if ms.running {
		return
	}
	go ms.register.Run()
	go ms.discovery.Run()
	ms.running = true
}

func (ms *microService) Stop() {
	if ms.running {
		ms.discovery.Stop()
		ms.register.Stop()
		ms.running = false
	}
}

func (ms *microService) Pick(name string, value string) *ServiceNode {
	return ms.manager.Pick(name, value)
}
