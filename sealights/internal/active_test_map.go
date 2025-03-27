package internal

import (
	"sync"
)

type activeTestMap struct {
	sync.RWMutex
	activeTests map[string]bool
}

func newActiveTestMap() *activeTestMap {
	return &activeTestMap{
		activeTests: make(map[string]bool),
	}
}

func (m *activeTestMap) add(testName string) {
	m.Lock()
	defer m.Unlock()
	m.activeTests[testName] = true
}

func (m *activeTestMap) remove(testName string) {
	m.Lock()
	defer m.Unlock()
	delete(m.activeTests, testName)
}

func (m *activeTestMap) list() []string {
	m.RLock()
	defer m.RUnlock()
	if len(m.activeTests) == 0 {
		return nil
	}
	result := make([]string, 0, len(m.activeTests))
	for testName := range m.activeTests {
		result = append(result, testName)
	}
	return result
}
