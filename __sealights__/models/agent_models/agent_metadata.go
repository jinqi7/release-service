package agent_models

import "github.com/konflux-ci/release-service/__sealights__/models"

type AgentMetadata struct {
	TechSpecificInfo *TechSpecificInfo `json:"techSpecificInfo,omitempty"`
	AgentInfo        *AgentInfo        `json:"agentInfo,omitempty"`
	MachineInfo      *MachineInfo      `json:"machineInfo,omitempty"`
}

func NewAgentMetadata() *AgentMetadata {
	a := &AgentMetadata{
		TechSpecificInfo: NewTechSpecificInfo(),
		AgentInfo:        NewAgentInfo(),
		MachineInfo:      NewMachineInfo(),
	}

	return a
}

func (a *AgentMetadata) Validate() error {
	if err := a.AgentInfo.Validate(); err != nil {
		return err
	}
	if err := a.MachineInfo.Validate(); err != nil {
		return err
	}
	if err := a.TechSpecificInfo.Validate(); err != nil {
		return err
	}
	return nil
}

var _ models.Model = (*AgentMetadata)(nil)
