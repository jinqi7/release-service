package execution_models

import (
	"encoding/json"
	"fmt"

	"github.com/konflux-ci/release-service/__sealights__/models"
)

type ExecutionState int

const (
	ExecutionStateUnknown ExecutionState = iota
	ExecutionStateNoOpenedExecution
	ExecutionStateRunning
	ExecutionStateFinished
)

type ExecutionResponse struct {
	Execution *struct {
		ExecutionId    string `json:"executionId"`
		AppName        string `json:"appName"`
		BranchName     string `json:"branchName"`
		BuildName      string `json:"buildName"`
		BuildSessionId string `json:"buildSessionId"`
		TestStage      string `json:"testStage"`
		Timestamp      int64  `json:"timestamp"`
		Status         string `json:"status"`
		TestGroupId    string `json:"testGroupId"`
	} `json:"execution"`
}

func NewExecutionResponse() *ExecutionResponse {
	return &ExecutionResponse{
		Execution: &struct {
			ExecutionId    string `json:"executionId"`
			AppName        string `json:"appName"`
			BranchName     string `json:"branchName"`
			BuildName      string `json:"buildName"`
			BuildSessionId string `json:"buildSessionId"`
			TestStage      string `json:"testStage"`
			Timestamp      int64  `json:"timestamp"`
			Status         string `json:"status"`
			TestGroupId    string `json:"testGroupId"`
		}{},
	}
}

func (e *ExecutionResponse) SetExecutionId(executionId string) *ExecutionResponse {
	e.Execution.ExecutionId = executionId
	return e
}

func (e *ExecutionResponse) SetAppName(appName string) *ExecutionResponse {
	e.Execution.AppName = appName
	return e
}

func (e *ExecutionResponse) SetBranchName(branchName string) *ExecutionResponse {
	e.Execution.BranchName = branchName
	return e
}

func (e *ExecutionResponse) SetBuildName(buildName string) *ExecutionResponse {
	e.Execution.BuildName = buildName
	return e
}

func (e *ExecutionResponse) SetBuildSessionId(buildSessionId string) *ExecutionResponse {
	e.Execution.BuildSessionId = buildSessionId
	return e
}

func (e *ExecutionResponse) SetTestStage(testStage string) *ExecutionResponse {
	e.Execution.TestStage = testStage
	return e
}

func (e *ExecutionResponse) SetTimestamp(timestamp int64) *ExecutionResponse {
	e.Execution.Timestamp = timestamp
	return e
}

func (e *ExecutionResponse) SetStatus(status string) *ExecutionResponse {
	e.Execution.Status = status
	return e
}
func (e *ExecutionResponse) SetTestGroupId(testGroupId string) *ExecutionResponse {
	e.Execution.TestGroupId = testGroupId
	return e
}

func (e *ExecutionResponse) GetExecutionId() (string, ExecutionState) {
	if e.Execution == nil {
		return "", ExecutionStateNoOpenedExecution
	}
	switch e.Execution.Status {
	case "":
		return "", ExecutionStateUnknown
	case "created":
		return e.Execution.ExecutionId, ExecutionStateRunning
	case "ended", "pendingDelete":
		return "", ExecutionStateFinished
	default:
		return "", ExecutionStateUnknown
	}
}

func (e *ExecutionResponse) Json() string {
	data, _ := json.MarshalIndent(e, "", "  ")
	return string(data)
}

func (e *ExecutionResponse) Validate() error {
	if e.Execution == nil {
		return nil
	}
	if e.Execution.ExecutionId == "" {
		return fmt.Errorf("ExecutionResponse: executionId is required")
	}
	if e.Execution.Status == "" {
		return fmt.Errorf("ExecutionResponse: status is required")
	}
	return nil
}

var _ models.Model = (*ExecutionResponse)(nil)
