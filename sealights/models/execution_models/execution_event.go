package execution_models

import (
	"fmt"

	"github.com/konflux-ci/release-service/__sealights__/models"
)

type ExecutionEvent struct {
	AppName     string `json:"appName"`
	BuildName   string `json:"buildName"`
	BranchName  string `json:"branchName"`
	LabId       string `json:"labId"`
	TestStage   string `json:"testStage"`
	ExecutionId string `json:"executionId"`
	AgentId     string `json:"agentId"`
	TestGroupId string `json:"testGroupId"`
}

func NewExecutionEvent() *ExecutionEvent {
	return &ExecutionEvent{}
}

func (e *ExecutionEvent) Validate() error {
	if e.AppName == "" {
		return fmt.Errorf("ExecutionEvent: appName is required")
	}
	if e.BuildName == "" {
		return fmt.Errorf("ExecutionEvent: buildName is required")
	}
	if e.BranchName == "" {
		return fmt.Errorf("ExecutionEvent: branchName is required")
	}
	if e.LabId == "" {
		return fmt.Errorf("ExecutionEvent: labId is required")
	}
	if e.TestStage == "" {
		return fmt.Errorf("ExecutionEvent: testStage is required")
	}
	if e.ExecutionId == "" {
		return fmt.Errorf("ExecutionEvent: executionId is required")
	}
	if e.AgentId == "" {
		return fmt.Errorf("ExecutionEvent: agentId is required")
	}

	return nil
}

func (e *ExecutionEvent) SetAppName(appName string) *ExecutionEvent {
	e.AppName = appName
	return e
}

func (e *ExecutionEvent) SetLabId(labId string) *ExecutionEvent {
	e.LabId = labId
	return e
}

func (e *ExecutionEvent) SetTestStage(testStage string) *ExecutionEvent {
	e.TestStage = testStage
	return e
}

func (e *ExecutionEvent) SetExecutionId(executionId string) *ExecutionEvent {
	e.ExecutionId = executionId
	return e
}

func (e *ExecutionEvent) SetAgentId(agentId string) *ExecutionEvent {
	e.AgentId = agentId
	return e
}

func (e *ExecutionEvent) SetBuildName(buildName string) *ExecutionEvent {
	e.BuildName = buildName
	return e
}

func (e *ExecutionEvent) SetBranchName(branchName string) *ExecutionEvent {
	e.BranchName = branchName
	return e
}
func (e *ExecutionEvent) SetTestGroupId(testGroupId string) *ExecutionEvent {
	e.TestGroupId = testGroupId
	return e
}

var _ models.Model = (*ExecutionEvent)(nil)
