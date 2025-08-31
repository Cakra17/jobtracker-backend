# JobTracker Backend

A minimal REST API built with Go for managing job applications, user authentication, session management, and CV generation. The backend supports user registration, job application tracking, rate limiting, JWT-based authentication, and ATS-friendly CV generation using templates.

## Features

- User management (register, login, session handling)
- Job application tracking (create, update, delete applications)
- CV generation with ATS-friendly templates
- Rate limiting for API protection
- JWT-based authentication
- PostgreSQL database with migrations

## Prerequisites

- Go (version 1.21 or higher)
- MySql (version 9 or higher)
- Migrate (https://github.com/golang-migrate/migrate) for database migrations
- Git

## Installation

1. Clone the repository:
```bash
   git clone git@github.com:Cakra17/jobtracker-backend.git
   cd jobtracker-backend
```
2. Install dependencies:
``` bash
   go mod download
```
3. Set up environment variables by creating a `.env` file in the project root. Use the following template:
```bash
# Database
DB_PASSWORD=
DB_USERNAME=
DB_PORT=
DB_HOST=
DB_NAME=
DB_MAX_OPEN_CONNECTION=
DB_MAX_IDLE_CONNECTION=
DB_MAX_CONN_LIFETIME=
DB_MAX_CONN_IDLE_TIME=

# App
APP_PORT=
ENV=

# JWT
JWT_SECRET=
BCRYPT_SALT=
```
4. Apply database migrations:
```bash
    migrate -database mysql://user:password@tcp(host:port)/dbname?query -path db/migration up
```

## Running the App

1. Build and run the application:
```bash
   go run cmd/jobtracker/main.go
```
   The server will start on `http://localhost:8080` (or the port specified in `.env`).

2. Test the API using tools like Postman or `curl`.

## API Endpoints

Below are the primary endpoints (assumed based on the codebase structure). All endpoints require JSON payloads and return JSON responses.

| Method | Endpoint | Description | Authentication |
| --- | --- | --- | --- |
| POST | `/api/v1/users/register` | Register a new user | None |
| POST | `/api/v1/users/login` | Login and receive JWT token | None |
| GET | `/api/v1/users/me` | Get current user details | JWT |
| POST | `/api/v1/jobs` | Create a new job application | JWT |
| POST | `/api/v1/jobs/bulk` | Create bulk new job application using xlsx file | JWT |
| GET | `/api/v1/jobs/users` | List all job applications by users | JWT |
| GET | `/api/v1/jobs/details/:id` | Get a specific job application | JWT |
| PUT | `/api/v1/jobs/:id` | Update a job application | JWT |
| DELETE | `/api/v1/jobs/:id` | Delete a job application | JWT |
| POST | `/api/v1/resumes/generate/pdf` | Generate an ATS-friendly CV | JWT |

**Example Request (Register User)**:
```bash
curl -X POST http://localhost:8080/api/v1/users/register \\\
-H "Content-Type: application/json" \\\
-d '{"email":"user@example.com","password":"securepassword","username":"John Doe", "confirm_password": "securepassword", "display_name": "john doe",}'
```
## Project Structure

- `cmd/jobtracker/`: Application entry point
- `internal/`: Core app logic (user, job, cvgenerator, ratelimiter)
-`pkg/`: Reusable utilities (JWT, validation)
- `db/migration/`: Database migration files
- `templates/`: HTML templates for CV generation