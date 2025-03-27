package internal

type EventType int

const (
	EventTypeTrace         EventType = 1
	EventTypeStartMain     EventType = 2
	EventTypeEndMain       EventType = 3
	EventTypeStartTest     EventType = 4
	EventTypeEndTest       EventType = 5
	EventTypeStartTestMain EventType = 6
	EventTypeEndTestMain   EventType = 7
)

type Event struct {
	eventType           EventType
	timestamp           int64
	function            string
	attributes          map[string]string
	activeTestList      []string
	footprintsStartTime int64
	footprintsEndTime   int64
}

func NewRawEvent() Event {
	return Event{
		eventType:           0,
		timestamp:           0,
		function:            "",
		attributes:          nil,
		activeTestList:      nil,
		footprintsStartTime: 0,
		footprintsEndTime:   0,
	}
}

func (e Event) SetEventType(eventType EventType) Event {
	e.eventType = eventType
	return e
}

func (e Event) SetTimestamp(timestamp int64) Event {
	e.timestamp = timestamp
	return e
}

func (e Event) SetFunction(function string) Event {
	e.function = function
	return e
}

func (e Event) SetAttribute(key string, value string) Event {
	if e.attributes == nil {
		e.attributes = make(map[string]string)
	}
	e.attributes[key] = value
	return e
}

func (e Event) SetFootprintsTimestamps(footprintsStartTime, footprintsEndTime int64) Event {
	e.footprintsStartTime = footprintsStartTime
	e.footprintsEndTime = footprintsEndTime
	return e
}

func (e Event) SetActiveTestList(activeTestList []string) Event {
	e.activeTestList = activeTestList
	return e
}
