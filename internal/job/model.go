package job

import (
	"database/sql"
	"encoding/json"
	"errors"
	"strings"

	"time"
)

type JobRequest struct {
	Postion        string      `json:"position" validate:"required"`
	Platform       string      `json:"platform" validate:"required"`
	Company        string      `json:"company" validate:"required"`
	Salary         NullFloat64 `json:"salary"`
	SalaryCurrency string      `json:"salary_currency" validate:"required"`
	Location       string      `json:"location" validate:"required"`
	EmploymentType string      `json:"employment_type" validate:"required"`
	WorkType       string      `json:"work_type" validate:"required"`
  Status         string      `json:"status" validate:"required"`
	Priority       string      `json:"priority" validate:"required"`
	AppliedDate    Date   		 `json:"applied_date" validate:"required"`
	Notes          NullString  `json:"notes"`
	IsActive       bool        `json:"is_active"`
}

type GetJob struct {
	Limit  uint `query:"limit"`
	Offset uint `query:"offset"`

	UserId  string
	Queries map[string]string
}

func (j *GetJob) Validate() error {
	queries := j.Queries

	if _, ok := queries["limit"]; !ok {
		return errors.New("limit is empty")
	}

	if _, ok := queries["offset"]; !ok {
		return errors.New("offset is empty")
	}

	return nil
}

type Job struct {
	ID             string      `db:"id"`
	User_ID        string      `db:"use_id"`
	Position       string      `db:"position"`
	Company        string      `db:"company"`
	Platform       string      `db:"platform"`
	Salary         NullFloat64 `db:"salary"`
	SalaryCurrency string      `db:"salary_currency"`
	Location       string      `db:"location"`
	EmploymentType string      `db:"employment_type"`
	WorkType       string      `db:"work_type"`
	Status         string      `db:"status"`
	Priority       string      `db:"priority"`
	AppliedDate    Date   		 `db:"applied_date"`
	Notes          NullString  `db:"notes"`
	IsActive       bool        `db:"is_active"`
	CreatedAt      time.Time   `db:"created_at"`
	UpdatedAt      time.Time   `db:"updated_at"`
}

type Stat struct {
  TotalApplication  int `db:"total_application"`
  Pending           int `db:"pending"`
  Interview         int `db:"interview"`
  Offer             int `db:"offer"`
  Rejected          int `db:"rejected"`
  WithDraw          int `db:"withdraw"`
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

type Date time.Time

func (d *Date) UnmarshalJSON(data []byte) error {
	str := strings.Trim(string(data), "\"")
	if str == "null" || str == "" {
		return errors.New("date is malformed")
	}

	layout := "2006-01-02"
	t, err := time.Parse(layout, str)
	if err != nil {
		return err
	}

	*d = Date(t)
	return nil
}

func (d Date) MarshalJSON() ([]byte, error) {
	return json.Marshal(time.Time(d).Format("2006-01-02"))
}
