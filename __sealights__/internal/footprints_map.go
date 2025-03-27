package internal

import (
	"sort"
	"sync"
)

type footprintsMap struct {
	sync.Mutex
	m map[string]Event
}

func newFootprintsMap() *footprintsMap {
	return &footprintsMap{
		m: make(map[string]Event),
	}
}

func (b *footprintsMap) add(events ...Event) {
	b.Lock()
	defer b.Unlock()
	for _, event := range events {
		b.m[event.function] = event
	}
}

func (b *footprintsMap) get() []Event {
	b.Lock()
	result := make([]Event, 0, len(b.m))
	for _, event := range b.m {
		result = append(result, event)
	}
	b.m = make(map[string]Event)
	b.Unlock()
	sort.Slice(result, func(i, j int) bool {
		return result[i].timestamp < result[j].timestamp
	})
	return result
}

func (b *footprintsMap) len() int {
	b.Lock()
	defer b.Unlock()
	return len(b.m)
}
