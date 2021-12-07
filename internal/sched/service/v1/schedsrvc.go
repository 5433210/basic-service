package servicev1

import (
	apiv1 "wailik.com/internal/sched/api/v1"
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
	return nil, nil
}

// (POST /jobs).
func (s *schedSrvc) CreateJob(job apiv1.Job) error {
	return nil
}

// (DELETE /jobs/{jobId}).
func (s *schedSrvc) DeleteJob(jobId string) error {
	return nil
}

// (GET /jobs/{jobId}).
func (s *schedSrvc) GetJobById(jobId string) (*apiv1.Job, error) {
	return nil, nil
}

// (PATCH /jobs/{jobId}).
func (s *schedSrvc) UpdateJob(jobId string, job apiv1.Job) error {
	return nil
}

// (GET /jobs/{jobId}/executions).
func (s *schedSrvc) GetJobAllExecutions(jobId string) (*[]apiv1.Execution, error) {
	return nil, nil
}
