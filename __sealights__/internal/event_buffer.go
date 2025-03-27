package internal

import "sync"

type eventsBuffer struct {
	sync.Mutex
	pre    []Event
	buffer []Event
	size   int
}

func newEventBuffer(size int) *eventsBuffer {
	return &eventsBuffer{
		pre:    make([]Event, 0, size),
		buffer: make([]Event, 0, size),
		size:   size,
	}
}

func (b *eventsBuffer) addPre(events ...Event) {
	b.Lock()
	defer b.Unlock()
	b.pre = append(b.pre, events...)
}

func (b *eventsBuffer) addBuffer(events ...Event) {
	b.Lock()
	defer b.Unlock()
	b.buffer = append(b.buffer, events...)
}

func (b *eventsBuffer) get() []Event {
	b.Lock()
	defer b.Unlock()
	result := make([]Event, 0, len(b.pre)+len(b.buffer))
	result = append(result, b.pre...)
	result = append(result, b.buffer...)
	b.pre = make([]Event, 0, b.size)
	b.buffer = make([]Event, 0, b.size)
	return result
}
