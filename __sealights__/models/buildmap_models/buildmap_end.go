package buildmap_models

import (
	"fmt"

	"github.com/konflux-ci/release-service/__sealights__/models"
)

type BuildMapEndStatus struct {
	Success  bool   `json:"success"`
	Message  string `json:"message"`
	Duration int64  `json:"duration"`
}

func NewBuildMapEndStatus() *BuildMapEndStatus {
	return &BuildMapEndStatus{}
}

type BuildMapEnd struct {
	AppName string             `json:"appName"`
	Build   string             `json:"build"`
	Branch  string             `json:"branch"`
	Status  *BuildMapEndStatus `json:"status"`
}

func NewBuildMapEnd() *BuildMapEnd {
	return &BuildMapEnd{
		AppName: "",
		Build:   "",
		Branch:  "",
		Status:  NewBuildMapEndStatus(),
	}
}

func (b *BuildMapEnd) SetAppName(appName string) *BuildMapEnd {
	b.AppName = appName
	return b
}

func (b *BuildMapEnd) SetBuild(build string) *BuildMapEnd {
	b.Build = build
	return b
}

func (b *BuildMapEnd) SetBranch(branch string) *BuildMapEnd {
	b.Branch = branch
	return b
}

func (b *BuildMapEnd) SetStatusSuccess(success bool) *BuildMapEnd {
	b.Status.Success = success
	return b
}

func (b *BuildMapEnd) SetStatusMessage(message string) *BuildMapEnd {
	b.Status.Message = message
	return b
}

func (b *BuildMapEnd) SetStatusDuration(duration int64) *BuildMapEnd {
	b.Status.Duration = duration
	return b
}

func (b *BuildMapEnd) Validate() error {
	if b.AppName == "" {
		return fmt.Errorf("BuildMapEnd: appName is required")
	}
	if b.Build == "" {
		return fmt.Errorf("BuildMapEnd: build is required")
	}
	if b.Branch == "" {
		return fmt.Errorf("BuildMapEnd: branch is required")
	}
	if b.Status.Duration == 0 {
		return fmt.Errorf("BuildMapEnd: status.duration is required")
	}
	return nil
}

var _ models.Model = (*BuildMapEnd)(nil)
