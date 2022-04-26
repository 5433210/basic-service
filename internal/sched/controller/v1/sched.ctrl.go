package controllerv1

import (
	"github.com/gofiber/fiber/v2"

	"wailik.com/internal/pkg/constant"
	"wailik.com/internal/pkg/core"
	"wailik.com/internal/pkg/errors"
	apiv1 "wailik.com/internal/sched/api/v1"
	servicev1 "wailik.com/internal/sched/service/v1"
)

type schedController struct {
	srvc servicev1.Service
}

func newSchedController(c *controller) *schedController {
	return &schedController{srvc: c.srvc}
}

// (GET /jobs).
func (s *schedController) GetAllJobs(c *fiber.Ctx) error {
	r, err := s.srvc.Sched().GetAllJobs()
	if err != nil {
		return core.WriteResponse(c, err, nil)
	}

	return core.WriteResponse(c, nil, r)
}

// (POST /jobs).
func (s *schedController) CreateJob(c *fiber.Ctx) error {
	var json apiv1.CreateJobJSONRequestBody
	if err := c.BodyParser(&json); err != nil {
		return core.WriteResponse(c, errors.NewErrorC(errors.ErrCdDataBind, err), nil)
	}
	err := s.srvc.Sched().CreateJob(apiv1.Job(json))
	if err != nil {
		return core.WriteResponse(c, err, nil)
	}

	return core.WriteResponse(c, nil, nil)
}

// (DELETE /jobs/{jobId}).
func (s *schedController) DeleteJob(c *fiber.Ctx) error {
	jobId := c.Params(constant.FldJobId)

	err := s.srvc.Sched().DeleteJob(jobId)
	if err != nil {
		return core.WriteResponse(c, err, nil)
	}

	return core.WriteResponse(c, nil, nil)
}

// (GET /jobs/{jobId}).
func (s *schedController) GetJobById(c *fiber.Ctx) error {
	jobId := c.Params(constant.FldJobId)
	r, err := s.srvc.Sched().GetJobById(jobId)
	if err != nil {
		return core.WriteResponse(c, err, nil)
	}

	return core.WriteResponse(c, nil, r)
}

// (PATCH /jobs/{jobId}).
func (s *schedController) UpdateJob(c *fiber.Ctx) error {
	var json apiv1.CreateJobJSONRequestBody
	if err := c.BodyParser(&json); err != nil {
		return core.WriteResponse(c, errors.NewErrorC(errors.ErrCdDataBind, err), nil)
	}
	job := apiv1.Job(json)
	err := s.srvc.Sched().UpdateJob(job.Id, job)
	if err != nil {
		return core.WriteResponse(c, err, nil)
	}

	return core.WriteResponse(c, nil, nil)
}

// (GET /jobs/{jobId}/executions).
func (s *schedController) GetJobAllExecutions(c *fiber.Ctx) error {
	jobId := c.Params(constant.FldJobId)
	r, err := s.srvc.Sched().GetJobAllExecutions(jobId)
	if err != nil {
		return core.WriteResponse(c, err, nil)
	}

	return core.WriteResponse(c, nil, r)
}
