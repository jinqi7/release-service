package footprint_models

import (
	"fmt"

	"github.com/konflux-ci/release-service/__sealights__/models"
)

type FootprintErrorEvent struct {
	Type           int    `json:"type"`
	UtcTimestampMs int64  `json:"utcTimestamp_ms"`
	Data           string `json:"data"`
}

func NewFootprintErrorEvent() *FootprintErrorEvent {
	return &FootprintErrorEvent{
		Type:           4006,
		UtcTimestampMs: models.UnixMilli(),
	}
}

func (f *FootprintErrorEvent) Validate() error {
	if f.Type != 4006 {
		return fmt.Errorf("FootprintErrorEvent: type is not valid")
	}
	if f.UtcTimestampMs == 0 {
		return fmt.Errorf("FootprintErrorEvent: utcTimestamp_ms is not valid")
	}
	if f.Data == "" {
		return fmt.Errorf("FootprintErrorEvent: data is not valid")
	}
	return nil
}

func (f *FootprintErrorEvent) SetData(value string) *FootprintErrorEvent {
	f.Data = value
	return f
}

var _ models.Model = (*FootprintErrorEvent)(nil)
