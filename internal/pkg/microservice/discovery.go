package microservice

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/coreos/etcd/clientv3"

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
		case clientv3.EventTypePut:
			node := &ServiceNode{}
			err := json.Unmarshal(ev.Kv.Value, node)
			if err != nil {
				log.Fatalf("json.Unmarshal err:%+v", err)

				continue
			}
			d.mgr.AddNode(node)
			log.Infof("new node:%s", string(ev.Kv.Value))
		case clientv3.EventTypeDelete:
			d.mgr.DelNode(string(ev.Kv.Key))
			log.Infof("del node:%s data:%s", string(ev.Kv.Key), string(ev.Kv.Value))
		}
	}
}
