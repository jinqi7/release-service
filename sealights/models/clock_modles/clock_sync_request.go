package clock_models

import (
	"fmt"
	"time"

	"github.com/konflux-ci/release-service/__sealights__/models"
)

type ClockSyncRequest struct {
	MyTime int64 `json:"myTime"`
}

func NewClockSyncRequest() *ClockSyncRequest {
	return &ClockSyncRequest{
		time.Now().UTC().UnixNano() / int64(time.Millisecond),
	}
}

func (r *ClockSyncRequest) SetMyTime(myTime int64) *ClockSyncRequest {
	r.MyTime = myTime
	return r
}

func (r *ClockSyncRequest) Validate() error {
	if r.MyTime <= 0 {
		return fmt.Errorf("ClockSyncRequest: myTime is required")
	}
	return nil
}

var _ models.Model = (*ClockSyncRequest)(nil)
