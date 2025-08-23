package cvgenerator

import "time"

type Resume struct {
  PersonalInfo PersonalInfo `json:"personal_info" validate:"required"`
  Summary  string `json:"summary" validate:"required"`
  Education []Education `json:"education"`
  Experience []Experience `json:"experience"`
  Project []Project `json:"project"`
  Skills []Skills `json:"skills"`
}

type PersonalInfo struct {
	FullName string `json:"full_name" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
  Phone    string `json:"phone"`
	Location string `json:"location" validate:"required"`
	LinkedIn string `json:"linkedin,omitempty"`
	Website  string `json:"website,omitempty"`
	Github   string `json:"github,omitempty"`
}

type Education struct {
	Institution string `json:"institution" validate:"required"`
  StartDate   time.Time `json:"start_date" validate:"required"`
  EndDate *time.Time `json:"end_date"`
  Degree string `json:"degree" validate:"required"`
  Major string `json:"major"`
  GPA *float64 `json:"gpa"`
  IsCurrent bool `json:"is_current"`
  Order    int    `json:"order"`
}

type Experience struct {
  JobTitle    string    `json:"job_title" validate:"required"`
  Company     string    `json:"company" validate:"required"`
  Location    string    `json:"location"`
  StartDate   time.Time `json:"start_date" validate:"required"`
  EndDate     *time.Time `json:"end_date,omitempty"`
  IsCurrent   bool      `json:"is_current"`
  Description string    `json:"description"`
  Achievements []string `json:"achievements"`
  Technologies []string `json:"technologies,omitempty"`
  Order    int    `json:"order"`
}

type Project struct {
  Name string `json:"project_name" validate:"required"`
  Description []string `json:"description"`
  StartDate time.Time `json:"start_date"`
}

type Skills struct {
  Name     string `json:"name" validate:"required"`
  Category string `json:"category"`
  Level    string `json:"level"`
  Order    int    `json:"order"`
}
