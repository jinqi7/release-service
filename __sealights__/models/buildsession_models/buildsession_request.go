package buildsession_models

import (
	"fmt"

	"github.com/konflux-ci/release-service/__sealights__/models"
)

type BuildSessionRequest struct {
	AppName        string `json:"appName"`
	BuildName      string `json:"buildName"`
	BranchName     string `json:"branchName"`
	BuildSessionId string `json:"buildSessionId"`
}

func NewBuildSessionRequest() *BuildSessionRequest {
	return &BuildSessionRequest{}
}

func (b *BuildSessionRequest) Validate() error {
	if b.AppName == "" {
		return fmt.Errorf("BuildSessionRequest: appName is required")
	}
	if b.BuildName == "" {
		return fmt.Errorf("BuildSessionRequest: buildName is required")
	}
	if b.BranchName == "" {
		return fmt.Errorf("BuildSessionRequest: branchName is required")
	}
	return nil
}

func (b *BuildSessionRequest) SetAppName(appName string) *BuildSessionRequest {
	b.AppName = appName
	return b
}

func (b *BuildSessionRequest) SetBuildName(buildName string) *BuildSessionRequest {
	b.BuildName = buildName
	return b
}

func (b *BuildSessionRequest) SetBranchName(branchName string) *BuildSessionRequest {
	b.BranchName = branchName
	return b
}

func (b *BuildSessionRequest) SetBuildSessionId(buildSessionId string) *BuildSessionRequest {
	b.BuildSessionId = buildSessionId
	return b
}

var _ models.Model = (*BuildSessionRequest)(nil)
