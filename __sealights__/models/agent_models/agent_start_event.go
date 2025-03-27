package agent_models

import (
	"fmt"

	"github.com/konflux-ci/release-service/__sealights__/models"
)

type AgentStartEvent struct {
	Type           int            `json:"type"`
	UtcTimestampMs int64          `json:"utcTimestamp_ms"`
	Data           *AgentMetadata `json:"data"`
}

func NewAgentStartEvent() *AgentStartEvent {
	return &AgentStartEvent{
		Type:           1001,
		UtcTimestampMs: models.UnixMilli(),
		Data:           NewAgentMetadata(),
	}
}

func (a *AgentStartEvent) Validate() error {
	if a.Type != 1001 {
		return fmt.Errorf("AgentStartEvent: type is not valid")
	}
	if a.UtcTimestampMs == 0 {
		return fmt.Errorf("AgentStartEvent: utcTimestamp_ms is not valid")
	}
	if a.Data == nil {
		return fmt.Errorf("AgentStartEvent: data is not valid")
	}
	if err := a.Data.Validate(); err != nil {
		return err
	}
	return nil
}

func (a *AgentStartEvent) SetTechnology(value string) *AgentStartEvent {
	a.Data.AgentInfo.Technology = value
	return a
}

func (a *AgentStartEvent) SetAgentType(value AgentType) *AgentStartEvent {
	a.Data.AgentInfo.SetAgentType(value)
	return a
}

func (a *AgentStartEvent) SetAgentVersion(value string) *AgentStartEvent {
	a.Data.AgentInfo.SetAgentVersion(value)
	return a
}

func (a *AgentStartEvent) SetLabId(value string) *AgentStartEvent {
	a.Data.AgentInfo.SetLabId(value)
	return a
}

func (a *AgentStartEvent) SetTestStage(value string) *AgentStartEvent {
	a.Data.AgentInfo.SetTestStage(value)
	return a
}

func (a *AgentStartEvent) AddAgentConfigKeyVal(key, val string) *AgentStartEvent {
	a.Data.AgentInfo.SetAgentConfigKeyVal(key, val)
	return a
}

func (a *AgentStartEvent) AddArgs(value string) *AgentStartEvent {
	a.Data.AgentInfo.Argv = append(a.Data.AgentInfo.Argv, value)
	return a
}

func (a *AgentStartEvent) SetBuildSessionId(value string) *AgentStartEvent {
	a.Data.AgentInfo.SetBuildSessionId(value)
	return a
}

var _ models.Model = (*MachineInfo)(nil)
