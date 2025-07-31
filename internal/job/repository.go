package job

import (
	"context"
	"database/sql"
	"fmt"
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

	query := `
		INSERT INTO 
			job_applications
			(id, user_id, job_title, location, salary_currency, employment_type, work_type, status, priority)
		VALUES
			(?, ?, ?, ?, ?, ?, ?, ?, ?)
	`
	_,err = tx.ExecContext(
		ctx, 
		query, 
		job.ID,
		job.User_ID,
		job.JobTitle,
		job.Location,
		job.SalaryCurrency,
		job.EmploymentType,
		job.WorkType,
		job.Status,
		job.Priority,
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
				job_title,
				location,
				salary_currency,
				employment_type,
				work_type,
				status,
				priority
		FROM
				job_applications
		WHERE
				user_id = ?
	`
	args := []interface{}{payload.UserId}

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
			&job.JobTitle,
			&job.Location,
			&job.SalaryCurrency,
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

func getLimitAndOffset(req GetJob) (string, []interface{}) {
	query := `LIMIT ? OFFSET ?`

	limit := req.Limit
	if limit <= 0 {
		limit = 10
	}

	offset := req.Offset

	args := []interface{}{limit, offset}

	return query, args
}

func (r *JobRepo) GetJobById(ctx context.Context, id string) (Job, error) {
	var job Job

	query := `
		SELECT
				id,
				user_id,
				job_title,
				job_url,
				job_description,
				location,
				salary_min,
				salary_max,
				salary_currency,
				employment_type,
				work_type,
				status,
				priority,
				applied_date,
				deadline,
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
		&job.JobTitle,
		&job.JobUrl,
		&job.JobDescription,
		&job.Location,
		&job.SalaryMin,
		&job.SalaryMax,
		&job.SalaryCurrency,
		&job.EmploymentType,
		&job.WorkType,
		&job.Status,
		&job.Priority,
		&job.AppliedDate,
		&job.Deadline,
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
				job_title = ?,
				job_url = ?,
				job_description = ?,
				location = ?,
				salary_min = ?,
				salary_max = ?,
				salary_currency = ?,
				employment_type = ?,
				work_type = ?,
				status = ?,
				priority = ?,
				applied_date = ?,
				deadline = ?,
				notes = ?
		WHERE
				id = ?
	`
	_, err := r.db.ExecContext(
		ctx, query, 
		job.JobTitle,
		job.JobUrl,
		job.JobDescription,
		job.Location,
		job.SalaryMin,
		job.SalaryMax,
		job.SalaryCurrency,
		job.EmploymentType,
		job.WorkType,
		job.Status,
		job.Priority,
		job.AppliedDate,
		job.Deadline,
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