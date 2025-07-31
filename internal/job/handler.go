package job

import (
	"fmt"
	"strconv"

	"github.com/Cakra17/JobTracker-Api/internal/model"
	"github.com/Cakra17/JobTracker-Api/internal/ratelimiter"
	"github.com/Cakra17/JobTracker-Api/pkg/jwt"
	"github.com/Cakra17/JobTracker-Api/pkg/validation"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type jobHandler struct {
	jobRepo *JobRepo
	jwtProvider *jwt.JWTProvider
	ratelimiter *ratelimiter.RateLimiter
}

type JobHandlerConfig struct {
	JobRepo *JobRepo
	JWTProvider *jwt.JWTProvider
	RateLimiter *ratelimiter.RateLimiter
}

func NewJobHandler(jhCfg JobHandlerConfig) jobHandler {
	return jobHandler{
		jobRepo: jhCfg.JobRepo,
		jwtProvider: jhCfg.JWTProvider,
		ratelimiter: jhCfg.RateLimiter,
	}
}

func (h *jobHandler) RegisterRoute(r *fiber.App) {
	authMiddleware := h.jwtProvider.Middleware()
	jobGroup := r.Group("v1/jobs")

	jobGroup.Post("/", authMiddleware, h.AddJob)
	jobGroup.Get("/users", authMiddleware, h.GetAllJobByUserId)
	jobGroup.Get("/details/:id", authMiddleware, h.GetJobById)
	jobGroup.Patch("/state/:id", authMiddleware, h.ChangeState)
	jobGroup.Delete("/:id", authMiddleware, h.DeleteJob)
	jobGroup.Put("/:id", authMiddleware, h.UpdateJobById)
}

func (h *jobHandler) AddJob(c *fiber.Ctx) error {
	if err := h.ratelimiter.Middleware(c.IP()); err != nil {
		return err
	}

	var payload JobRequest

	claim, err := jwt.GetLoggedInUser(c)
	if err != nil {
		return c.Status(fiber.StatusForbidden).JSON(model.ErrorResponse{
			Status: "fail",
			Message: "Forbidden access",
		})
	}

	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse{
			Status: "fail",
			Message: "Request malformated",
		})
	}

	if err := validation.Validate(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse{
			Status: "fail",
			Message: err.Error(),
		})
	}

	job := Job{
		ID: uuid.NewString(),
		User_ID: claim.UserID,
		JobTitle: payload.JobTitle,
		Location: payload.Location,
		SalaryCurrency: payload.SalaryCurrency,
		EmploymentType: payload.EmploymentType,
		WorkType: payload.WorkType,
		Status: payload.Status,
		Priority: payload.Priority,
	}

	ctx := c.Context()

	err = h.jobRepo.AddJob(ctx, job)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(model.ErrorResponse{
			Status: "fail",
			Message: err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(model.DataResponse {
		Status: "success",
		Message: "Job added successfully",
	})
}


func (h *jobHandler) GetAllJobByUserId(c *fiber.Ctx) error {
	if err := h.ratelimiter.Middleware(c.IP()); err != nil {
		return err
	}

	claim, err := jwt.GetLoggedInUser(c)
	if err != nil {
		return c.Status(fiber.StatusForbidden).JSON(model.ErrorResponse{
			Status: "fail",
			Message: "Forbidden request",
		})
	}

	var jobQueries GetJob

	jobQueries.UserId = claim.UserID
	jobQueries.Queries = c.Queries()

	if err := jobQueries.Validate(); err != nil {
		return err
	}

	limit, _ := strconv.ParseUint(jobQueries.Queries["limit"], 10, 32)
	offset, _ := strconv.ParseUint(jobQueries.Queries["offset"], 10, 32)

	jobQueries.Limit = uint(limit)
	jobQueries.Offset = uint(offset)

	jobs, err := h.jobRepo.GetJobByUserId(c.Context(), jobQueries)
	if err != nil {
		return err
	}

	response := []PartialResponse{}
	for _, job := range jobs {
		response = append(response, PartialResponse{
			JobTitle: job.JobTitle,
			SalaryCurrency: job.SalaryCurrency,
			Location: job.Location,
			EmploymentType: job.EmploymentType,
			WorkType: job.WorkType,
			Status: job.Status,
			Priority: job.Priority,
		})
	}

	return c.Status(fiber.StatusOK).JSON(model.DataResponse{
		Status: "success",
		Data: response,
	})
}

func (h *jobHandler) GetJobById(c *fiber.Ctx) error {
	if err := h.ratelimiter.Middleware(c.IP()); err != nil {
		return err
	}

	_, err := jwt.GetLoggedInUser(c)
	if err != nil {
		return c.Status(fiber.StatusForbidden).JSON(model.ErrorResponse{
			Status: "fail",
			Message: "Forbidden request",
		})
	}

	id := c.Params("id", "")
	
	job, err := h.jobRepo.GetJobById(c.Context(), id)
	if err != nil {
		return err
	}

	response := FullResponse{
		JobTitle: job.JobTitle,
		JobUrl: job.JobUrl.String,
		JobDescription: job.JobDescription.String,
		SalaryMin: job.SalaryMin.Float64,
		SalaryMax: job.SalaryMax.Float64,
		SalaryCurrency: job.SalaryCurrency,
		Location: job.Location,
		EmploymentType: job.EmploymentType,
		WorkType: job.WorkType,
		Status: job.Status,
		Priority: job.Priority,
		AppliedDate: job.AppliedDate,
		Deadline: job.Deadline,
		Notes: job.Notes.String,
		CreatedAt: job.CreatedAt,
		UpdatedAt: job.UpdatedAt,
	}
	return c.Status(fiber.StatusOK).JSON(model.DataResponse{
		Status: "success",
		Data: response,
	})
}

func (h *jobHandler) UpdateJobById(c *fiber.Ctx) error {
	if err := h.ratelimiter.Middleware(c.IP()); err != nil {
		return err
	}

	var payload JobRequest
	_, err := jwt.GetLoggedInUser(c)
	if err != nil {
		return c.Status(fiber.StatusForbidden).JSON(model.ErrorResponse{
			Status: "fail",
			Message: "Forbidden request",
		})
	}

	id := c.Params("id", "")

	if err := c.BodyParser(&payload); err != nil {
		fmt.Println(err)
		return c.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse{
			Status: "fail",
			Message: "Request malformated",
		})
	}

	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse{
			Status: "fail",
			Message: err.Error(),
		})
	}

	job := Job{
		ID: id,
		JobTitle: payload.JobTitle,
		JobUrl: payload.JobUrl,
		JobDescription: payload.JobDescription,
		SalaryMin: payload.SalaryMin,
		SalaryMax: payload.SalaryMax,
		SalaryCurrency: payload.SalaryCurrency,
		Location: payload.Location,
		EmploymentType: payload.EmploymentType,
		WorkType: payload.WorkType,
		Status: payload.Status,
		Priority: payload.Priority,
		AppliedDate: payload.AppliedDate,
		Deadline: payload.Deadline,
		Notes: payload.Notes,
	}

	err = h.jobRepo.UpdateJob(c.Context(), job)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(model.DataResponse{
		Status: "success",
		Message: "job updated successfully",
	})
}

