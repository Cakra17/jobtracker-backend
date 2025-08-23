package cvgenerator

type ResponseStat struct {
  TotalApplication  int `json:"total_application"`
  Pending           int `json:"pending"`
  Interview         int `json:"interview"`
  Offer             int `json:"offer"`
  Rejected          int `json:"rejected"`
  WithDraw          int `json:"withdraw"`
}
