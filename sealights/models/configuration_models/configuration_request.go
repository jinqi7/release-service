package configuration_models

import "github.com/konflux-ci/release-service/__sealights__/models"

type ConfigurationRequest struct {
	AppName      string `json:"appName"`
	BranchName   string `json:"branchName"`
	BuildName    string `json:"buildName"`
	TestStage    string `json:"testStage"`
	LabId        string `json:"labId"`
	AgentId      string `json:"agentId"`
	AgentType    string `json:"agentType"`
	AgentVersion string `json:"agentVersion"`
}

func NewConfigurationRequest() *ConfigurationRequest {
	return &ConfigurationRequest{}
}

func (c *ConfigurationRequest) SetAppName(appName string) *ConfigurationRequest {
	c.AppName = appName
	return c
}

func (c *ConfigurationRequest) SetBranchName(branchName string) *ConfigurationRequest {
	c.BranchName = branchName
	return c
}

func (c *ConfigurationRequest) SetBuildName(buildName string) *ConfigurationRequest {
	c.BuildName = buildName
	return c
}

func (c *ConfigurationRequest) SetTestStage(testStage string) *ConfigurationRequest {
	c.TestStage = testStage
	return c
}

func (c *ConfigurationRequest) SetLabId(labId string) *ConfigurationRequest {
	c.LabId = labId
	return c
}

func (c *ConfigurationRequest) SetAgentId(agentId string) *ConfigurationRequest {
	c.AgentId = agentId
	return c
}

func (c *ConfigurationRequest) SetAgentType(agentType string) *ConfigurationRequest {
	c.AgentType = agentType
	return c
}

func (c *ConfigurationRequest) SetAgentVersion(agentVersion string) *ConfigurationRequest {
	c.AgentVersion = agentVersion
	return c
}

func (c *ConfigurationRequest) Validate() error {
	return nil
}

func (c *ConfigurationRequest) ToQueryMap() map[string]string {
	queryMap := map[string]string{}
	if c.AppName != "" {
		queryMap["appName"] = c.AppName
	}
	if c.BranchName != "" {
		queryMap["branchName"] = c.BranchName
	}
	if c.BuildName != "" {
		queryMap["buildName"] = c.BuildName
	}
	if c.TestStage != "" {
		queryMap["testStage"] = c.TestStage
	}
	if c.LabId != "" {
		queryMap["labId"] = c.LabId
	}
	if c.AgentId != "" {
		queryMap["agentId"] = c.AgentId
	}
	if c.AgentType != "" {
		queryMap["agentType"] = c.AgentType
	}
	if c.AgentVersion != "" {
		queryMap["agentVersion"] = c.AgentVersion
	}
	return queryMap
}

var _ models.Model = (*ConfigurationRequest)(nil)
