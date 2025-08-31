package job

import (
	"bytes"
	"database/sql"
	"io"
	"strconv"
	"time"

	"github.com/360EntSecGroup-Skylar/excelize"
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
	jobGroup := r.Group("api/v1/jobs")

	jobGroup.Post("/", authMiddleware, h.AddJob)
	jobGroup.Get("/users", authMiddleware, h.GetAllJobByUserId)
	jobGroup.Get("/details/:id", authMiddleware, h.GetJobById)
	jobGroup.Patch("/state/:id", authMiddleware, h.ChangeState)
	jobGroup.Delete("/:id", authMiddleware, h.DeleteJob)
	jobGroup.Put("/:id", authMiddleware, h.UpdateJobById)

  jobGroup.Post("/bulk", authMiddleware, h.BulkCreateApplications)
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
		Position: payload.Position,
		Company: payload.Company,
		Platform: payload.Platform,
		Location: payload.Location,
		SalaryCurrency: payload.SalaryCurrency,
		EmploymentType: payload.EmploymentType,
		WorkType: payload.WorkType,
		Status: payload.Status,
		Priority: payload.Priority,
		AppliedDate: time.Time(payload.AppliedDate),
		Salary: payload.Salary,
		Notes: payload.Notes,
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
		Data: PartialResponse{
			Id: job.ID,
			Position: job.Position,
			Platform: job.Platform,
			Company: job.Company,
			Salary: job.Salary.Float64,
			SalaryCurrency: job.SalaryCurrency,
			Location: job.Location,
			EmploymentType: job.EmploymentType,
			WorkType: job.WorkType,
			Status: job.Status,
			Priority: job.Priority,
		},
	})
}

func (h *jobHandler) BulkCreateApplications(c *fiber.Ctx) error {
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
  
  userId := claim.UserID
  formFile, err := c.FormFile("data")
  if err != nil {
    return err
  }

  file, err := formFile.Open()
  if err != nil {
    return err
  }

  fileContent, err := io.ReadAll(file)
  if err != nil {
    return err
  }
  
  excel, err := excelize.OpenReader(bytes.NewReader(fileContent))
  if err != nil {
    return err
  }

  sheetName := "Sheet1"

  error := make([]BulkError, 0)
  rows := excel.GetRows(sheetName)
  success := 0

  for i, row := range rows {
    if  i == 0 {
      continue
    } 
    
    isSalaryValid := true
    salary, err := strconv.ParseFloat(row[3], 64)
    if err != nil {
      isSalaryValid = false
    }
    
    AppliedDate := ConvertFromRawExcelToDate(row[10])

    cell := Job{
      ID: uuid.NewString(),
      User_ID: userId,
      Position: row[0],
      Platform: row[1],
      Company: row[2],
      Salary: NullFloat64{ sql.NullFloat64{Float64: salary, Valid: isSalaryValid} },
      SalaryCurrency: row[4],
      Location: row[5],
      EmploymentType: row[6],
      WorkType: row[7],
      Status: row[8],
      Priority: row[9],
      AppliedDate: AppliedDate,
      Notes: NullString{ sql.NullString{String: row[11], Valid: true} },
    }

    err = h.jobRepo.AddJob(c.Context(), cell)
    if err != nil {
      e := BulkError {
        Index: i,
        Error: err.Error(),
      }
      error = append(error, e)
    } else {
      success++
    }
  }

  response := BulkApplicationResponse {
    TotalRequested: success + len(error),
    Successful: success,
    Failed: len(error),
    Errors: error,
  }

  if response.Failed == 0 {
    return c.Status(fiber.StatusOK).JSON(model.DataResponse{
      Status: "success",
      Data: response,
    })
  } else if response.Successful == 0 {
    return c.Status(fiber.StatusBadRequest).JSON(model.DataResponse{
      Status: "fail",
      Data: response,
    })
  } else {
    return c.Status(fiber.StatusMultiStatus).JSON(model.DataResponse{
      Status: "multi-status",
      Data: response,
    })  
  }
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
  var stat Stat

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
			Id: job.ID,
			Position: job.Position,
			Platform: job.Platform,
			Company: job.Company,
			SalaryCurrency: job.SalaryCurrency,
			Salary: job.Salary.Float64,
			Location: job.Location,
			EmploymentType: job.EmploymentType,
			WorkType: job.WorkType,
			Status: job.Status,
			Priority: job.Priority,
		})
	}

  stat, err = h.jobRepo.GetStat(c.Context(), claim.UserID) 
  if err != nil {
    return err
  }

  statRes := &model.ResponseStat{
    TotalApplication: stat.TotalApplication,
    Pending: stat.Pending,
    Interview: stat.Interview,
    Rejected: stat.Rejected,
    WithDraw: stat.WithDraw,
    Offer: stat.Offer,
  }

	return c.Status(fiber.StatusOK).JSON(model.DataResponse{
		Status: "success",
		Data: response, 
    Stat: statRes,
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
		Position: job.Position,
		Company: job.Company,
		Platform: job.Platform,
		Salary: job.Salary.Float64,
		SalaryCurrency: job.SalaryCurrency,
		Location: job.Location,
		EmploymentType: job.EmploymentType,
		WorkType: job.WorkType,
		Status: job.Status,
		Priority: job.Priority,
		AppliedDate: job.AppliedDate,
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
		Position: payload.Position,
		Company: payload.Company,
		Platform: payload.Platform,
		Salary: payload.Salary,
		SalaryCurrency: payload.SalaryCurrency,
		Location: payload.Location,
		EmploymentType: payload.EmploymentType,
		WorkType: payload.WorkType,
		Status: payload.Status,
		Priority: payload.Priority,
		AppliedDate: time.Time(payload.AppliedDate),
		Notes: payload.Notes,
	}

	err = h.jobRepo.UpdateJob(c.Context(), job)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(model.DataResponse{
		Status: "success",
		Message: "job updated successfully",
		Data: PartialResponse{
			Id: job.ID,
			Position: job.Position,
			Platform: job.Platform,
			Company: job.Company,
			Salary: job.Salary.Float64,
			SalaryCurrency: job.SalaryCurrency,
			Location: job.Location,
			EmploymentType: job.EmploymentType,
			WorkType: job.WorkType,
			Status: job.Status,
			Priority: job.Priority,
		},
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
