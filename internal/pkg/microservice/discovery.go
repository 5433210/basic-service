package microservice

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/mvcc/mvccpb"

	"wailik.com/internal/pkg/constant"
	"wailik.com/internal/pkg/log"
)

type Discovery struct {
	client    *clientv3.Client
	manager   *NodeManager
	closeChan chan error
	timer     *time.Ticker
	watcher   clientv3.Watcher
	watchChan clientv3.WatchChan
}

func NewDiscovery(conf ClientConfig, manager *NodeManager) (*Discovery, error) {
	if manager == nil {
		return nil, fmt.Errorf("mgr == nil")
	}

	client, err := clientv3.New(clientv3.Config(conf))
	if err != nil {
		return nil, err
	}

	return &Discovery{manager: manager, client: client, closeChan: make(chan error)}, nil
}

func (d *Discovery) pull() error {
	kv := clientv3.NewKV(d.client)
	resp, err := kv.Get(context.TODO(), constant.DiscoveryPrifex+"/", clientv3.WithPrefix())
	if err != nil {
		log.Fatalf("kv.Get err:%+v", err)

		return err
	}
	for _, v := range resp.Kvs {
		node := &ServiceNode{}
		err = json.Unmarshal(v.Value, node)
		if err != nil {
			log.Fatalf("json.Unmarshal err:%+v", err)

			continue
		}
		d.manager.AddNode(node)
		log.Debugf("pull node:%+v", node)
	}

	// resp, err = kv.Get(context.TODO(), constant.MasterPrifex+"/", clientv3.WithPrefix())
	// if err != nil {
	// 	log.Fatalf("kv.Get err:%+v", err)

	// 	return err
	// }
	// for _, v := range resp.Kvs {
	// 	nodeId := string(v.Value)
	// 	if nodeId == "" {
	// 		break
	// 	}

	// 	if !d.manager.HasNode(nodeId) {
	// 		log.Debugf("try to remove node:(%+v) of service(%+v) unavailable", nodeId, v.Key)
	// 		if err := d.removeMasterNode(nodeId); err != nil {
	// 			return err
	// 		}
	// 		log.Infof("removed master node(%+v) of service(%+v) unavailable ", string(v.Value), string(v.Key))
	// 	}
	// }

	return nil
}

func (d *Discovery) Run() {
	log.Info("discovery running...")
	d.watcher = clientv3.NewWatcher(d.client)
	d.watchChan = d.watcher.Watch(context.TODO(), constant.DiscoveryPrifex, clientv3.WithPrefix())
	d.timer = time.NewTicker(_ttl * time.Second / 2)

	for {
		select {
		case <-d.timer.C:
			if err := d.pull(); err != nil {
				log.Fatalf("pull error:%+v", err)
			}
		case resp := <-d.watchChan:
			d.watchEvent(resp.Events)
		case <-d.closeChan:
			goto EXIT
		}
	}

EXIT:
	log.Infof("discovery exited")
}

func (d *Discovery) Stop() {
	d.watcher.Close()
	close(d.closeChan)
	log.Info("discovery stopped")
}

func (d *Discovery) watchEvent(evs []*clientv3.Event) {
	for _, ev := range evs {
		switch ev.Type {
		case mvccpb.PUT:
			node := &ServiceNode{}
			err := json.Unmarshal(ev.Kv.Value, node)
			if err != nil {
				log.Fatalf("json.Unmarshal err:%+v", err)

				continue
			}
			d.manager.AddNode(node)
			log.Infof("add node:%s data:%+v", string(ev.Kv.Value), node)
		case mvccpb.DELETE:
			nodeId := string(ev.Kv.Key)
			d.manager.DelNode(nodeId)
			// err := d.removeMasterNode(nodeId)
			// if err != nil {
			// 	continue
			// }
			log.Infof("del node:%s", nodeId)
		}
	}
}

// func (d *Discovery) removeMasterNode(masterNodeId string) error {
// 	client := d.client
// 	nodeName := GetNodeNameFromNodeId(masterNodeId)
// 	session, err := concurrency.NewSession(client, concurrency.WithTTL(_ttl))
// 	if err != nil {
// 		log.Fatalf("create new session err:%+v", err)

// 		return err
// 	}
// 	defer session.Close()

// 	masterNodeKey := constant.MasterPrifex + "/" + nodeName

// 	mutex := concurrency.NewMutex(session, masterNodeKey)
// 	if err = mutex.Lock(context.TODO()); err != nil {
// 		log.Fatalf("lock mutex err:%+v", err)

// 		return err
// 	}
// 	defer func() {
// 		err = mutex.Unlock(context.TODO())
// 		if err != nil {
// 			log.Fatalf("unlock mutex err:%+v", err)
// 		}
// 	}()

// 	kv := clientv3.NewKV(client)
// 	resp, err := kv.Get(context.TODO(), masterNodeKey)
// 	if err != nil {
// 		log.Fatalf("kv get err:%+v", err)

// 		return err
// 	}
// 	if len(resp.Kvs) > 0 {
// 		for _, v := range resp.Kvs {
// 			nodeId := string(v.Value)
// 			if nodeId == masterNodeId {
// 				_, err = kv.Delete(context.TODO(), masterNodeKey)
// 				if err != nil {
// 					log.Fatalf("kv delete err:%+v", err)

// 					return err
// 				}

// 				log.Infof("removed master node(%+v) of service(%+v) unavailable ", masterNodeId, masterNodeKey)
// 			}
// 		}
// 	}

// 	return nil
// }
