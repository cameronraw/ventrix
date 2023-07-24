package main

import "github.com/cameronraw/ventrix/cmd/queue"


type ListenEventRequest struct {
	ServiceName string `json:"service_name"`
	Type        string `json:"type"`
}

type RegisterEventRequest struct {
	EventType string `json:"event_type"`
}

type RegisterServiceRequest struct {
	Name     string `json:"name"`
	Endpoint string `json:"endpoint"`
}

type QueueEventRequest struct {
	Type    string            `json:"type"`
	Payload map[string]string `json:"payload"`
	Timeout int               `json:"timeout"`
}

func (request *QueueEventRequest) ToEvent() queue.Event {
	return queue.Event{
		Type:    request.Type,
		Payload: request.Payload,
		Timeout: request.Timeout,
	}
}

