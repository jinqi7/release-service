package __sealights__

import (
	"fmt"
	"sync/atomic"
)

type collectorStats struct {
	mainNotReady uint32
	mainReady    uint32
	mainNotInMap uint32
}

func newCollectorStats() *collectorStats {
	return &collectorStats{
		mainNotReady: 0,
		mainReady:    0,
		mainNotInMap: 0,
	}
}

func (s *collectorStats) setMainNotReady() {
	atomic.AddUint32(&s.mainNotReady, 1)
}

func (s *collectorStats) setMainReady() {
	atomic.AddUint32(&s.mainReady, 1)
}

func (s *collectorStats) setMainNotInMap() {
	atomic.AddUint32(&s.mainNotInMap, 1)
}

func (s *collectorStats) String() string {
	return fmt.Sprintf("mainNotReady: %d, mainReady: %d, mainNotInMap: %d", atomic.LoadUint32(&s.mainNotReady), atomic.LoadUint32(&s.mainReady), atomic.LoadUint32(&s.mainNotInMap))
}
