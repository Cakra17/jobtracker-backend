package job

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"
)

type JobRepo struct {
	db *sql.DB
}

func NewJobRepo(db *sql.DB) JobRepo {
	return JobRepo{db: db}
}

func (r *JobRepo) AddJob(ctx context.Context, job Job) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	query, value := r.insertBuilder(job)

	_,err = tx.ExecContext(
		ctx, 
		query, 
		value...
	)

	if err != nil {
		return err
	}

	return tx.Commit()
}

func (r *JobRepo) GetJobByUserId(ctx context.Context, payload GetJob) ([]Job, error) {
	var jobs []Job

	query := `
		SELECT
				id,
				user_id,
				position,
				company,
				platform,
				location,
				salary_currency,
				salary,
				employment_type,
				work_type,
				status,
				priority
		FROM
				job_applications
		WHERE
				user_id = ?
		AND
				is_active = 1
		ORDER BY created_at DESC
	`
	args := []any{payload.UserId}

	limitQuery, limitArgs := getLimitAndOffset(payload)
	args = append(args, limitArgs...)

	completeQuery := fmt.Sprintf("%s %s", query, limitQuery)

	rows, err := r.db.QueryContext(ctx, completeQuery, args...)
	if err != nil {
		return jobs, nil
	}	

	for rows.Next() {
		var job Job
		if err := rows.Scan(
			&job.ID,
			&job.User_ID,
			&job.Position,
			&job.Company,
			&job.Platform,
			&job.Location,
			&job.SalaryCurrency,
			&job.Salary,
			&job.EmploymentType,
			&job.WorkType,
			&job.Status,
			&job.Priority,
		); err != nil {
			return jobs, err
		}

		jobs = append(jobs, job)
	}
	
	return jobs, nil
}

func getLimitAndOffset(req GetJob) (string, []any) {
	query := `LIMIT ? OFFSET ?`

	limit := req.Limit
	if limit <= 0 {
		limit = 10
	}

	offset := req.Offset

	args := []any{limit, offset}

	return query, args
}

func(r *JobRepo) GetStat(ctx context.Context, id string) (Stat, error) {
  var stat Stat
  query := `
    SELECT 
        COUNT(*) AS total_application,
        SUM(status = "pending") AS pending,
        SUM(status = "interview") AS interview,
        SUM(status = "offer") AS offer,
        SUM(status = "rejected") AS rejected,
        SUM(status = "withdraw") AS withdraw
    FROM
        job_applications
    WHERE
        user_id = ?
  `
  row := r.db.QueryRowContext(ctx, query, id)
  if err := row.Scan(
    &stat.TotalApplication,
    &stat.Pending,
    &stat.Interview,
    &stat.Offer,
    &stat.Rejected,
    &stat.WithDraw,
  ); err != nil {
    return stat, err
  }

  return stat, nil
}

func (r *JobRepo) GetJobById(ctx context.Context, id string) (Job, error) {
	var job Job

	query := `
		SELECT
				id,
				user_id,
				position,
				company,
				platform,
				location,
				salary,
				salary_currency,
				employment_type,
				work_type,
				status,
				priority,
				applied_date,
				notes,
				is_active,
				created_at,
				updated_at
		FROM
				job_applications
		WHERE
				id = ?
		LIMIT 1
	`
	row := r.db.QueryRowContext(ctx, query, id)
	if err := row.Scan(
		&job.ID,
		&job.User_ID,
		&job.Position,
		&job.Company,
		&job.Position,
		&job.Location,
		&job.Salary,
		&job.SalaryCurrency,
		&job.EmploymentType,
		&job.WorkType,
		&job.Status,
		&job.Priority,
		&job.AppliedDate,
		&job.Notes,
		&job.IsActive,
		&job.CreatedAt,
		&job.UpdatedAt,
	); err != nil {
		return job, err
	}

	return job, nil
}

func (r *JobRepo) UpdateJob(ctx context.Context, job Job) error {
	query := `
		UPDATE 
			job_applications SET
				position = ?,
				company = ?,
				platform = ?,
				location = ?,
				salary = ?,
				salary_currency = ?,
				employment_type = ?,
				work_type = ?,
				status = ?,
				priority = ?,
				applied_date = ?,
				notes = ?
		WHERE
				id = ?
	`
	appliedDate := time.Time(job.AppliedDate).Format("2006-01-02")
	_, err := r.db.ExecContext(
		ctx, query, 
		job.Position,
		job.Company,
		job.Platform,
		job.Location,
		job.Salary,
		job.SalaryCurrency,
		job.EmploymentType,
		job.WorkType,
		job.Status,
		job.Priority,
		appliedDate,
		job.Notes,
		job.ID,
	)
	if err != nil {
		return err
	}
	
	return nil
}

func (r *JobRepo) ChangeState(ctx context.Context, id string, value bool) error {
	query := `
		UPDATE
			job_applications SET
				is_active = ?
		WHERE
			id = ?	
	`
	_, err := r.db.ExecContext(ctx, query, value, id)
	if err != nil {
		return err
	}
	return nil
}

func (r *JobRepo) HardDelete(ctx context.Context, id string) error {
	query := `
		DELETE 
			FROM 
				job_applications
		WHERE
			id = ?	
	`

	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}
	return nil
}

func (r *JobRepo) insertBuilder(job Job) (string, []any) {
	baseColumn := []string{
		"id", "user_id", "position", "company", "platform",
		"location", "employment_type", "work_type", "status",
		"priority", "applied_date", "salary_currency",
	}

	appliedDate := time.Time(job.AppliedDate).Format("2006-01-02")

	baseValue := []any{
		job.ID, job.User_ID, job.Position, job.Company, job.Platform,
		job.Location, job.EmploymentType, job.WorkType, job.Status,
		job.Priority, appliedDate, job.SalaryCurrency,
	}

	if job.Salary.Valid {
		baseColumn = append(baseColumn, "salary")
		baseValue = append(baseValue, job.Salary.Float64)
	}

	if job.Notes.Valid {
		baseColumn = append(baseColumn, "notes")
		baseValue = append(baseValue, job.Notes.String)
	}

	columns := strings.Join(baseColumn, ",")
	placeholder := strings.Repeat("?, ", len(baseValue)-1) + "?"

	query := fmt.Sprintf(`
		INSERT INTO job_applications (%s)
		VALUES (%s)`, columns, placeholder)
	
	return query, baseValue
}
