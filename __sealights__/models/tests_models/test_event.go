package tests_models

import (
	"fmt"

	"github.com/konflux-ci/release-service/__sealights__/models"
)

type TestEvent struct {
	Type        string `json:"type"`
	Timestamp   int64  `json:"timestamp"`
	ExecutionId string `json:"executionId"`
	TestName    string `json:"testName"`
	Result      string `json:"result,omitempty"`
	Duration    int64  `json:"duration,omitempty"`
}

func NewTestEvent(eventType TestEventType) *TestEvent {
	return &TestEvent{
		Type:        string(eventType),
		Timestamp:   0,
		ExecutionId: "",
		TestName:    "",
		Result:      "",
		Duration:    0,
	}
}

func (t *TestEvent) Validate() error {
	if t.Type == "" {
		return fmt.Errorf("TestEvent: type is empty")
	}
	if t.Timestamp == 0 {
		return fmt.Errorf("TestEvent: timestamp is empty")
	}
	if t.ExecutionId == "" {
		return fmt.Errorf("TestEvent: executionId is empty")
	}
	if t.TestName == "" {
		return fmt.Errorf("TestEvent: testName is empty")
	}
	if t.Type == string(TestEventTypeTestEnd) {
		if t.Result == "" {
			return fmt.Errorf("TestEvent: result is empty for test end event")
		}
		if t.Duration == 0 {
			return fmt.Errorf("TestEvent: duration is empty for test end event")
		}
	}
	return nil
}

func (t *TestEvent) SetExecutionId(executionId string) *TestEvent {
	t.ExecutionId = executionId
	return t
}

func (t *TestEvent) SetTestName(testName string) *TestEvent {
	t.TestName = testName
	return t
}

func (t *TestEvent) SetResult(result TestEventResult) *TestEvent {
	t.Result = string(result)
	return t
}

func (t *TestEvent) SetDuration(duration int64) *TestEvent {
	t.Duration = duration
	return t
}

func (t *TestEvent) SetTimestamp(timestamp int64) *TestEvent {
	t.Timestamp = timestamp
	return t
}

var _ models.Model = (*TestEvent)(nil)
