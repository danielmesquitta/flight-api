package dto

type ErrorResponse struct {
	Message string      `json:"message,omitzero"`
	Errors  []ErrorItem `json:"errors,omitzero"`
}

type ErrorItem struct {
	Name   string `json:"name,omitzero"`
	Reason string `json:"reason,omitzero"`
}
