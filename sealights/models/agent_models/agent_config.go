package agent_models

import (
	"fmt"
	"sync"
	"time"

	"github.com/konflux-ci/release-service/__sealights__/models"
)

type AgentConfig struct {
	sync.RWMutex
	AppName                         string `json:"appName"`
	ModuleName                      string `json:"moduleName"`
	Build                           string `json:"build"`
	Branch                          string `json:"branch"`
	AgentId                         string `json:"agentId"`
	BuildSessionId                  string `json:"buildSessionId"`
	LabId                           string `json:"labId"`
	ServerUrl                       string `json:"serverUrl"`
	ProxyUrl                        string `json:"proxyUrl"`
	TestStage                       string `json:"testStage"`
	LogLevel                        string `json:"logLevel"`
	AgentVersion                    string `json:"agentVersion"`
	DisableAgent                    bool   `json:"disableAgent,omitempty"`
	AgentHeartbeatIntervalSec       int    `json:"agentHeartbeatIntervalSec,omitempty"`
	EnableLogs                      bool   `json:"EnableReportLogs,omitempty"`
	LogsIntervalSec                 int    `json:"logsIntervalSec,omitempty"`
	DisableFootprints               bool   `json:"disableFootprints,omitempty"`
	FootprintsReportingIntervalSec  int    `json:"footprintsReportingIntervalSec,omitempty"`
	FootprintsCollectionIntervalSec int    `json:"footprintsCollectionIntervalSec,omitempty"`
	DisableTests                    bool   `json:"disableTests,omitempty"`
	TestsIntervalSec                int    `json:"testsIntervalSec,omitempty"`
	CheckExecutionIntervalSec       int    `json:"checkExecutionIntervalSec,omitempty"`
	ExecutionId                     string `json:"executionId,omitempty"`
	ExecutionOpenTime               int64  `json:"executionOpenTime,omitempty"`
	EnableRemoteConfig              bool   `json:"enableRemoteConfig,omitempty"`
	RemoteConfigIntervalSec         int64  `json:"remoteConfigIntervalSec,omitempty"`
	TestSelection                   bool   `json:"testSelection,omitempty"`
	ExecutionBuildSessionId         string `json:"executionBuildSessionId,omitempty"`
	TestGroupId                     string `json:"testGroupId,omitempty"`
	TestRecommendationSleepSeconds  int64  `json:"testRecommendationSleepSeconds,omitempty"`
	TestsRunnerMode                 bool   `json:"testsRunnerMode,omitempty"`
	GinkgoEnabled                   bool   `json:"ginkgoEnabled,omitempty"`
}

func NewAgentConfig() *AgentConfig {
	return &AgentConfig{
		RWMutex:                         sync.RWMutex{},
		AppName:                         "",
		ModuleName:                      "",
		Build:                           "",
		Branch:                          "",
		AgentId:                         "",
		BuildSessionId:                  "",
		LabId:                           "",
		ServerUrl:                       "",
		TestStage:                       "",
		LogLevel:                        "",
		AgentVersion:                    "",
		DisableAgent:                    false,
		AgentHeartbeatIntervalSec:       120,
		EnableLogs:                      false,
		LogsIntervalSec:                 60,
		DisableFootprints:               false,
		FootprintsReportingIntervalSec:  5,
		FootprintsCollectionIntervalSec: 1,
		DisableTests:                    false,
		TestsIntervalSec:                1,
		CheckExecutionIntervalSec:       5,
		ExecutionId:                     "",
		ExecutionOpenTime:               0,
		EnableRemoteConfig:              false,
		RemoteConfigIntervalSec:         300,
		TestSelection:                   true,
		ExecutionBuildSessionId:         "",
		TestGroupId:                     "",
		TestRecommendationSleepSeconds:  60,
		TestsRunnerMode:                 false,
		GinkgoEnabled:                   false,
	}
}

func (a *AgentConfig) SetAppName(v string) *AgentConfig {
	a.Lock()
	defer a.Unlock()
	a.AppName = v
	return a
}

func (a *AgentConfig) GetAppName() string {
	a.RLock()
	defer a.RUnlock()
	return a.AppName
}

func (a *AgentConfig) SetModuleName(v string) *AgentConfig {
	a.Lock()
	defer a.Unlock()
	a.ModuleName = v
	return a
}

