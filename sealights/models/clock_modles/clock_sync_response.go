package clock_models

import (
	"fmt"

	"github.com/konflux-ci/release-service/__sealights__/models"
)

type ClockSyncResponse struct {
	YourTime    int64 `json:"your_time,omitempty"`
	CurrentTime int64 `json:"current_time,omitempty"`
	Offset      int64 `json:"offset,omitempty"`
}

func NewClockSyncResponse() *ClockSyncResponse {
	return &ClockSyncResponse{}
}

func (r *ClockSyncResponse) SetYourTime(yourTime int64) *ClockSyncResponse {
	r.YourTime = yourTime
	return r
}

func (r *ClockSyncResponse) SetCurrentTime(currentTime int64) *ClockSyncResponse {
	r.CurrentTime = currentTime
	return r
}

func (r *ClockSyncResponse) SetOffset(offset int64) *ClockSyncResponse {
	r.Offset = offset
	return r
}

func (r *ClockSyncResponse) Validate() error {
	if r.CurrentTime <= 0 {
		return fmt.Errorf("ClockSyncResponse: currentTime is required")
	}
	return nil
}

var _ models.Model = (*ClockSyncResponse)(nil)
