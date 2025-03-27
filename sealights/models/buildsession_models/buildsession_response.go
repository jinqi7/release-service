package buildsession_models

import (
	"fmt"

	"github.com/konflux-ci/release-service/__sealights__/models"
)

type BuildSessionResponse struct {
	BuildSessionId string `json:"buildSessionId"`
}

func NewGetBuildSessionIdResponse() *BuildSessionResponse {
	return &BuildSessionResponse{BuildSessionId: ""}
}

func (b *BuildSessionResponse) SetBuildSessionId(buildSessionId string) *BuildSessionResponse {
	b.BuildSessionId = buildSessionId
	return b
}

func (b *BuildSessionResponse) Validate() error {
	if b.BuildSessionId == "" {
		return fmt.Errorf("BuildSessionResponse: BuildSessionId is required")
	}
	return nil
}

var _ models.Model = (*BuildSessionResponse)(nil)
