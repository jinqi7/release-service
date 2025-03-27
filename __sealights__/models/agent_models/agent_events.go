package agent_models

import (
	"fmt"

	"github.com/konflux-ci/release-service/__sealights__/models"
)

type AgentEvents struct {
	Events         []models.Model `json:"events"`
	AgentId        string         `json:"agentId"`
	BuildSessionId string         `json:"buildSessionId,omitempty"`
	AppName        string         `json:"appName,omitempty"`
}

func NewAgentEvents() *AgentEvents {
	return &AgentEvents{}
}

func (a *AgentEvents) Validate() error {
	if a.AgentId == "" {
		return fmt.Errorf("AgentEvents: agentId is required")
	}
	if len(a.Events) == 0 {
		return fmt.Errorf("AgentEvents: at list one event is required")
	} else {
		for _, event := range a.Events {
			if err := event.Validate(); err != nil {
				return err
			}
		}
	}
	return nil
}

func (a *AgentEvents) AddEvent(events ...models.Model) *AgentEvents {
	a.Events = append(a.Events, events...)
	return a
}

func (a *AgentEvents) SetAgentId(id string) *AgentEvents {
	a.AgentId = id
	return a
}

func (a *AgentEvents) SetBuildSessionId(id string) *AgentEvents {
	a.BuildSessionId = id
	return a
}

func (a *AgentEvents) SetAppName(name string) *AgentEvents {
	a.AppName = name
	return a
}

var _ models.Model = (*AgentEvents)(nil)
