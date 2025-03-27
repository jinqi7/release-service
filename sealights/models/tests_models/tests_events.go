package tests_models

import (
	"fmt"

	"github.com/konflux-ci/release-service/__sealights__/models"
)

type (
	TestEventType   string
	TestEventResult string
)

const (
	TestEventTypeTestStart TestEventType = "testStart"
	TestEventTypeTestEnd   TestEventType = "testEnd"

	TestEventResultPassed  TestEventResult = "passed"
	TestEventResultFailed  TestEventResult = "failed"
	TestEventResultSkipped TestEventResult = "skipped"
)

type TestsEvents struct {
	AppName             string `json:"appName"`
	Branch              string `json:"branch"`
	Build               string `json:"build"`
	TestSelectionStatus string `json:"testSelectionStatus"`
	Environment         struct {
		TestStage string `json:"testStage"`
		LabId     string `json:"labId"`
		AgentId   string `json:"agentId"`
	} `json:"environment"`
	Events []*TestEvent `json:"events"`
}

func NewTestsEvents() *TestsEvents {
	te := &TestsEvents{
		AppName: "",
		Branch:  "",
		Build:   "",
		Environment: struct {
			TestStage string `json:"testStage"`
			LabId     string `json:"labId"`
			AgentId   string `json:"agentId"`
		}{
			TestStage: "",
			LabId:     "",
			AgentId:   "",
		},
		Events: nil,
	}
	return te
}

func (t *TestsEvents) SetAppName(appName string) *TestsEvents {
	t.AppName = appName
	return t
}

func (t *TestsEvents) SetBranch(branch string) *TestsEvents {
	t.Branch = branch
	return t
}

func (t *TestsEvents) SetBuild(build string) *TestsEvents {
	t.Build = build
	return t
}

func (t *TestsEvents) SetTestStage(testStage string) *TestsEvents {
	t.Environment.TestStage = testStage
	return t
}

func (t *TestsEvents) SetLabId(labId string, buildSessionId string) *TestsEvents {
	if labId != "" {
		t.Environment.LabId = labId
	} else {
		t.Environment.LabId = buildSessionId
	}

	return t
}

func (t *TestsEvents) SetAgentId(agentId string) *TestsEvents {
	t.Environment.AgentId = agentId
	return t
}

func (t *TestsEvents) SetTestSelectionStatus(testSelectionStatus string) *TestsEvents {
	t.TestSelectionStatus = testSelectionStatus
	return t
}

func (t *TestsEvents) Validate() error {
	if t.AppName == "" {
		return fmt.Errorf("TestsEvents: appName is empty")
	}
	if t.Branch == "" {
		return fmt.Errorf("TestsEvents: branch is empty")
	}
	if t.Build == "" {
		return fmt.Errorf("TestsEvents: build is empty")
	}
	if t.Environment.LabId == "" {
		return fmt.Errorf("TestsEvents: labId is empty")
	}
	if t.Environment.AgentId == "" {
		return fmt.Errorf("TestsEvents: agentId is empty")
	}
	if t.Environment.TestStage == "" {
		return fmt.Errorf("TestsEvents: testStage is empty")
	}
	if len(t.Events) == 0 {
		return fmt.Errorf("TestsEvents: events is empty")
	}
	for _, event := range t.Events {
		if err := event.Validate(); err != nil {
			return err
		}
	}
	return nil
}

func (t *TestsEvents) AddEvent(event *TestEvent) *TestsEvents {
	t.Events = append(t.Events, event)
	return t
}

var _ models.Model = (*TestsEvents)(nil)
