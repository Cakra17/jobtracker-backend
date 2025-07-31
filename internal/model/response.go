package model

type ResponseMeta struct {
	Limit uint `json:"limit"`
	Offset uint `json:"offset"`
	Total uint `json:"page"`
}

type DataResponse struct {
	Status string `json:"status"`
	Message string `json:"message"`
	Data any `json:"data,omitempty"`
	Meta *ResponseMeta `json:"meta,omitempty"`
}

type ErrorResponse struct {
	Status string `json:"status"`
	Message string `json:"message"`
}