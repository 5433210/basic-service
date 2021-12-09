package microservice

import (
	"strings"
	"sync"

	"stathat.com/c/consistent"

	"wailik.com/internal/pkg/log"
	"wailik.com/internal/pkg/util"
)

type ServiceNode struct {
	Addr     string
	Name     string
	UniqueId string
	RunMode  SrvcRunMode
	PickMode SrvcPickMode
	IsMaster bool
}

type ServiceGroup struct {
	Name     string
	RunMode  SrvcRunMode
	PickMode SrvcPickMode
}

type NodeManager struct {
	sync.RWMutex
	// <name,<id,node>>
	nodes  map[string]map[string]*ServiceNode
	groups map[string]*ServiceGroup

	hashNodes map[string]*consistent.Consistent
}

func NewNodeManager() *NodeManager {
	return &NodeManager{
		nodes:     map[string]map[string]*ServiceNode{},
		hashNodes: map[string]*consistent.Consistent{},
		groups:    map[string]*ServiceGroup{},
	}
}

func (mgr *NodeManager) GetNameById(id string) string {
	sli := strings.Split(id, "/")

	return sli[len(sli)-2]
}

func (mgr *NodeManager) HasMaster(name string) bool {
	mgr.Lock()
	defer mgr.Unlock()
	if mgr.nodes == nil {
		return false
	}
	log.Debugf("%+v", mgr.nodes)
	nodes := mgr.nodes[name]
	for _, node := range nodes {
		if node.IsMaster {
			log.Debug("node:" + node.UniqueId + " is master")

			return true
		}
	}

	return false
}

func (mgr *NodeManager) HasNode(nodeId string) bool {
	mgr.Lock()
	defer mgr.Unlock()
	if mgr.nodes == nil {
		return false
	}
	log.Debugf("%+v", mgr.nodes)
	log.Debugf("%+v", nodeId)
	name := GetNodeNameByNodeId(nodeId)
	log.Debugf("node name:%+v", name)
	nodes := mgr.nodes[name]
	for _, node := range nodes {
		if node.UniqueId == nodeId {
			return true
		}
	}

	return false
}

func (mgr *NodeManager) GetMaster(name string) string {
	mgr.Lock()
	defer mgr.Unlock()
	nodes := mgr.nodes[name]

	for _, node := range nodes {
		if node.IsMaster {
			return node.UniqueId
		}
	}

	return ""
}

func (mgr *NodeManager) SetMaster(name string, node *ServiceNode) {
	node.IsMaster = true
	mgr.AddNode(node)
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

	if _, exist := mgr.groups[node.Name]; !exist {
		mgr.groups[node.Name] = &ServiceGroup{
			Name:     node.Name,
			RunMode:  node.RunMode,
			PickMode: node.PickMode,
		}
	}

	mgr.nodes[node.Name][node.UniqueId] = node

	if node.PickMode == SrvcPickModeHash {
		if _, exist := mgr.hashNodes[node.Name]; !exist {
			mgr.hashNodes[node.Name] = consistent.New()
		}
		mgr.hashNodes[node.Name].Add(node.UniqueId)
	}
}

func (mgr *NodeManager) DelNode(id string) {
	sli := strings.Split(id, "/")
	name := sli[len(sli)-2]
	mgr.Lock()
	defer mgr.Unlock()
	if _, exist := mgr.nodes[name]; exist {
		if node, exist := mgr.nodes[name][id]; exist {
			if node.PickMode == SrvcPickModeHash {
				mgr.hashNodes[node.Name].Remove(node.UniqueId)
			}
		}
		delete(mgr.nodes[name], id)
	}
}

func (mgr *NodeManager) Pick(name string, hashValue string) *ServiceNode {
	mgr.RLock()
	defer mgr.RUnlock()
	nodes, exist := mgr.nodes[name]
	if !exist {
		return nil
	}
	mode := mgr.groups[name].PickMode
	// 纯随机取节点
	if mode == SrvcPickModeRandom {
		idx := util.RandomInt(len(nodes), 0)
		for _, v := range nodes {
			if idx == 0 {
				return v
			}
			idx--
		}
	}

	if mode == SrvcPickModeHash {
		id, err := mgr.hashNodes[name].Get(hashValue)
		if err != nil {
			return nil
		}

		return mgr.nodes[name][id]
	}

	if mode == SrvcPickModeMaster {
		nodes := mgr.nodes[name]
		for _, v := range nodes {
			if v.IsMaster {
				return v
			}
		}
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
