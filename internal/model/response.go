package model

type ResponseMeta struct {
	Limit  uint `json:"limit"`
	Offset uint `json:"offset"`
	Total  uint `json:"page"`
}

type DataResponse struct {
	Status  string        `json:"status"`
	Message string        `json:"message"`
	Data    any           `json:"data,omitempty"`
	Meta    *ResponseMeta `json:"meta,omitempty"`
  Stat    *ResponseStat `json:"stat,omitempty"`
}

type ResponseStat struct {
  TotalApplication  int `json:"total_application"`
  Pending           int `json:"pending"`
  Interview         int `json:"interview"`
  Offer             int `json:"offer"`
  Rejected          int `json:"rejected"`
  WithDraw          int `json:"withdraw"`
}

type ErrorResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}
