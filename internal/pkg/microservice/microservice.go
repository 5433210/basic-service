package microservice

import (
	"encoding/json"
	"io/ioutil"

	"github.com/coreos/etcd/clientv3"
	"gopkg.in/yaml.v2"

	"wailik.com/internal/pkg/log"
)

type ClientConfig clientv3.Config

type microService struct {
	register  *Register
	discovery *Discovery
	manager   *NodeManager
	running   bool
}

type MicroService interface {
	Run()
	Stop()
	Pull() error
	Pick(name string) *ServiceNode
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

func New(node ServiceNode, clientConfigPath string) (*microService, error) {
	conf, err := loadClientConfig(clientConfigPath)
	if err != nil {
		return nil, err
	}

	manager := NewNodeManager()

	discovery, err := NewDiscovery(*conf, manager)
	if err != nil {
		return nil, err
	}

	if err = discovery.Pull(); err != nil {
		return nil, err
	}

	register, err := NewRegister(&node, *conf)
	if err != nil {
		return nil, err
	}

	if err = register.register(); err != nil {
		return nil, err
	}

	return &microService{discovery: discovery, register: register, manager: manager, running: false}, nil
}

func (ms *microService) Run() {
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

func (ms *microService) Pull() error {
	return ms.discovery.Pull()
}

func (ms *microService) Pick(name string) *ServiceNode {
	return ms.manager.Pick(name)
}

func loadClientConfig(path string) (*ClientConfig, error) {
	conf := clientv3.Config{}
	confY := make(map[interface{}]interface{})
	f, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	if err = yaml.Unmarshal(f, &confY); err != nil {
		return nil, err
	}

	log.Debugf("confY:%+v", confY)

	confJ := make(map[string]interface{})
	for key, value := range confY {
		switch key := key.(type) {
		case string:
			confJ[key] = value
		}
	}

	j, err := json.Marshal(&confJ)
	if err != nil {
		return nil, err
	}

	log.Debugf("json:%+v", string(j))

	if err := json.Unmarshal(j, &conf); err != nil {
		return nil, err
	}

	log.Debugf("config:%+v", conf)

	return (*ClientConfig)(&conf), nil
}
