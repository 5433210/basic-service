package store

import (
	"context"

	"github.com/tikv/client-go/v2/tikv"

	"wailik.com/internal/pkg/util"
)

type Store struct {
	pool *clientPool
}

type Client struct {
	client *tikv.KVStore
	txn    *Transaction
}

type clientPool struct {
	size    int
	clients []*Client
}

type Transaction struct {
	txn *tikv.KVTxn
}

func (t *Transaction) Commit(ctx context.Context) error {
	return t.txn.Commit(ctx)
}

func (t *Transaction) Rollback() error {
	return t.txn.Rollback()
}

func (t *Transaction) Get(ctx context.Context, k []byte) ([]byte, error) {
	return t.txn.Get(ctx, k)
}

func (t *Transaction) Delete(k []byte) error {
	return t.txn.Delete(k)
}

func (t *Transaction) Put(k []byte, v []byte) error {
	return t.txn.Set(k, v)
}

func newPool(endpoint []string, size int) (pool *clientPool, err error) {
	pool = &clientPool{
		size: size,
	}
	for i := 0; i < size; i++ {
		pool.clients[i], err = pool.newClient(endpoint)
		if err != nil {
			return nil, err
		}
	}

	return pool, nil
}

func (p *clientPool) newClient(endpoint []string) (*Client, error) {
	client, err := tikv.NewTxnClient(endpoint)
	if err != nil {
		return nil, err
	}

	return &Client{
		client: client,
	}, nil
}

func NewStore(endpoint []string, size int) (*Store, error) {
	pool, err := newPool(endpoint, size)
	if err != nil {
		return nil, err
	}

	return &Store{
		pool: pool,
	}, nil
}

func (r *Store) Obtain() *Client {
	return r.pool.clients[util.RandomInt(r.pool.size, 0)]
}

func (c *Client) Txn(isPessimistic bool) (*Transaction, error) {
	c.txn = &Transaction{}
	txn, err := c.client.Begin()
	if err != nil {
		return nil, err
	}
	txn.SetPessimistic(isPessimistic)
	c.txn.txn = txn

	return c.txn, nil
}
