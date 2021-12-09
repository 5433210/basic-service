package microservice

import (
	"context"
	"encoding/json"
	"time"

	"github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/clientv3/concurrency"
	"github.com/coreos/etcd/etcdserver/api/v3rpc/rpctypes"

	"wailik.com/internal/pkg/errors"
	"wailik.com/internal/pkg/log"
)

const (
	_ttl = 10
)

type Register struct {
	cli       *clientv3.Client
	leaseId   clientv3.LeaseID
	lease     clientv3.Lease
	node      *ServiceNode
	closeChan chan error
}

func NewRegister(node *ServiceNode, conf ClientConfig, mgr *NodeManager) (reg *Register, err error) {
	log.Info("new register")
	r := &Register{}
	r.closeChan = make(chan error)
	r.node = node
	r.cli, err = clientv3.New(clientv3.Config(conf))

	return r, err
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
	log.Infof("register exit...")
}

func (r *Register) Stop() {
	if err := r.revoke(); err != nil {
		return
	}
	close(r.closeChan)
	log.Info("stop register")
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

	return r.updateNode()
}

func (r *Register) registerMasterSlave() error {
	log.Debug("registerMasterSlave")
	if err := r.setMaster(); err != nil {
		return err
	}
	if err := r.updateNode(); err != nil {
		return err
	}

	return nil
}

func (r *Register) setMaster() error {
	if r.node.RunMode != SrvcRunModeMasterSlave {
		return nil
	}
	session, err := concurrency.NewSession(r.cli)
	if err != nil {
		return err
	}

	defer session.Close()

	masterKey := "master/" + r.node.Name
	mutex := concurrency.NewMutex(session, masterKey)
	if err = mutex.Lock(context.TODO()); err != nil {
		return err
	}

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

	if err = mutex.Lock(context.TODO()); err != nil {
		return err
	}

	return nil
}

func (r *Register) updateNode() error {
	r.leaseId = 0
	kv := clientv3.NewKV(r.cli)
	r.lease = clientv3.NewLease(r.cli)
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
	_, err := r.cli.Revoke(context.TODO(), r.leaseId)
	if err != nil {
		return err
	}
	log.Infof("revoked leaseId:%+v", r.leaseId)

	return nil
}

func (r *Register) getValue(key string) (string, error) {
	kv := clientv3.NewKV(r.cli)
	resp, err := kv.Get(context.TODO(), key)
	if err != nil {
		return "", err
	}
	if len(resp.Kvs) > 0 {
		for _, v := range resp.Kvs {
			// log.Debugf("get value(%+v) of key(%+v)", string(v.Value), key)

			return string(v.Value), nil
		}
	}

	log.Debug("got no value")

	return "", nil
}

func (r *Register) setValue(key string, value string) error {
	kv := clientv3.NewKV(r.cli)
	_, err := kv.Put(context.TODO(), key, value)
	if err != nil {
		return err
	}

	// log.Debugf("put key(%+v) value(%+v)", key, value)

	return nil
}
