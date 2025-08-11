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
	AppliedDate    Date 		 `json:"applied_date"`
	Notes          string    `json:"notes"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}