func (h *jobHandler) ChangeState(c *fiber.Ctx) error {
	if err := h.ratelimiter.Middleware(c.IP()); err != nil {
		return err
	}

	_, err := jwt.GetLoggedInUser(c)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse{
			Status: "fail",
			Message: "Request malformated",
		})
	}

	id := c.Params("id", "")
	job, err := h.jobRepo.GetJobById(c.Context(), id)
	state := !job.IsActive

	err = h.jobRepo.ChangeState(c.Context(), id, state)
	if err != nil {
		return err
	}

	if !state {
		return c.Status(fiber.StatusOK).JSON(model.DataResponse{
			Status: "success",
			Message: "removed successfully",
		})
	}

	return c.Status(fiber.StatusOK).JSON(model.DataResponse{
		Status: "success",
		Message: "restored successfully",
	})
} 

func (h *jobHandler) DeleteJob(c *fiber.Ctx) error {
	if err := h.ratelimiter.Middleware(c.IP()); err != nil {
		return err
	}

	_, err := jwt.GetLoggedInUser(c)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse{
			Status: "fail",
			Message: "Request malformated",
		})
	}

	id := c.Params("id", "")

	err = h.jobRepo.HardDelete(c.Context(), id)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(model.DataResponse{
		Status: "success",
		Message: "deleted successfully",
	})
}