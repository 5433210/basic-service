package store

import (
	"context"
	"encoding/json"

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
	if err = txn.Put([]byte("job:"+job.Id), bytes); err != nil {
		return err
	}

	return nil
}

func (j *JobStore) Retrieve(txn *Transaction, jobId string) (*apiv1.Job, error) {
	bytes, err := txn.Get(context.TODO(), []byte("job:"+jobId))
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

	return txn.Put([]byte("job:"+job.Id), bytes)
}

func (j *JobStore) Delete(txn *Transaction, jobId string) error {
	return txn.Delete([]byte("job:" + jobId))
}

func (j *JobStore) RetrieveAll(txn *Transaction) (*[]apiv1.Job, error) {
	jobs := make([]apiv1.Job, 0)
	it, err := txn.txn.Iter([]byte("job:"), nil)
	if err != nil {
		return nil, err
	}
	for it.Valid() {
		job := &apiv1.Job{}
		if err = json.Unmarshal(it.Value(), job); err != nil {
			it.Close()
			jobs = append(jobs, *job)

			return nil, err
		}
	}

	return &jobs, nil
}

func (j *JobStore) RetrieveExecutions(txn *Transaction, jobId string) (*[]apiv1.Execution, error) {
	executions := make([]apiv1.Execution, 0)
	it, err := txn.txn.Iter([]byte("execution:"+jobId+":"), nil)
	if err != nil {
		return nil, err
	}
	for it.Valid() {
		execution := &apiv1.Execution{}
		if err = json.Unmarshal(it.Value(), execution); err != nil {
			it.Close()
			executions = append(executions, *execution)

			return nil, err
		}
	}

	return &executions, nil
}
