package models

type Error struct {
	// message
	Message string `json:"message,omitempty"`

	// request id
	RequestID string `json:"request_id,omitempty"`
}
