package microservice

import (
	"context"
	"encoding/json"
	"time"

	"github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/clientv3/concurrency"
	"github.com/coreos/etcd/etcdserver/api/v3rpc/rpctypes"

	"wailik.com/internal/pkg/constant"
	"wailik.com/internal/pkg/errors"
	"wailik.com/internal/pkg/log"
)

const (
	_ttl = 10
)

type Register struct {
	client    *clientv3.Client
	leaseId   clientv3.LeaseID
	lease     clientv3.Lease
	node      *ServiceNode
	closeChan chan error
}

func NewRegister(node *ServiceNode, conf ClientConfig) (*Register, error) {
	log.Info("new register")
	client, err := clientv3.New(clientv3.Config(conf))
	if err != nil {
		return nil, err
	}

	r := &Register{
		client:    client,
		closeChan: make(chan error),
		node:      node,
	}

	return r, nil
}

func (r *Register) Run() {
	log.Info("register running...")
	dur := time.Second * 1
	timer := time.NewTicker(dur)

	if err := r.register(); err != nil {
		return
	}
	for {
		select {
		case <-timer.C:
			if err := r.keepalive(); err != nil {
				return
			}
		case <-r.closeChan:
			goto EXIT
		}
	}
EXIT:
	log.Infof("register exited")
}

func (r *Register) Stop() {
	if err := r.revoke(); err != nil {
		return
	}
	close(r.closeChan)
	log.Info("register stoped")
}

func (r *Register) register() error {
	if r.node.RunMode == SrvcRunModeFair {
		return r.registerFair()
	}

	if r.node.RunMode == SrvcRunModeMasterSlave {
		return r.registerMasterSlave()
	}

	return nil
}

func (r *Register) registerFair() error {
	log.Debug("registerFair")

	return r.saveNode()
}

func (r *Register) registerMasterSlave() error {
	log.Debug("registerMasterSlave")
	if err := r.setMaster(); err != nil {
		return err
	}
	if err := r.saveNode(); err != nil {
		return err
	}

	return nil
}

func (r *Register) setMaster() error {
	if r.node.RunMode != SrvcRunModeMasterSlave {
		return nil
	}
	session, err := concurrency.NewSession(r.client, concurrency.WithTTL(_ttl))
	if err != nil {
		return err
	}

	defer session.Close()

	masterKey := constant.MasterPrifex + "/" + r.node.Name
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

	masterValue, err := r.getValue(masterKey)
	if err != nil {
		return err
	}

	if masterValue == "" {
		log.Debug("no master")
		err = r.setValue(masterKey, r.node.UniqueId)
		if err != nil {
			return err
		}
		r.node.IsMaster = true
		log.Debug("set master:" + r.node.UniqueId)
	}

	return nil
}

func (r *Register) saveNode() error {
	r.leaseId = 0
	kv := clientv3.NewKV(r.client)
	r.lease = clientv3.NewLease(r.client)
	leaseResp, err := r.lease.Grant(context.TODO(), _ttl)
	if err != nil {
		return err
	}

	data, _ := json.Marshal(r.node)
	_, err = kv.Put(context.TODO(), r.node.UniqueId, string(data), clientv3.WithLease(leaseResp.ID))
	if err != nil {
		return err
	}
	r.leaseId = leaseResp.ID
	log.Infof("registered leaseId:%+v", r.leaseId)

	return nil
}

func (r *Register) keepalive() error {
	_, err := r.lease.KeepAliveOnce(context.TODO(), r.leaseId)
	if err != nil {
		if errors.Is(err, rpctypes.ErrLeaseNotFound) {
			log.Debug("keep alive")
			if err = r.register(); err != nil {
				return err
			}
		} else {
			return err
		}
	}

	if err = r.setMaster(); err != nil {
		return err
	}

	return err
}

func (r *Register) revoke() error {
	_, err := r.client.Revoke(context.TODO(), r.leaseId)
	if err != nil {
		return err
	}
	log.Infof("revoked leaseId:%+v", r.leaseId)

	return nil
}

func (r *Register) getValue(key string) (string, error) {
	kv := clientv3.NewKV(r.client)
	resp, err := kv.Get(context.TODO(), key)
	if err != nil {
		return "", err
	}
	if len(resp.Kvs) > 0 {
		for _, v := range resp.Kvs {
			return string(v.Value), nil
		}
	}

	log.Debug("got no value")

	return "", nil
}

func (r *Register) setValue(key string, value string) error {
	kv := clientv3.NewKV(r.client)
	_, err := kv.Put(context.TODO(), key, value)
	if err != nil {
		return err
	}

	return nil
}
