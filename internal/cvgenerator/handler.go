package cvgenerator

import (
	"fmt"
	"strings"

	"github.com/Cakra17/JobTracker-Api/internal/model"
	"github.com/Cakra17/JobTracker-Api/internal/ratelimiter"
	"github.com/Cakra17/JobTracker-Api/pkg/jwt"
	"github.com/Cakra17/JobTracker-Api/pkg/validation"
	"github.com/gofiber/fiber/v2"
)

type resumeHandler struct {
	pdfGenerator *PDFGenerator
  jwtProvider *jwt.JWTProvider
	rateLimiter *ratelimiter.RateLimiter
}

type ResumeHandlerConfig struct {
  PdfGenerator *PDFGenerator
  JWTProvider *jwt.JWTProvider
	RateLimiter *ratelimiter.RateLimiter
}

func NewCVGenerator(rhc ResumeHandlerConfig) resumeHandler {
	return resumeHandler{
    pdfGenerator: rhc.PdfGenerator,
    jwtProvider: rhc.JWTProvider,
    rateLimiter: rhc.RateLimiter,
  }
}

func (h *resumeHandler) RegisterRoute(r *fiber.App) {
  authMiddleware := h.jwtProvider.Middleware()
  resumeGroup := r.Group("api/v1/resumes")

  resumeGroup.Post("/generate/pdf", authMiddleware, h.GeneratePDF)
}

func (h *resumeHandler) GeneratePDF(c *fiber.Ctx) error {
  if err := h.rateLimiter.Middleware(c.IP()); err != nil {
    return err
  }
  
  _, err := jwt.GetLoggedInUser(c)
  if err != nil {
    return c.Status(fiber.StatusForbidden).JSON(model.ErrorResponse{
			Status: "fail",
			Message: "Forbidden access",
		})
  }

  var resume Resume

  if err := c.BodyParser(&resume); err != nil {
    return c.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse{
			Status: "fail",
			Message: "Request malformated",
		})
  }

  if err := validation.Validate(&resume); err != nil {
    return c.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse{
			Status: "fail",
			Message: err.Error(),
		})
  }

  pdf, err := h.pdfGenerator.GeneratePDF(c.Context(), &resume)
  if err != nil {
    return err
  }

  filename := fmt.Sprintf("%s_resume.pdf", strings.ReplaceAll(resume.PersonalInfo.FullName, " ", "_"))
  c.Set("Content-Type", "application/pdf")
  c.Set("Content-Disposition", fmt.Sprintf("attachment; filename:\"%s\"", filename))
  return c.Status(fiber.StatusCreated).Send(pdf)
}
