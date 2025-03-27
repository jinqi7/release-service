package clock

import "sync"

var GlobalClockSync *ClockSync = NewClockSync()

type ClockSync struct {
	sync.RWMutex
	CurrentOffset int64
}

func NewClockSync() *ClockSync {
	return &ClockSync{}
}

func (c *ClockSync) SetCurrentOffset(offset int64) {
	c.Lock()
	defer c.Unlock()
	c.CurrentOffset = offset
}

func (c *ClockSync) GetCurrentOffset() int64 {
	c.RLock()
	defer c.RUnlock()
	return c.CurrentOffset
}