func (a *AgentConfig) GetModuleName() string {
	a.RLock()
	defer a.RUnlock()
	return a.ModuleName
}

func (a *AgentConfig) SetBuild(v string) *AgentConfig {
	a.Lock()
	defer a.Unlock()
	a.Build = v
	return a
}

func (a *AgentConfig) GetBuild() string {
	a.RLock()
	defer a.RUnlock()
	return a.Build
}

func (a *AgentConfig) SetBranch(v string) *AgentConfig {
	a.Lock()
	defer a.Unlock()
	a.Branch = v
	return a
}

func (a *AgentConfig) GetBranch() string {
	a.RLock()
	defer a.RUnlock()
	return a.Branch
}

func (a *AgentConfig) SetAgentId(v string) *AgentConfig {
	a.Lock()
	defer a.Unlock()
	a.AgentId = v
	return a
}

func (a *AgentConfig) GetAgentId() string {
	a.RLock()
	defer a.RUnlock()
	return a.AgentId
}

func (a *AgentConfig) SetBuildSessionId(v string) *AgentConfig {
	a.Lock()
	defer a.Unlock()
	a.BuildSessionId = v
	return a
}

func (a *AgentConfig) GetBuildSessionId() string {
	a.RLock()
	defer a.RUnlock()
	return a.BuildSessionId
}

func (a *AgentConfig) SetLabId(v string) *AgentConfig {
	a.Lock()
	defer a.Unlock()
	a.LabId = v
	return a
}

func (a *AgentConfig) GetLabId() string {
	a.RLock()
	defer a.RUnlock()
	return a.LabId
}

func (a *AgentConfig) SetTestStage(v string) *AgentConfig {
	a.Lock()
	defer a.Unlock()
	a.TestStage = v
	return a
}

func (a *AgentConfig) GetTestStage() string {
	a.RLock()
	defer a.RUnlock()
	return a.TestStage
}

func (a *AgentConfig) SetLogLevel(v string) *AgentConfig {
	a.Lock()
	defer a.Unlock()
	a.LogLevel = v
	return a
}

func (a *AgentConfig) GetLogLevel() string {
	a.RLock()
	defer a.RUnlock()
	return a.LogLevel
}

func (a *AgentConfig) SetAgentVersion(v string) *AgentConfig {
	a.Lock()
	defer a.Unlock()
	a.AgentVersion = v
	return a
}

func (a *AgentConfig) GetAgentVersion() string {
	a.RLock()
	defer a.RUnlock()
	return a.AgentVersion
}

func (a *AgentConfig) SetDisableAgent(v bool) *AgentConfig {
	a.Lock()
	defer a.Unlock()
	a.DisableAgent = v
	return a
}

func (a *AgentConfig) GetDisableAgent() bool {
	a.RLock()
	defer a.RUnlock()
	return a.DisableAgent
}

func (a *AgentConfig) SetAgentHeartbeatIntervalSec(v int) *AgentConfig {
	a.Lock()
	defer a.Unlock()
	a.AgentHeartbeatIntervalSec = v
	return a
}

func (a *AgentConfig) GetAgentHeartbeatIntervalDuration() time.Duration {
	a.RLock()
	defer a.RUnlock()
	return time.Duration(a.AgentHeartbeatIntervalSec) * time.Second
}

func (a *AgentConfig) SetServerUrl(v string) *AgentConfig {
	a.Lock()
	defer a.Unlock()
	a.ServerUrl = v
	return a
}

func (a *AgentConfig) GetServerUrl() string {
	a.RLock()
	defer a.RUnlock()
	return a.ServerUrl
}

func (a *AgentConfig) SetProxyUrl(v string) *AgentConfig {
	a.Lock()
	defer a.Unlock()
	a.ProxyUrl = v
	return a
}

func (a *AgentConfig) GetProxyUrl() string {
	a.RLock()
	defer a.RUnlock()
	return a.ProxyUrl
}

func (a *AgentConfig) SetEnableLogs(v bool) *AgentConfig {
	a.Lock()
	defer a.Unlock()
	a.EnableLogs = v
	return a
}

func (a *AgentConfig) GetEnableLogs() bool {
	a.RLock()
	defer a.RUnlock()
	return a.EnableLogs
}

