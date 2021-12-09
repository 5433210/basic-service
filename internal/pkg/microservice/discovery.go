package microservice

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/clientv3/concurrency"
	"github.com/coreos/etcd/mvcc/mvccpb"

	"wailik.com/internal/pkg/constant"
	"wailik.com/internal/pkg/log"
)

type Discovery struct {
	cli *clientv3.Client
	mgr *NodeManager
	cls chan error
}

func NewDiscovery(conf ClientConfig, mgr *NodeManager) (*Discovery, error) {
	if mgr == nil {
		return nil, fmt.Errorf("mgr == nil")
	}

	cli, err := clientv3.New(clientv3.Config(conf))
	if err != nil {
		return nil, err
	}

	return &Discovery{mgr: mgr, cli: cli, cls: make(chan error)}, nil
}

func (d *Discovery) Pull() error {
	kv := clientv3.NewKV(d.cli)
	resp, err := kv.Get(context.TODO(), "discovery/", clientv3.WithPrefix())
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
		d.mgr.AddNode(node)
		log.Debugf("pull node:%+v", node)
	}

	resp, err = kv.Get(context.TODO(), constant.MasterPrifex+"/", clientv3.WithPrefix())
	if err != nil {
		log.Fatalf("kv.Get err:%+v", err)

		return err
	}
	for _, v := range resp.Kvs {
		nodeId := string(v.Value)
		if nodeId == "" {
			log.Debug("nodeId is null")

			break
		}
		log.Debug("nodeId:" + nodeId)
		if !d.mgr.HasNode(nodeId) {
			if err := d.removeMasterKey(nodeId); err != nil {
				log.Fatalf("removeMasterKey(%+v) failed", nodeId)

				return err
			}
			log.Debugf("remove not use master key(%+v)", string(v.Key))
		}
	}

	return nil
}

func (d *Discovery) Run() {
	watcher := clientv3.NewWatcher(d.cli)
	watchChan := watcher.Watch(context.TODO(), "discovery", clientv3.WithPrefix())
	for {
		select {
		case resp := <-watchChan:
			d.watchEvent(resp.Events)
		case <-d.cls:
			goto EXIT
		}
	}

EXIT:
	log.Infof("discovery exit...")
}

func (d *Discovery) Stop() {
	close(d.cls)
	log.Info("stop discovery")
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
			d.mgr.AddNode(node)
			log.Infof("new node:%s", string(ev.Kv.Value))
		case mvccpb.DELETE:
			id := string(ev.Kv.Key)
			d.mgr.DelNode(id)
			err := d.removeMasterKey(id)
			if err != nil {
				continue
			}
			log.Infof("del node:%s data:%s", string(ev.Kv.Key), string(ev.Kv.Value))
		}
	}
}

func (d *Discovery) removeMasterKey(key string) error {
	client := d.cli
	nodeName := GetNodeNameByNodeId(key)
	session, err := concurrency.NewSession(client, concurrency.WithTTL(_ttl))
	if err != nil {
		return err
	}
	defer session.Close()

	masterKey := constant.MasterPrifex + "/" + nodeName

	mutex := concurrency.NewMutex(session, masterKey)
	if err = mutex.Lock(context.TODO()); err != nil {
		return err
	}
	defer func() {
		err = mutex.Unlock(context.TODO())
		if err != nil {
			log.Fatal("unlock mutex failed")
		}
	}()

	kv := clientv3.NewKV(client)
	resp, err := kv.Get(context.TODO(), masterKey)
	if err != nil {
		return err
	}
	if len(resp.Kvs) > 0 {
		for _, v := range resp.Kvs {
			if string(v.Value) == key {
				_, err = kv.Delete(context.TODO(), masterKey)
				if err != nil {
					return err
				}

				log.Debugf("remove masterKey:%+v value:%+v", v.Key, string(v.Value))
			}
		}
	}

	return nil
}
