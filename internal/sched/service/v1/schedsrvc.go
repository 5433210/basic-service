package servicev1

import (
	"context"

	apiv1 "wailik.com/internal/sched/api/v1"
	"wailik.com/internal/sched/store"
)

type schedSrvc struct {
	service *service
}

func newSchedSrvc(s *service) *schedSrvc {
	return &schedSrvc{
		service: s,
	}
}

// (GET /jobs).
func (s *schedSrvc) GetAllJobs() (*[]apiv1.Job, error) {
	client := s.service.store.Obtain()
	txn, err := client.Txn(true)
	if err != nil {
		return nil, err
	}
	jobStore := &store.JobStore{}

	jobs, err := jobStore.RetrieveAll(txn)
	if err != nil {
		_ = txn.Rollback()

		return nil, err
	}

	_ = txn.Commit(context.Background())

	return jobs, nil
}

// (POST /jobs).
func (s *schedSrvc) CreateJob(job apiv1.Job) error {
	client := s.service.store.Obtain()
	txn, err := client.Txn(true)
	if err != nil {
		return err
	}
	jobStore := &store.JobStore{}

	err = jobStore.Create(txn, job)

	if err != nil {
		_ = txn.Rollback()

		return err
	}
	_ = txn.Commit(context.Background())

	err = s.service.Cron().Add(job)
	if err != nil {
		return err
	}

	return nil
}

// (DELETE /jobs/{jobId}).
func (s *schedSrvc) DeleteJob(jobId string) error {
	client := s.service.store.Obtain()
	txn, err := client.Txn(true)
	if err != nil {
		return err
	}
	jobStore := &store.JobStore{}

	err = jobStore.Delete(txn, jobId)
	if err != nil {
		_ = txn.Rollback()

		return err
	}
	_ = txn.Commit(context.Background())

	err = s.service.Cron().Remove(jobId)
	if err != nil {
		return err
	}

	return nil
}

// (GET /jobs/{jobId}).
func (s *schedSrvc) GetJobById(jobId string) (*apiv1.Job, error) {
	client := s.service.store.Obtain()
	txn, err := client.Txn(true)
	if err != nil {
		return nil, err
	}
	jobStore := &store.JobStore{}

	job, err := jobStore.Retrieve(txn, jobId)
	if err != nil {
		_ = txn.Rollback()

		return nil, err
	}
	_ = txn.Commit(context.Background())

	return job, nil
}

// (PATCH /jobs/{jobId}).
func (s *schedSrvc) UpdateJob(jobId string, job apiv1.Job) error {
	client := s.service.store.Obtain()
	txn, err := client.Txn(true)
	if err != nil {
		return err
	}
	jobStore := &store.JobStore{}

	err = jobStore.Update(txn, job)
	if err != nil {
		_ = txn.Rollback()

		return err
	}
	_ = txn.Commit(context.Background())

	err = s.service.Cron().Remove(jobId)
	if err != nil {
		return err
	}
	err = s.service.Cron().Add(job)
	if err != nil {
		return err
	}

	return nil
}

// (GET /jobs/{jobId}/executions).
func (s *schedSrvc) GetJobAllExecutions(jobId string) (*[]apiv1.Execution, error) {
	client := s.service.store.Obtain()
	txn, err := client.Txn(true)
	if err != nil {
		return nil, err
	}
	jobStore := &store.JobStore{}

	executions, err := jobStore.RetrieveExecutions(txn, jobId)
	if err != nil {
		_ = txn.Rollback()

		return nil, err
	}
	_ = txn.Commit(context.Background())

	return executions, nil
}
