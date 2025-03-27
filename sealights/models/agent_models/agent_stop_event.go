package agent_models

import (
	"fmt"

	"github.com/konflux-ci/release-service/__sealights__/models"
)

type AgentStopEvent struct {
	Type           int   `json:"type"`
	UtcTimestampMs int64 `json:"utcTimestamp_ms"`
}

func NewAgentStopEvent() *AgentStopEvent {
	return &AgentStopEvent{
		Type:           1002,
		UtcTimestampMs: models.UnixMilli(),
	}
}

func (a *AgentStopEvent) Validate() error {
	if a.Type != 1002 {
		return fmt.Errorf("AgentStopEvent: type is not valid")
	}
	if a.UtcTimestampMs == 0 {
		return fmt.Errorf("AgentStopEvent: utcTimestamp_ms is not valid")
	}

	return nil
}

var _ models.Model = (*AgentStopEvent)(nil)
