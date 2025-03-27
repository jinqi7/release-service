package tests_models

import (
	"fmt"

	"github.com/konflux-ci/release-service/__sealights__/models"
)

type TestErrorEvent struct {
	Type           int    `json:"type"`
	UtcTimestampMs int64  `json:"utcTimestamp_ms"`
	Data           string `json:"data"`
}

func NewTestErrorEventEvent() *TestErrorEvent {
	return &TestErrorEvent{
		Type:           4007,
		UtcTimestampMs: models.UnixMilli(),
	}
}

func (t *TestErrorEvent) SetData(value string) *TestErrorEvent {
	t.Data = value
	return t
}

func (t *TestErrorEvent) Validate() error {
	if t.Type != 4007 {
		return fmt.Errorf("TestErrorEvent: type is not valid")
	}
	if t.UtcTimestampMs == 0 {
		return fmt.Errorf("TestErrorEvent: utcTimestamp_ms is not valid")
	}
	if t.Data == "" {
		return fmt.Errorf("TestErrorEvent: data is not valid")
	}
	return nil
}

var _ models.Model = (*TestErrorEvent)(nil)
