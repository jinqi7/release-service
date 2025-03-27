package agent_models

import (
	"fmt"
	"os"
	"strings"

	"github.com/konflux-ci/release-service/__sealights__/models"
)

type AgentType string

const (
	AgentTypeTestListener     AgentType = "TestListener"
	AgentTypeBuildScanner     AgentType = "BuildScanner"
	AgentTypeCollectorService AgentType = "CollectorService"
	AgentTypeCLI              AgentType = "CLI"
	AgentTypeRunner           AgentType = "TestsRunner"
)

type AgentInfo struct {
	Technology     string            `json:"technology"`
	AgentType      string            `json:"agentType"`
	AgentVersion   string            `json:"agentVersion"`
	LabId          string            `json:"labId"`
	TestStage      string            `json:"testStage"`
	Tags           TagsInfo          `json:"tags,omitempty"`
	AgentConfig    map[string]string `json:"agentConfig"`
	Argv           []string          `json:"argv"`
	EnvVars        map[string]string `json:"envVars"`
	SendsPing      bool              `json:"sendsPing"`
	BuildSessionId string            `json:"buildSessionId"`
}

func NewAgentInfo() *AgentInfo {
	ai := &AgentInfo{
		Technology:   "golang",
		AgentType:    "",
		AgentVersion: "",
		LabId:        "",
		TestStage:    "",
		Tags:         createTagsFromEnv(),
		AgentConfig:  map[string]string{},
		Argv:         nil,
		EnvVars:      getEnvVarToMap(),
		SendsPing:    true,
	}
	ai.Argv = append(ai.Argv, os.Args...)
	return ai
}

func (a *AgentInfo) SetTechnology(technology string) *AgentInfo {
	a.Technology = technology
	return a
}

func (a *AgentInfo) SetAgentType(agentType AgentType) *AgentInfo {
	a.AgentType = string(agentType)
	return a
}

func (a *AgentInfo) SetAgentVersion(agentVersion string) *AgentInfo {
	a.AgentVersion = agentVersion
	return a
}

func (a *AgentInfo) SetLabId(labId string) *AgentInfo {
	a.LabId = labId
	return a
}

func (a *AgentInfo) SetTags(tags TagsInfo) *AgentInfo {
	a.Tags = tags
	return a
}

func (a *AgentInfo) SetAgentConfig(agentConfig map[string]string) *AgentInfo {
	a.AgentConfig = agentConfig
	return a
}

func (a *AgentInfo) SetAgentConfigKeyVal(key, val string) *AgentInfo {
	a.AgentConfig[key] = val
	return a
}

func (a *AgentInfo) SetArgv(argv []string) *AgentInfo {
	a.Argv = argv
	return a
}

func (a *AgentInfo) SetTestStage(testStage string) *AgentInfo {
	a.TestStage = testStage
	return a
}

func (a *AgentInfo) SetSendPing(sendPing bool) *AgentInfo {
	a.SendsPing = sendPing
	return a
}

func (a *AgentInfo) SetBuildSessionId(buildSessionId string) *AgentInfo {
	a.BuildSessionId = buildSessionId
	return a
}

func (a *AgentInfo) Validate() error {
	if a.Technology == "" {
		return fmt.Errorf("AgentInfo: technology is not valid")
	}
	if a.AgentType == "" {
		return fmt.Errorf("AgentInfo: agentType is not valid")
	}
	if a.AgentVersion == "" {
		return fmt.Errorf("AgentInfo: agentVersion is not valid")
	}

	return nil
}

type TagInfo struct {
	Name string `json:"name"`
}
type TagsInfo []*TagInfo

func createTagsFromEnv() TagsInfo {
	tags := TagsInfo{}
	envTags := os.Getenv("SEALIGHTS_CLIENT_TAGS")
	if envTags != "" {
		for _, tag := range strings.Split(envTags, ",") {
			tags = append(tags, &TagInfo{Name: tag})
		}
	}
	if len(tags) == 0 {
		return nil
	}
	return tags
}

var _ models.Model = &AgentInfo{}
