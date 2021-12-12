// Package apiv1 provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen version v1.8.3 DO NOT EDIT.
package apiv1

import (
	"time"
)

// ErrorResponse defines model for ErrorResponse.
type ErrorResponse struct {
	Code    string                  `json:"code"`
	Data    *map[string]interface{} `json:"data,omitempty"`
	Message string                  `json:"message"`
}

// Execution defines model for Execution.
type Execution struct {
	FinishedAt *string `json:"finished_at,omitempty"`
	Id         *string `json:"id,omitempty"`
	JobId      *string `json:"job_id,omitempty"`
	Output     *string `json:"output,omitempty"`
	ReturnCode *string `json:"return_code,omitempty"`
	StartedAt  *string `json:"started_at,omitempty"`
}

// Executor defines model for Executor.
type Executor struct {
	Config map[string]interface{} `json:"config"`
	Name   string                 `json:"name"`
}

// Job defines model for Job.
type Job struct {
	Data     *map[string]interface{} `json:"data,omitempty"`
	Disabled bool                    `json:"disabled"`
	Executor Executor                `json:"executor"`
	Id       string                  `json:"id"`
	Name     *string                 `json:"name,omitempty"`
	Next     *time.Time              `json:"next,omitempty"`
	Onetime  bool                    `json:"onetime"`
	Owner    string                  `json:"owner"`
	Schedule string                  `json:"schedule"`
	Timezone string                  `json:"timezone"`
}

// OkResponse defines model for OkResponse.
type OkResponse struct {
	Code    *string `json:"code,omitempty"`
	Message *string `json:"message,omitempty"`
}

// Error defines model for Error.
type Error ErrorResponse

// Ok defines model for Ok.
type Ok OkResponse

// CreateJobJSONBody defines parameters for CreateJob.
type CreateJobJSONBody Job

// UpdateJobJSONBody defines parameters for UpdateJob.
type UpdateJobJSONBody Job

// CreateJobJSONRequestBody defines body for CreateJob for application/json ContentType.
type CreateJobJSONRequestBody CreateJobJSONBody

// UpdateJobJSONRequestBody defines body for UpdateJob for application/json ContentType.
type UpdateJobJSONRequestBody UpdateJobJSONBody
