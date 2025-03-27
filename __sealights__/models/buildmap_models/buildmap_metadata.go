package buildmap_models

import (
	"fmt"

	"github.com/konflux-ci/release-service/__sealights__/models"
)

type BuildMapMetadata struct {
	CustomerId     string         `json:"customerId"`
	Technology     string         `json:"technology"`
	AppName        string         `json:"appName"`
	UniqueModuleID string         `json:"uniqueModuleId"`
	Build          string         `json:"build"`
	Branch         string         `json:"branch"`
	Generated      int64          `json:"generated"`
	AgentId        string         `json:"agentId"`
	BuildSessionId string         `json:"buildSessionId"`
	Commit         string         `json:"commit"`
	CommitsLog     []*CommitLog   `json:"commitsLog"`
	Contributors   []*Contributor `json:"contributors"`
	History        []string       `json:"history"`
	RepositoryUrl  string         `json:"repositoryUrl"`
	SCMProvider    string         `json:"scmProvider"`
	SCMVersion     string         `json:"scmVersion"`
	SCMBaseUrl     string         `json:"scmBaseUrl"`
	SCM            string         `json:"scm"`
}

func NewBuildMapMetadata() *BuildMapMetadata {
	return &BuildMapMetadata{
		CustomerId:     "",
		Technology:     "GoLang",
		AppName:        "",
		Build:          "",
		Branch:         "",
		Generated:      models.UnixMilli(),
		AgentId:        "",
		BuildSessionId: "",
		UniqueModuleID: "",
		Commit:         "",
		SCMProvider:    "",
		SCMVersion:     "",
		SCMBaseUrl:     "",
		SCM:            "",
	}
}

func (b *BuildMapMetadata) SetCustomerId(value string) *BuildMapMetadata {
	b.CustomerId = value
	return b
}

func (b *BuildMapMetadata) SetAppName(value string) *BuildMapMetadata {
	b.AppName = value
	return b
}

func (b *BuildMapMetadata) SetBuild(value string) *BuildMapMetadata {
	b.Build = value
	return b
}

func (b *BuildMapMetadata) SetBranch(value string) *BuildMapMetadata {
	b.Branch = value
	return b
}

func (b *BuildMapMetadata) SetAgentId(value string) *BuildMapMetadata {
	b.AgentId = value
	return b
}

func (b *BuildMapMetadata) SetGenerated(value int64) *BuildMapMetadata {
	b.Generated = value
	return b
}

func (b *BuildMapMetadata) SetBuildSessionId(value string) *BuildMapMetadata {
	b.BuildSessionId = value
	return b
}

func (b *BuildMapMetadata) SetUniqueModuleId(value string) *BuildMapMetadata {
	b.UniqueModuleID = value
	return b
}

func (b *BuildMapMetadata) SetCommit(value string) *BuildMapMetadata {
	b.Commit = value
	return b
}

func (b *BuildMapMetadata) SetCommitsLog(value []*CommitLog) *BuildMapMetadata {
	b.CommitsLog = value
	return b
}

func (b *BuildMapMetadata) SetContributors(value []*Contributor) *BuildMapMetadata {
	b.Contributors = value
	return b
}

func (b *BuildMapMetadata) SetHistory(value []string) *BuildMapMetadata {
	b.History = value
	return b
}
func (b *BuildMapMetadata) SetSCMProvider(value string) *BuildMapMetadata {
	b.SCMProvider = value
	return b
}

func (b *BuildMapMetadata) SetSCMVersion(value string) *BuildMapMetadata {
	b.SCMVersion = value
	return b
}

func (b *BuildMapMetadata) SetSCMBaseUrl(value string) *BuildMapMetadata {
	b.SCMBaseUrl = value
	return b
}

func (b *BuildMapMetadata) SetSCMType(value string) *BuildMapMetadata {
	b.SCM = value
	return b
}

func (b *BuildMapMetadata) SetRepositoryUrl(value string) *BuildMapMetadata {
	b.RepositoryUrl = value
	return b
}

func (b *BuildMapMetadata) Validate() error {
	if b.CustomerId == "" {
		return fmt.Errorf("BuildMapMetadata: customerId is required")
	}
	if b.AppName == "" {
		return fmt.Errorf("BuildMapMetadata: appName is required")
	}
	if b.Build == "" {
		return fmt.Errorf("BuildMapMetadata: build is required")
	}
	if b.Branch == "" {
		return fmt.Errorf("BuildMapMetadata: branch is required")
	}
	if b.AgentId == "" {
		return fmt.Errorf("BuildMapMetadata: agentId is required")
	}
	if b.BuildSessionId == "" {
		return fmt.Errorf("BuildMapMetadata: buildSessionId is required")
	}
	return nil
}

var _ models.Model = (*BuildMapMetadata)(nil)
