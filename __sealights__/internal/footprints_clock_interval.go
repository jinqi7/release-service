package internal

import (
	"context"
	"sync"
	"time"

	"github.com/konflux-ci/release-service/__sealights__/pkg/clock"
)

type footprintsCollectionClockInterval struct {
	sync.RWMutex
	startTimestamp int64
	endTimestamp   int64
	interval       time.Duration
	timer          *time.Timer
}

func newFootprintsClockInterval(interval time.Duration) *footprintsCollectionClockInterval {
	currentTimeSec := clock.UnixSeconds()
	c := &footprintsCollectionClockInterval{
		interval:       interval,
		startTimestamp: currentTimeSec,
		endTimestamp:   currentTimeSec + int64(interval.Seconds()),
		timer:          time.NewTimer(interval),
	}
	return c
}

func (c *footprintsCollectionClockInterval) setTimeInterval(interval time.Duration) *footprintsCollectionClockInterval {
	c.interval = interval
	if c.timer != nil {
		c.timer.Reset(interval)
	} else {
		c.timer = time.NewTimer(interval)
	}
	c.setTimestamps()
	return c
}

func (c *footprintsCollectionClockInterval) setTimestamps() {
	c.Lock()
	defer c.Unlock()
	c.startTimestamp = clock.UnixSeconds()
	c.endTimestamp = c.startTimestamp + int64(c.interval.Seconds())
}

func (c *footprintsCollectionClockInterval) run(ctx context.Context) {
	go func() {
		for {
			select {
			case <-ctx.Done():
				c.timer.Stop()
				return
			case <-c.timer.C:
				c.timer.Reset(c.interval)
				c.setTimestamps()
			}
		}
	}()
}

func (c *footprintsCollectionClockInterval) getsTimestamps() (int64, int64) {
	c.RLock()
	defer c.RUnlock()
	return c.startTimestamp, c.endTimestamp
}
