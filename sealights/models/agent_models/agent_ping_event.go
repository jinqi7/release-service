package agent_models

import (
	"fmt"

	"github.com/konflux-ci/release-service/__sealights__/models"
)

type AgentPingEvent struct {
	Type           int   `json:"type"`
	UtcTimestampMs int64 `json:"utcTimestamp_ms"`
	Data           struct {
		AgentInfo struct {
			LabID string `json:"labId"`
		} `json:"agentInfo"`
		LabID string `json:"labId"`
	} `json:"data"`
}

func NewAgentPingEvent() *AgentPingEvent {
	return &AgentPingEvent{
		Type:           1003,
		UtcTimestampMs: models.UnixMilli(),
	}
}

func (a *AgentPingEvent) SetLabId(labId string) *AgentPingEvent {
	a.Data.AgentInfo.LabID = labId
	a.Data.LabID = labId
	return a
}

func (a *AgentPingEvent) Validate() error {
	if a.Type != 1003 {
		return fmt.Errorf("AgentPingEvent: type is not valid")
	}
	if a.UtcTimestampMs == 0 {
		return fmt.Errorf("AgentPingEvent: utcTimestamp_ms is not valid")
	}
	return nil
}

var _ models.Model = (*AgentPingEvent)(nil)
