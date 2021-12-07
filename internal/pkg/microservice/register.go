package microservice

import (
	"context"
	"encoding/json"
	"time"

	"github.com/coreos/etcd/clientv3"
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

func NewRegister(node *ServiceNode, conf ClientConfig) (reg *Register, err error) {
	r := &Register{}
	r.closeChan = make(chan error)
	r.node = node
	r.cli, err = clientv3.New(clientv3.Config(conf))

	return r, err
}

func (r *Register) Run() {
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
			if err = r.register(); err != nil {
				return err
			}
			err = nil
		} else {
			return err
		}
	}

	// log.Infof("keepalived leaseId:%+v", r.leaseId)

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
