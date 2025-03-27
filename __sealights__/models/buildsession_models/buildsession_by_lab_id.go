package buildsession_models

import (
	"fmt"

	"github.com/konflux-ci/release-service/__sealights__/models"
)

type BuildSessionByLabIdRequest struct {
	LabId string `json:"labId"`
}

func NewBuildSessionByLabIdRequest() *BuildSessionByLabIdRequest {
	return &BuildSessionByLabIdRequest{}
}

func (b *BuildSessionByLabIdRequest) SetLabId(labId string) *BuildSessionByLabIdRequest {
	b.LabId = labId
	return b
}
func (b *BuildSessionByLabIdRequest) Validate() error {
	if b.LabId == "" {
		return fmt.Errorf("BuildSessionByLabIdRequest: labId is not valid")
	}
	return nil
}

var _ models.Model = (*BuildSessionRequest)(nil)
