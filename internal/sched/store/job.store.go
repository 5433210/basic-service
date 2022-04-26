package store

import (
	"bytes"
	"context"
	"encoding/json"

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
	startKey := []byte("job:")
	jobs := make([]apiv1.Job, 0)
	it, err := txn.txn.Iter(startKey, nil)
	if err != nil {
		return nil, err
	}
	defer it.Close()

	for it.Valid() {
		if !bytes.HasPrefix(it.Key(), startKey) {
			break
		}

		job := &apiv1.Job{}
		if err = json.Unmarshal(it.Value(), job); err != nil {
			log.Debugf("%+v", string(it.Value()))

			return nil, err
		}
		jobs = append(jobs, *job)

		if err = it.Next(); err != nil {
			return nil, err
		}
	}

	return &jobs, nil
}

func (j *JobStore) RetrieveExecutions(txn *Transaction, jobId string) (*[]apiv1.Execution, error) {
	startKey := []byte("execution:" + jobId + ":")
	executions := make([]apiv1.Execution, 0)
	it, err := txn.txn.Iter(startKey, nil)
	if err != nil {
		return nil, err
	}
	defer it.Close()
	for it.Valid() {
		if !bytes.HasPrefix(it.Key(), startKey) {
			break
		}
		execution := &apiv1.Execution{}
		if err = json.Unmarshal(it.Value(), execution); err != nil {
			return nil, err
		}
		executions = append(executions, *execution)
		if err = it.Next(); err != nil {
			return nil, err
		}
	}

	return &executions, nil
}
