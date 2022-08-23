package common

type Destination struct {
	To        string   `json:"to" validate:"required,uuid4"`
	Endpoints []string `json:"endpoints" validate:"required,dive,uri"`
}
type EventRequest struct {
	Event     string                 `json:"event" validate:"required"`
	Route     string                 `json:"route" validate:"required"`
	From      string                 `json:"from" validate:"required,uuid4"`
	To        []Destination          `json:"to" validate:"required,dive"`
	Reference string                 `json:"reference" validate:"required"`
	CreatedAt string                 `json:"created_at" validate:"required,datetime=2006-01-02T15:04:05Z07:00"`
	Payload   map[string]interface{} `json:"payload"`
}

type Event struct {
	EventRequest

	RequestId string `json:"request_id"`
}