func (a *AgentConfig) SetLogsIntervalSec(v int) *AgentConfig {
	a.Lock()
	defer a.Unlock()
	a.LogsIntervalSec = v
	return a
}

func (a *AgentConfig) GetLogsIntervalDuration() time.Duration {
	a.RLock()
	defer a.RUnlock()
	return time.Duration(a.LogsIntervalSec) * time.Second
}

func (a *AgentConfig) SetDisableFootprints(v bool) *AgentConfig {
	a.Lock()
	defer a.Unlock()
	a.DisableFootprints = v
	return a
}

func (a *AgentConfig) SetTestsRunnerMode(v bool) *AgentConfig {
	a.Lock()
	defer a.Unlock()
	a.TestsRunnerMode = v
	return a
}
func (a *AgentConfig) GetDisableFootprints() bool {
	a.RLock()
	defer a.RUnlock()
	return a.DisableFootprints
}

func (a *AgentConfig) SetFootprintsReportingIntervalSec(v int) *AgentConfig {
	a.Lock()
	defer a.Unlock()
	a.FootprintsReportingIntervalSec = v
	return a
}

func (a *AgentConfig) GetFootprintsReportingIntervalDuration() time.Duration {
	a.RLock()
	defer a.RUnlock()
	return time.Duration(a.FootprintsReportingIntervalSec) * time.Second
}

func (a *AgentConfig) SetFootprintsCollectionInterval(v int) *AgentConfig {
	a.Lock()
	defer a.Unlock()
	a.FootprintsCollectionIntervalSec = v
	return a
}

func (a *AgentConfig) GetFootprintsCollectionIntervalDuration() time.Duration {
	a.RLock()
	defer a.RUnlock()
	return time.Duration(a.FootprintsCollectionIntervalSec) * time.Second
}

func (a *AgentConfig) SetDisableTests(v bool) *AgentConfig {
	a.Lock()
	defer a.Unlock()
	a.DisableTests = v
	return a
}

func (a *AgentConfig) GetDisableTests() bool {
	a.RLock()
	defer a.RUnlock()
	return a.DisableTests
}

func (a *AgentConfig) SetTestsIntervalSec(v int) *AgentConfig {
	a.Lock()
	defer a.Unlock()
	a.TestsIntervalSec = v
	return a
}

func (a *AgentConfig) GetTestsIntervalDuration() time.Duration {
	a.RLock()
	defer a.RUnlock()
	return time.Duration(a.TestsIntervalSec) * time.Second
}

func (a *AgentConfig) SetCheckExecutionIntervalSec(v int) *AgentConfig {
	a.Lock()
	defer a.Unlock()
	a.CheckExecutionIntervalSec = v
	return a
}

func (a *AgentConfig) GetCheckExecutionIntervalDuration() time.Duration {
	a.RLock()
	defer a.RUnlock()
	return time.Duration(a.CheckExecutionIntervalSec) * time.Second
}

func (a *AgentConfig) SetExecutionId(v string) *AgentConfig {
	a.Lock()
	defer a.Unlock()
	a.ExecutionId = v
	return a
}

func (a *AgentConfig) GetExecutionId() string {
	a.RLock()
	defer a.RUnlock()
	return a.ExecutionId
}

func (a *AgentConfig) SetExecutionOpenTime(v int64) *AgentConfig {
	a.Lock()
	defer a.Unlock()
	a.ExecutionOpenTime = v
	return a
}

func (a *AgentConfig) GetExecutionOpenTime() int64 {
	a.RLock()
	defer a.RUnlock()
	return a.ExecutionOpenTime
}

func (a *AgentConfig) SetEnableRemoteConfig(v bool) *AgentConfig {
	a.Lock()
	defer a.Unlock()
	a.EnableRemoteConfig = v
	return a
}

func (a *AgentConfig) GetEnableRemoteConfig() bool {
	a.RLock()
	defer a.RUnlock()
	return a.EnableRemoteConfig
}

func (a *AgentConfig) SetRemoteConfigIntervalSec(v int64) *AgentConfig {
	a.Lock()
	defer a.Unlock()
	a.RemoteConfigIntervalSec = v
	return a
}

func (a *AgentConfig) GetRemoteConfigIntervalDuration() time.Duration {
	a.RLock()
	defer a.RUnlock()
	return time.Duration(a.RemoteConfigIntervalSec) * time.Second
}

