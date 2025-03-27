package execution_models

import (
	"fmt"

	"github.com/konflux-ci/release-service/__sealights__/models"
)

type ExecutionRequest struct {
	LabId string `json:"labId"`
}

func NewExecutionRequest() *ExecutionRequest {
	return &ExecutionRequest{}
}

func (e *ExecutionRequest) Validate() error {
	if e.LabId == "" {
		return fmt.Errorf("ExecutionRequest: labId is required")
	}
	return nil
}

func (e *ExecutionRequest) SetLabId(labId string) *ExecutionRequest {
	e.LabId = labId
	return e
}

var _ models.Model = (*ExecutionRequest)(nil)
