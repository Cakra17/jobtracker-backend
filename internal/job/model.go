package job

import (
	"database/sql"
	"encoding/json"
	"errors"


	"time"
)

type JobRequest struct {
	JobTitle       string      `json:"job_title" validation:"required"`
	JobUrl         NullString  `json:"job_url"`
	JobDescription NullString  `json:"job_description"`
	SalaryMin      NullFloat64 `json:"salary_min"`
	SalaryMax      NullFloat64 `json:"salary_max"`
	SalaryCurrency string      `json:"salary_currency"`
	Location       string      `json:"location"`
	EmploymentType string      `json:"employment_type"`
	WorkType       string      `json:"work_type"`
	Status         string      `json:"status"`
	Priority       string      `json:"priority"`
	AppliedDate    NullTime    `json:"applied_date"`
	Deadline       NullTime    `json:"deadline"`
	Notes          NullString  `json:"notes"`
	IsActive       bool        `json:"is_active"`
}

type GetJob struct {
	Limit  uint `query:"limit"`
	Offset uint `query:"offset"`

	UserId string
	Queries map[string]string
}

func (j *GetJob) Validate() error {
	queries := j.Queries

	if _,ok := queries["limit"]; !ok {
		return errors.New("limit is empty")
	}

	if _,ok := queries["offset"]; !ok {
		return errors.New("offset is empty")
	}

	return nil
}

type Job struct {
	ID             string      `db:"id"`
	User_ID        string      `db:"use_id"`
	JobTitle       string      `db:"job_title"`
	JobUrl         NullString  `db:"job_url"`
	JobDescription NullString  `db:"job_description"`
	SalaryMin      NullFloat64 `db:"salary_min"`
	SalaryMax      NullFloat64 `db:"salary_max"`
	SalaryCurrency string      `db:"salary_currency"`
	Location       string      `db:"location"`
	EmploymentType string      `db:"employment_type"`
	WorkType       string      `db:"work_type"`
	Status         string      `db:"status"`
	Priority       string      `db:"priority"`
	AppliedDate    NullTime    `db:"applied_date"`
	Deadline       NullTime    `db:"deadline"`
	Notes          NullString  `db:"notes"`
	IsActive       bool        `db:"is_active"`
	CreatedAt      time.Time   `db:"created_at"`
	UpdatedAt      time.Time   `db:"updated_at"`
}

// Null String custom Marshal/Unmarshaler
type NullString struct {
	sql.NullString
}

func (ns *NullString) MarshalJSON() ([]byte, error) {
	if ns.Valid {
		return json.Marshal(ns.String)
	}
	return json.Marshal(nil)
}

func (ns *NullString) UnmarshalJSON(data []byte) error {
	var s *string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	if s != nil {
		ns.String = *s
		ns.Valid = true
	} else {
		ns.Valid = false
	}

	return nil
}

// NullTime Custom Marshal/Unmarshaler
type NullTime struct {
	sql.NullTime
}

func (nt *NullTime) MarshalJSON() ([]byte, error) {
	if nt.Valid {
		return json.Marshal(nt.Time)
	}
	return json.Marshal(nil)
}

func (nt *NullTime) UnmarshalJSON(data []byte) error {
	var s *string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	if s == nil || *s == "" {
		nt.Valid = false
		return nil
	}

	t, err := time.Parse(time.RFC3339, *s)
	if err != nil {
		// If RFC3339 fails, try other common formats
		formats := []string{
			"2006-01-02T15:04:05Z07:00",
			"2006-01-02T15:04:05",
			"2006-01-02 15:04:05",
			"2006-01-02",
		}

		for _, format := range formats {
			if t, err = time.Parse(format, *s); err == nil {
				break
			}
		}

		if err != nil {
			return err
		}
	}

	nt.Time = t
	nt.Valid = true

	return nil
}

// NullFloat Custom Marshal/Unmarshaler
type NullFloat64 struct {
	sql.NullFloat64
}

func (nf *NullFloat64) MarshalJSON() ([]byte, error) {
	if nf.Valid {
		return json.Marshal(nf.Float64)
	}
	return json.Marshal(nil)
}

func (nf *NullFloat64) UnmarshalJSON(data []byte) error {
	var f *float64
	if err := json.Unmarshal(data, &f); err != nil {
		return err
	}

	if f != nil {
		nf.Float64 = *f
		nf.Valid = true
	} else {
		nf.Valid = false
	}

	return nil
}
