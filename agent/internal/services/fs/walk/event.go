package walk

import (
	"vcx/pkg/logging"
	"vcx/pkg/toolkit/debugkit"
)


type EventType string

const (
	ERROR    EventType = "error"
	FILE     EventType = "file"
	DIR      EventType = "dir"
	SKIP     EventType = "skip"
	SYM      EventType = "sym"
)

var EVENT_TYPES = map[EventType]struct{}{
    ERROR:    {},
    FILE:     {},
    DIR:      {},
    SKIP:     {},
    SYM:      {},
}


type Event struct {
	Type EventType `json:"type"`
	Data string    `json:"data"`
}



func NewEvent(eventType EventType, message string) Event {
    if _, ok := EVENT_TYPES[eventType]; !ok {
        eventType = ERROR // default to ERROR if invalid
    }
	return newEvent(eventType, message)
}


func newEvent(msgType EventType, message string) Event {
	return Event{
		Type:    msgType,
		Data:    message,
	}
}


func Error(message string) Event {
	return newEvent(ERROR, message)
}


func File(message string) Event {
	return newEvent(FILE, message)
}


func Dir(message string) Event {
	return newEvent(DIR, message)
}


func Skip(message string) Event {
	return newEvent(SKIP, message)
}


func Sym(message string) Event {
	return newEvent(SYM, message)
}


func (e Event) Log() {
    // TEMP for debugging
    log := logging.GetLogger()
    log.Debug("Walk Event", "type", e.Type, "data", e.Data, "caller", debugkit.GetCaller(1))
}
