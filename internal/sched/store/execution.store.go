package store

import (
	"context"
	"encoding/json"

	apiv1 "wailik.com/internal/sched/api/v1"
)

// execution
// key:	executions:job_id:execution_id
// create
// retrieve
// update
// delete.
type ExecutionStore struct{}

func (e *ExecutionStore) Create(txn *Transaction, execution apiv1.Execution) error {
	bytes, err := json.Marshal(execution)
	if err != nil {
		return err
	}
	if err = txn.Put(context.Background(), []byte("execution:"+execution.JobId+":"+execution.Id), bytes); err != nil {
		return err
	}

	return nil
}

func (e *ExecutionStore) Retrieve(txn *Transaction, jobId string, id string) (*apiv1.Execution, error) {
	bytes, err := txn.Get(context.Background(), []byte("execution:"+jobId+":"+id))
	if err != nil {
		return nil, err
	}
	execution := &apiv1.Execution{}
	err = json.Unmarshal(bytes, execution)
	if err != nil {
		return nil, err
	}

	return execution, nil
}

func (e *ExecutionStore) Update(txn *Transaction, execution apiv1.Execution) error {
	bytes, err := json.Marshal(execution)
	if err != nil {
		return err
	}
	if err = txn.Put(context.Background(), []byte("execution:"+execution.JobId+":"+execution.Id), bytes); err != nil {
		return err
	}

	return nil
}

func (e *ExecutionStore) Delete(txn *Transaction, jobId string, id string) error {
	return txn.Delete(context.Background(), []byte("execution:"+jobId+":"+id))
}
