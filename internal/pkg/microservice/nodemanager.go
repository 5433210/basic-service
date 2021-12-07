package microservice

import (
	"strings"
	"sync"

	"wailik.com/internal/pkg/log"
	"wailik.com/internal/pkg/util"
)

type ServiceNode struct {
	Addr     string
	Name     string
	UniqueId string
}

type NodeManager struct {
	sync.RWMutex
	// <name,<id,node>>
	nodes map[string]map[string]*ServiceNode
}

func NewNodeManager() *NodeManager {
	return &NodeManager{
		nodes: map[string]map[string]*ServiceNode{},
	}
}

func (mgr *NodeManager) AddNode(node *ServiceNode) {
	if node == nil {
		return
	}
	mgr.Lock()
	defer mgr.Unlock()
	if _, exist := mgr.nodes[node.Name]; !exist {
		mgr.nodes[node.Name] = map[string]*ServiceNode{}
	}
	mgr.nodes[node.Name][node.UniqueId] = node
}

func (mgr *NodeManager) DelNode(id string) {
	sli := strings.Split(id, "/")
	name := sli[len(sli)-2]
	mgr.Lock()
	defer mgr.Unlock()
	if _, exist := mgr.nodes[name]; exist {
		delete(mgr.nodes[name], id)
	}
}

func (mgr *NodeManager) Pick(name string) *ServiceNode {
	mgr.RLock()
	defer mgr.RUnlock()
	nodes, exist := mgr.nodes[name]
	if !exist {
		return nil
	}
	// 纯随机取节点
	idx := util.RandomInt(len(nodes), 0)

	for _, v := range nodes {
		if idx == 0 {
			return v
		}
		idx--
	}

	return nil
}

func (mgr *NodeManager) Dump() {
	for k, v := range mgr.nodes {
		for kk, vv := range v {
			log.Infof("[NodeManager] Name:%s Id:%s Node:%+v", k, kk, vv)
		}
	}
}
