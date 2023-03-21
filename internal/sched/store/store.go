package store

import (
	"context"

	redis "github.com/redis/go-redis/v9"

	"wailik.com/internal/pkg/util"
)

type Store struct {
	pool *clientPool
}

type Client struct {
	client *redis.Client
}

type clientPool struct {
	size    int
	clients []*Client
}

type Transaction struct {
	txn *redis.Client
}

func newTransaction(client *redis.Client) *Transaction {
	return &Transaction{client}
}

func (t *Transaction) Commit(ctx context.Context) error {
	return nil
}

func (t *Transaction) Rollback() error {
	return nil
}

func (t *Transaction) Get(ctx context.Context, k []byte) ([]byte, error) {
	return t.txn.Get(ctx, string(k)).Bytes()
}

func (t *Transaction) Delete(ctx context.Context, k []byte) error {
	return t.txn.Del(ctx, string(k)).Err()
}

func (t *Transaction) Put(ctx context.Context, k []byte, v []byte) error {
	return t.txn.Set(ctx, string(k), string(v), 0).Err()
}

func newPool(endpoint []string, size int) (pool *clientPool, err error) {
	pool = &clientPool{
		clients: make([]*Client, size),
		size:    size,
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
	return &Client{
		client: redis.NewClient(&redis.Options{
			Addr:     endpoint[0],
			Password: "",
			DB:       0}),
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
	txn := newTransaction(c.client)
	return txn, nil
}