func (a *AgentConfig) SetTestSelection(v bool) *AgentConfig {
	a.Lock()
	defer a.Unlock()
	a.TestSelection = v
	return a
}

func (a *AgentConfig) GetTestSelection() bool {
	a.RLock()
	defer a.RUnlock()
	return a.TestSelection
}
func (a *AgentConfig) SetTestGroupId(v string) *AgentConfig {
	a.Lock()
	defer a.Unlock()
	a.TestGroupId = v
	return a
}

func (a *AgentConfig) GetTestGroupId() string {
	a.RLock()
	defer a.RUnlock()
	return a.TestGroupId
}
func (a *AgentConfig) SetExecutionBuildSessionId(v string) *AgentConfig {
	a.Lock()
	defer a.Unlock()
	a.ExecutionBuildSessionId = v
	return a
}

func (a *AgentConfig) GetExecutionBuildSessionId() string {
	a.RLock()
	defer a.RUnlock()
	return a.ExecutionBuildSessionId
}
func (a *AgentConfig) SetTestRecommendationSleepSeconds(v int64) *AgentConfig {
	a.Lock()
	defer a.Unlock()
	a.TestRecommendationSleepSeconds = v
	return a
}

func (a *AgentConfig) GetTestRecommendationSleepSeconds() int64 {
	a.RLock()
	defer a.RUnlock()
	return a.TestRecommendationSleepSeconds
}

func (a *AgentConfig) GetTestsRunnerMode() bool {
	a.RLock()
	defer a.RUnlock()
	return a.TestsRunnerMode
}

func (a *AgentConfig) SetGinkgoEnabled(v bool) *AgentConfig {
	a.Lock()
	defer a.Unlock()
	a.GinkgoEnabled = v
	return a
}

func (a *AgentConfig) GetGinkgoEnabled() bool {
	a.RLock()
	defer a.RUnlock()
	return a.GinkgoEnabled
}
func (a *AgentConfig) Validate() error {
	a.RLock()
	defer a.RUnlock()
	if a.AppName == "" {
		return fmt.Errorf("AgentConfig:AppName is required")
	}
	if a.AgentId == "" {
		return fmt.Errorf("AgentConfig: agentId is required")
	}
	if a.Build == "" {
		return fmt.Errorf("AgentConfig: buildId is required")
	}
	if a.Branch == "" {
		return fmt.Errorf("AgentConfig: branch is required")
	}
	if a.BuildSessionId == "" {
		return fmt.Errorf("AgentConfig: buildSessionId is required")
	}
	if a.AgentVersion == "" {
		return fmt.Errorf("AgentConfig: agentVersion is required")
	}
	if a.ServerUrl == "" {
		return fmt.Errorf("AgentConfig: serverUrl is required")
	}

	if a.AgentHeartbeatIntervalSec <= 0 {
		return fmt.Errorf("AgentConfig: agentHeartbeatIntervalSec must be greater than 0")
	}
	if a.EnableLogs && a.LogsIntervalSec <= 0 {
		return fmt.Errorf("AgentConfig: logsIntervalSec must be greater than 0")
	}
	if !a.DisableFootprints && a.FootprintsReportingIntervalSec <= 0 {
		return fmt.Errorf("AgentConfig: footprintsIntervalSec must be greater than 0")
	}
	if !a.DisableFootprints && a.FootprintsCollectionIntervalSec <= 0 {
		return fmt.Errorf("AgentConfig: footprintsCollectionIntervalSec must be greater than 0")
	}

	if !a.DisableFootprints && a.FootprintsCollectionIntervalSec >= a.FootprintsReportingIntervalSec {
		return fmt.Errorf("AgentConfig: footprintsCollectionIntervalSec must be less than footprintsIntervalSec")
	}

	if !a.DisableTests && a.TestsIntervalSec <= 0 {
		return fmt.Errorf("AgentConfig: testsIntervalSec must be greater than 0")
	}
	if a.CheckExecutionIntervalSec <= 0 {
		return fmt.Errorf("AgentConfig: checkExecutionIntervalSec must be greater than 0")
	}
	if !a.EnableRemoteConfig && a.RemoteConfigIntervalSec <= 0 {
		return fmt.Errorf("AgentConfig: RemoteConfigIntervalSec must be greater than 0")
	}
	return nil
}

var _ models.Model = (*AgentConfig)(nil)
