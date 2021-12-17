package microservice

import (
	"context"
	"encoding/json"
	"time"

	"go.etcd.io/etcd/api/v3/v3rpc/rpctypes"
	"go.etcd.io/etcd/clientv3"
	"go.etcd.io/etcd/clientv3/concurrency"

	"wailik.com/internal/pkg/constant"
	"wailik.com/internal/pkg/errors"
	"wailik.com/internal/pkg/log"
)

type Register struct {
	client     *clientv3.Client
	leaseId    clientv3.LeaseID
	lease      clientv3.Lease
	node       *ServiceNode
	closeChan  chan error
	timer      *time.Ticker
	registered Registered
	beenMaster BeenMaster
}

func NewRegister(node *ServiceNode, conf ClientConfig, registered Registered, beenMaster BeenMaster) (*Register, error) {
	log.Info("new register")
	client, err := clientv3.New(clientv3.Config(conf))
	if err != nil {
		return nil, err
	}

	r := &Register{
		client:     client,
		closeChan:  make(chan error),
		node:       node,
		registered: registered,
		beenMaster: beenMaster,
	}

	return r, nil
}

func (r *Register) Run() {
	log.Info("register running...")
	// 续租间隔为租约周期的一半，保证租约持续有效
	r.timer = time.NewTicker(time.Second * _ttl / 2)

	if err := r.register(); err != nil {
		return
	}
	for {
		select {
		case <-r.timer.C:
			if err := r.keepalive(); err != nil {
				log.Fatalf("keepalive error:%+v", err)
			}
		case <-r.closeChan:
			goto EXIT
		}
	}
EXIT:
	log.Infof("register exited")
}

func (r *Register) Stop() {
	r.lease.Close()
	if err := r.revoke(); err != nil {
		return
	}
	r.timer.Stop()
	close(r.closeChan)
	log.Info("register stopped")
}

func (r *Register) register() error {
	if r.node.RunMode == SrvcRunModeFair {
		return r.registerFair()
	}

	if r.node.RunMode == SrvcRunModeMasterSlave {
		return r.registerMasterSlave()
	}

	if r.registered != nil {
		if err := r.registered(); err != nil {
			return err
		}
	}

	return nil
}

func (r *Register) registerFair() error {
	log.Debug("registerFair")
	if err := r.grant(); err != nil {
		return err
	}

	return r.saveNode()
}

func (r *Register) registerMasterSlave() error {
	log.Debug("registerMasterSlave")
	if err := r.grant(); err != nil {
		return err
	}
	if _, err := r.campaignMaster(); err != nil {
		return err
	}

	return r.saveNode()
}

func (r *Register) campaignMaster() (bool, error) {
	session, err := concurrency.NewSession(r.client,
		concurrency.WithTTL(_ttl))
	if err != nil {
		log.Fatalf("create session error:%+v", err)

		return false, err
	}

	defer session.Close()

	masterNodeKey := constant.MasterPrifex + "/" + r.node.Name
	mutex := concurrency.NewMutex(session, masterNodeKey)
	if err = mutex.Lock(context.TODO()); err != nil {
		log.Fatalf("lock mutex error:%+v", err)

		return false, err
	}

	defer func() {
		err = mutex.Unlock(context.TODO())
		if err != nil {
			log.Fatalf("unlock mutex error:%+v", err)
		}
	}()

	masterNodeId, err := r.getValue(masterNodeKey)
	if err != nil {
		return false, err
	}

	// log.Debugf("masterNodeId:%+v", masterNodeId)

	if masterNodeId == "" {
		log.Debugf("no master node for service(%s) yet", r.node.Name)
		err = r.setValue(masterNodeKey, r.node.UniqueId)
		if err != nil {
			return false, err
		}
		r.node.IsMaster = true
		log.Debugf("node(%s) campaigned service(%s) master node!",
			r.node.UniqueId, r.node.Name)
		if r.beenMaster != nil {
			if err = r.beenMaster(); err != nil {
				return false, err
			}
		}

		return true, nil
	}

	return false, nil
}

func (r *Register) saveNode() error {
	kv := clientv3.NewKV(r.client)
	data, _ := json.Marshal(r.node)
	_, err := kv.Put(context.TODO(), r.node.UniqueId, string(data), clientv3.WithLease(r.leaseId))
	if err != nil {
		return err
	}
	log.Infof("registered leaseId:%+v", r.leaseId)

	return nil
}

func (r *Register) grant() error {
	r.leaseId = 0
	r.lease = clientv3.NewLease(r.client)
	resp, err := r.lease.Grant(context.TODO(), _ttl)
	if err != nil {
		return err
	}
	r.leaseId = resp.ID

	return nil
}

func (r *Register) keepalive() (err error) {
	_, err = r.lease.KeepAliveOnce(context.TODO(), r.leaseId)
	if err != nil {
		if errors.Is(err, rpctypes.ErrLeaseNotFound) {
			log.Infof("lease(%s) not found, register again", r.leaseId)
			if err = r.register(); err != nil {
				return err
			}
		} else {
			return err
		}
	}

	if r.node.RunMode == SrvcRunModeMasterSlave {
		var success bool
		success, err = r.campaignMaster()
		if err != nil {
			return err
		}
		if success {
			if err = r.saveNode(); err != nil {
				return err
			}
		}
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

	return "", nil
}

func (r *Register) setValue(key string, value string) error {
	log.Debugf("lease id:%+v", r.leaseId)
	kv := clientv3.NewKV(r.client)
	// 通过租约来进行节点存活的判断。租约过期键值对自动失效，意味着对应节点已经失效，无法续租保活
	_, err := kv.Put(context.TODO(), key, value, clientv3.WithLease(r.leaseId))
	if err != nil {
		return err
	}

	return nil
}
