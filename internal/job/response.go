package job

import (
	"time"
)

type PartialResponse struct {
	JobTitle       string `json:"job_title"`
	SalaryCurrency string `json:"salary_currency"`
	Location       string `json:"location"`
	EmploymentType string `json:"employment_type"`
	WorkType       string `json:"work_type"`
	Status         string `json:"status"`
	Priority       string `json:"priority"`
}

type FullResponse struct {
	JobTitle       string    `json:"job_title" validation:"required"`
	JobUrl         string    `json:"job_url"`
	JobDescription string    `json:"job_description"`
	SalaryMin      float64   `json:"salary_min"`
	SalaryMax      float64   `json:"salary_max"`
	SalaryCurrency string    `json:"salary_currency"`
	Location       string    `json:"location"`
	EmploymentType string    `json:"employment_type"`
	WorkType       string    `json:"work_type"`
	Status         string    `json:"status"`
	Priority       string    `json:"priority"`
	AppliedDate    NullTime   `json:"applied_date"`
	Deadline       NullTime    `json:"deadline"`
	Notes          string    `json:"notes"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}