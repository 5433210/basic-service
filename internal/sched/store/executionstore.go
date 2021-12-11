package store

import apiv1 "wailik.com/internal/sched/api/v1"

// execution
// key:	executions:job_id:execution_id
// create
// retrieve
// update
// delete.
type ExecutionDao struct{}

func (e *ExecutionDao) Create(txn *Transaction, execution apiv1.Execution) error {
	return nil
}

func (e *ExecutionDao) Retrieve(txn *Transaction, id string) (*apiv1.Execution, error) {
	return nil, nil
}

func (e *ExecutionDao) Update(txn *Transaction, execution apiv1.Execution) error {
	return nil
}

func (e *ExecutionDao) Delete(txn *Transaction, id string) error {
	return nil
}
