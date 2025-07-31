package user

import (
	"context"
	"database/sql"
)

type UserRepo struct {
	db *sql.DB
}

func NewUserRepo(db *sql.DB) UserRepo {
	return UserRepo{db: db}
}

func (r *UserRepo) CreateUser(ctx context.Context, user User) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	query := `
		INSERT INTO users
			(id, email, username, display_name, password_hash)
		VALUES
			(?, ?, ?, ?, ?)
	`

	_, err = tx.ExecContext(ctx, query, user.ID, user.Email, user.Username, user.DisplayName, user.PasswordHash)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (r *UserRepo) GetUserByEmail(ctx context.Context, email string) (User, error) {
	var result User
	query := `
		SELECT 
				id,
				email,
				username,
				display_name,
				password_hash,
				avatar_url,
				email_verified
		FROM
				users
		WHERE
				email = ?
		LIMIT 1
	`
	row:= r.db.QueryRowContext(ctx, query, email)
	err := row.Scan(
		&result.ID, 
		&result.Email, 
		&result.Username,
		&result.DisplayName, 
		&result.PasswordHash,
		&result.AvatarUrl, 
		&result.EmailVerified,
	)
	if err != nil {
		return result, err
	}

	return result, nil
}

func (r *UserRepo) UpdateName(ctx context.Context, user ChangeNameRequest) error {
	query := `
		UPDATE 
				users 
		SET
				username = ?,
				display_name = ?
		WHERE 
				id = ?
	`
	_, err := r.db.ExecContext(ctx, query, user.Username, user.DisplayName, user.ID)
	if err != nil {
		return err
	}

	return nil
}

func (r *UserRepo) DeleteUser(ctx context.Context, id string) error {
	query := `
		DELETE FROM 
				users
		WHERE 
				id = ?
	`
	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	return nil
}