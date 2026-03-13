package message

import "time"


type MsgType string


const (
	ERROR MsgType = "error"
	LOG   MsgType = "log"
)


type Event struct {
	Type      MsgType   `json:"type"`
	Text      string    `json:"text"`
    Timestamp time.Time `json:"timestamp"`
}


func NewEvent(msgType MsgType, message string) Event {
	if msgType != ERROR && msgType != LOG {
		msgType = LOG // default to progress if invalid
	}
	return Event{
		Type:    msgType,
		Text:    message,
        Timestamp: time.Now(),
	}
}


func Error(message string) Event {
	return NewEvent(ERROR, message)
}

func Log(message string) Event {
	return NewEvent(LOG, message)
}
