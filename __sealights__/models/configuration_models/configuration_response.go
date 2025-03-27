package configuration_models

import (
	"errors"

	"github.com/konflux-ci/release-service/__sealights__/models"
)

type ConfigurationResponse struct {
	DisableAgent                  bool   `json:"disableAgent,omitempty"`
	AgentHeartbeatIntervalSec     int    `json:"agentHeartbeatIntervalSec,omitempty"`
	EnableLogs                    bool   `json:"enableReportLogs,omitempty"`
	LogsIntervalSec               int    `json:"logsIntervalSec,omitempty"`
	LogLevel                      string `json:"logLevel,omitempty"`
	DisableFootprints             bool   `json:"disableFootprints,omitempty"`
	FootprintsIntervalSec         int    `json:"footprintsIntervalSec,omitempty"`
	FootprintsCollectIntervalSecs int    `json:"footprintsCollectIntervalSecs,omitempty"`
	DisableTests                  bool   `json:"disableTests,omitempty"`
	TestsIntervalSec              int    `json:"testsIntervalSec,omitempty"`
	CheckExecutionIntervalSec     int    `json:"checkExecutionIntervalSec,omitempty"`
}

func NewConfigurationResponse() *ConfigurationResponse {
	return &ConfigurationResponse{}
}

func (c *ConfigurationResponse) SetDisableAgent(disableAgent bool) *ConfigurationResponse {
	c.DisableAgent = disableAgent
	return c
}

func (c *ConfigurationResponse) SetAgentHeartbeatIntervalSec(agentHeartbeatIntervalSec int) *ConfigurationResponse {
	c.AgentHeartbeatIntervalSec = agentHeartbeatIntervalSec
	return c
}

func (c *ConfigurationResponse) SetEnableLogs(enableLogs bool) *ConfigurationResponse {
	c.EnableLogs = enableLogs
	return c
}

func (c *ConfigurationResponse) SetLogsIntervalSec(logsIntervalSec int) *ConfigurationResponse {
	c.LogsIntervalSec = logsIntervalSec
	return c
}

func (c *ConfigurationResponse) SetLogLevel(logLevel string) *ConfigurationResponse {
	c.LogLevel = logLevel
	return c
}

func (c *ConfigurationResponse) SetDisableFootprints(disableFootprints bool) *ConfigurationResponse {
	c.DisableFootprints = disableFootprints
	return c
}

func (c *ConfigurationResponse) SetFootprintsIntervalSec(footprintsIntervalSec int) *ConfigurationResponse {
	c.FootprintsIntervalSec = footprintsIntervalSec
	return c
}

func (c *ConfigurationResponse) SetDisableTests(disableTests bool) *ConfigurationResponse {
	c.DisableTests = disableTests
	return c
}

func (c *ConfigurationResponse) SetTestsIntervalSec(testsIntervalSec int) *ConfigurationResponse {
	c.TestsIntervalSec = testsIntervalSec
	return c
}

func (c *ConfigurationResponse) SetFootprintsCollectIntervalSecs(footprintsCollectIntervalSecs int) *ConfigurationResponse {
	c.FootprintsCollectIntervalSecs = footprintsCollectIntervalSecs
	return c
}

func (c *ConfigurationResponse) SetCheckExecutionIntervalSec(checkExecutionIntervalSec int) *ConfigurationResponse {
	c.CheckExecutionIntervalSec = checkExecutionIntervalSec
	return c
}

func (c *ConfigurationResponse) Validate() error {
	if c.AgentHeartbeatIntervalSec < 0 {
		return errors.New("ConfigurationResponse: agentHeartbeatIntervalSec cannot be negative")
	}
	if c.LogsIntervalSec < 0 {
		return errors.New("ConfigurationResponse: logsIntervalSec cannot be negative")
	}
	if c.FootprintsIntervalSec < 0 {
		return errors.New("ConfigurationResponse: footprintsIntervalSec cannot be negative")
	}
	if c.FootprintsCollectIntervalSecs < 0 {
		return errors.New("ConfigurationResponse: footprintsCollectIntervalSecs cannot be negative")
	}

	if c.FootprintsIntervalSec > 0 && c.FootprintsCollectIntervalSecs > 0 && c.FootprintsIntervalSec < c.FootprintsCollectIntervalSecs {
		return errors.New("ConfigurationResponse: footprintsIntervalSec cannot be less than footprintsCollectIntervalSecs")
	}

	if c.TestsIntervalSec < 0 {
		return errors.New("ConfigurationResponse: testsIntervalSec cannot be negative")
	}
	if c.CheckExecutionIntervalSec < 0 {
		return errors.New("ConfigurationResponse: checkExecutionIntervalSec cannot be negative")
	}
	switch c.LogLevel {
	case "", "debug", "info", "warn", "error", "fatal":

	default:
		return errors.New("ConfigurationResponse: logLevel must be one of debug, info, warn, error, fatal")
	}
	return nil
}

var _ models.Model = (*ConfigurationResponse)(nil)
