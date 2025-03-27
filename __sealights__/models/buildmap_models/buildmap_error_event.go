package buildmap_models

import (
	"fmt"

	"github.com/konflux-ci/release-service/__sealights__/models"
)

type BuildMapErrorEvent struct {
	Type           int    `json:"type"`
	UtcTimestampMs int64  `json:"utcTimestamp_ms"`
	Data           string `json:"data"`
}

func NewBuildMapErrorEvent() *BuildMapErrorEvent {
	return &BuildMapErrorEvent{
		Type:           4005,
		UtcTimestampMs: models.UnixMilli(),
	}
}

func (b *BuildMapErrorEvent) Validate() error {
	if b.Type != 4005 {
		return fmt.Errorf("BuildMapErrorEvent: type is not valid")
	}
	if b.UtcTimestampMs == 0 {
		return fmt.Errorf("BuildMapErrorEvent: utcTimestamp_ms is not valid")
	}
	if b.Data == "" {
		return fmt.Errorf("BuildMapErrorEvent: data is not valid")
	}
	return nil
}

func (b *BuildMapErrorEvent) SetData(value string) *BuildMapErrorEvent {
	b.Data = value
	return b
}

var _ models.Model = (*BuildMapErrorEvent)(nil)
