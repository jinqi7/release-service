package logs_models

import (
	"fmt"

	"github.com/konflux-ci/release-service/__sealights__/models"
)

type LogsSubmission struct {
	AppName     string `json:"appName"`
	BranchName  string `json:"branchName"`
	BuildName   string `json:"buildName"`
	Environment struct {
		EnvironmentName string `json:"environmentName"`
		MachineName     string `json:"machineName"`
	} `json:"environment"`
	Log []string `json:"logs"`
}

func NewLogsSubmission() *LogsSubmission {
	return &LogsSubmission{}
}

func (s *LogsSubmission) SetAppName(appName string) *LogsSubmission {
	s.AppName = appName
	return s
}

func (s *LogsSubmission) SetBranchName(branchName string) *LogsSubmission {
	s.BranchName = branchName
	return s
}

func (s *LogsSubmission) SetBuildName(buildName string) *LogsSubmission {
	s.BuildName = buildName
	return s
}

func (s *LogsSubmission) SetEnvironmentName(environmentName string) *LogsSubmission {
	s.Environment.EnvironmentName = environmentName
	return s
}

func (s *LogsSubmission) SetMachineName(machineName string) *LogsSubmission {
	s.Environment.MachineName = machineName
	return s
}

func (s *LogsSubmission) SetLog(log []string) *LogsSubmission {
	s.Log = log
	return s
}

func (s *LogsSubmission) Validate() error {
	if s.AppName == "" {
		return fmt.Errorf("LogsSubmission: appName is required")
	}
	if s.BranchName == "" {
		return fmt.Errorf("LogsSubmission: branchName is required")
	}
	if s.BuildName == "" {
		return fmt.Errorf("LogsSubmission: buildName is required")
	}
	if len(s.Log) == 0 {
		return fmt.Errorf("LogsSubmission: log is required")
	}
	return nil
}

var _ models.Model = (*LogsSubmission)(nil)
