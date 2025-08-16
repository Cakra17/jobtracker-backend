package job

import (
	"time"
)

type PartialResponse struct {
	Id             string `json:"id"`
	Position       string `json:"position"`
	Platform       string `json:"platform"`
	Company        string `json:"company"`
	Salary         float64 `json:"salary"`
	SalaryCurrency string `json:"salary_currency"`
	Location       string `json:"location"`
	EmploymentType string `json:"employment_type"`
	WorkType       string `json:"work_type"`
	Status         string `json:"status"`
	Priority       string `json:"priority"`
}

type FullResponse struct {
	Id             string    `json:"id"`
	Position       string    `json:"position"`
	Company        string    `json:"company"`
	Platform       string    `json:"platform"`
	Salary         float64   `json:"salary"`
	SalaryCurrency string    `json:"salary_currency"`
	Location       string    `json:"location"`
	EmploymentType string    `json:"employment_type"`
	WorkType       string    `json:"work_type"`
	Status         string    `json:"status"`
	Priority       string    `json:"priority"`
	AppliedDate    time.Time `json:"applied_date"`
	Notes          string    `json:"notes"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type BulkApplicationResponse struct {
  TotalRequested int          `json:"total_requested"`
  Successful     int          `json:"successful"`
  Failed         int          `json:"failed"`
  Errors         []BulkError  `json:"errors"`
}

type BulkError struct {
  Index int    `json:"index"`
  Error string `json:"error"`
}
