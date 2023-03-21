package store

import (
	"context"
	"encoding/json"
	"fmt"

	"wailik.com/internal/pkg/log"
	apiv1 "wailik.com/internal/sched/api/v1"
)

// job
// key:	jos:job_id
// create
// retrieve
// update
// delete.
type JobStore struct{}

func (j *JobStore) Create(txn *Transaction, job apiv1.Job) error {
	bytes, err := json.Marshal(job)
	if err != nil {
		return err
	}
	if err = txn.Put(context.Background(), []byte("job:"+job.Id), bytes); err != nil {
		return err
	}

	return nil
}

func (j *JobStore) Retrieve(txn *Transaction, jobId string) (*apiv1.Job, error) {
	bytes, err := txn.Get(context.Background(), []byte("job:"+jobId))
	if err != nil {
		return nil, err
	}
	job := &apiv1.Job{}
	err = json.Unmarshal(bytes, job)
	if err != nil {
		return nil, err
	}

	return job, nil
}

func (j *JobStore) Update(txn *Transaction, job apiv1.Job) error {
	bytes, err := json.Marshal(job)
	if err != nil {
		return err
	}

	return txn.Put(context.Background(), []byte("job:"+job.Id), bytes)
}

func (j *JobStore) Delete(txn *Transaction, jobId string) error {
	return txn.Delete(context.Background(), []byte("job:"+jobId))
}

func (j *JobStore) RetrieveAll(txn *Transaction) (*[]apiv1.Job, error) {
	ctx := context.Background()
	jobs := make([]apiv1.Job, 0)
	var cursor uint64
	for {
		var keys []string
		var err error
		keys, cursor, err = txn.txn.Scan(ctx, cursor, "job:*", 0).Result()
		if err != nil {
			return nil, err
		}

		for _, key := range keys {
			fmt.Println("key", key)
			val := txn.txn.Get(ctx, key).Val()
			if val == "" {
				break
			}

			job := &apiv1.Job{}
			if err := json.Unmarshal([]byte(val), job); err != nil {
				return nil, err
			}
			jobs = append(jobs, *job)
		}

		if cursor == 0 { // no more keys
			break
		}
	}

	log.Debugf("scan finished")
	return &jobs, nil
}

func (j *JobStore) RetrieveExecutions(txn *Transaction, jobId string) (*[]apiv1.Execution, error) {
	ctx := context.Background()
	executions := make([]apiv1.Execution, 0)
	var cursor uint64
	for {
		var keys []string
		var err error
		keys, cursor, err = txn.txn.Scan(ctx, cursor, "execution:"+jobId+":", 0).Result()
		if err != nil {
			return nil, err
		}

		for _, key := range keys {
			fmt.Println("key", key)
			val := txn.txn.Get(ctx, key).Val()
			if val == "" {
				break
			}

			execution := &apiv1.Execution{}
			if err := json.Unmarshal([]byte(val), execution); err != nil {
				return nil, err
			}
			executions = append(executions, *execution)
		}

		if cursor == 0 { // no more keys
			break
		}
	}

	log.Debugf("scan finished")
	return &executions, nil
}
