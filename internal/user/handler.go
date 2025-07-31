package user

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/Cakra17/JobTracker-Api/internal/model"
	"github.com/Cakra17/JobTracker-Api/internal/ratelimiter"
	"github.com/Cakra17/JobTracker-Api/pkg/jwt"
	"github.com/Cakra17/JobTracker-Api/pkg/validation"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type userHandler struct {
	userRepo *UserRepo
	jwtProvider *jwt.JWTProvider
	saltCost int
	ratelimiter *ratelimiter.RateLimiter
}

type UserHandlerConfig struct {
	UserRepo *UserRepo
	JWTProvider *jwt.JWTProvider
	SaltCost int
	RateLimiter *ratelimiter.RateLimiter
}

func NewUserHandler(cfg UserHandlerConfig) userHandler {
	return userHandler{
		userRepo: cfg.UserRepo,
		jwtProvider: cfg.JWTProvider,
		saltCost: cfg.SaltCost,
		ratelimiter: cfg.RateLimiter,
	}
}

func (h *userHandler) RegisterRoute(r *fiber.App) {
	authMiddleware := h.jwtProvider.Middleware()
	userGroup := r.Group("v1/user")

	userGroup.Post("/register", h.RegisterUser)
	userGroup.Post("/login", h.LoginUser)
	userGroup.Get("/auth", authMiddleware, h.VerifyAuth)
	userGroup.Patch("/name", authMiddleware, h.ChangeName)
}

func (h *userHandler) RegisterUser(c *fiber.Ctx) error {
	if err := h.ratelimiter.Middleware(c.IP()); err != nil {
		return err
	}

	var payload RegisterUserRequest
	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse{
			Status: "fail",
			Message: "Request Malformed",
		})
	}

	if err := validation.Validate(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse{
			Status: "fail",
			Message: err.Error(),
		})
	}

	if payload.Password != payload.ConfirmPassword {
		return c.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse{
			Status: "fail",
			Message: "Make sure the password and confirm password are match",
		})
	}

	user, accessToken, err := h.createUser(c.Context(), payload)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse{
			Status: "fail",
			Message: err.Error(),
		})
	}

	cookie := fiber.Cookie{
		Name: "access_token",
		Value: accessToken,
		Expires: time.Now().Add(time.Hour * 8),
		HTTPOnly: false,
		Path: "/",
		SameSite: "Lax",
		Secure: false,
	}

	c.Cookie(&cookie)

	return c.Status(fiber.StatusCreated).JSON(model.DataResponse{
		Status: "success",
		Message: "User registered successfully",
		Data: UserResponse{
			Email: user.Email,
			Username: user.Username,
			AccessToken: accessToken,
		},
	})
}

func (h *userHandler) createUser(ctx context.Context, payload RegisterUserRequest) (User, string, error) {
	_, err := h.userRepo.GetUserByEmail(ctx, payload.Email)
	if err != nil && err != sql.ErrNoRows {
		return User{}, "", errors.New(
			fmt.Sprintf("%s: %s", err, "GetUserByEmail error"),
		)
	}

	if err == nil {
		return User{}, "", fiber.NewError(fiber.StatusConflict, "User already exist")
	}
	
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(payload.Password), h.saltCost)
	if err != nil {
		return User{}, "", err
	}

	user := User {
		ID: uuid.NewString(),
		Email: payload.Email,
		Username: payload.Username,
		PasswordHash: string(hashedPassword),
		DisplayName: payload.DisplayName,
	}

	err = h.userRepo.CreateUser(ctx, user)
	if err != nil {
		return user, "", err
	}

	accessToken, err := h.generateAccessTokenFromUser(user)
	if err != nil {
		return user, "", err
	}

	return user, accessToken, nil
}

func (h *userHandler) LoginUser(c *fiber.Ctx) error {
	if err := h.ratelimiter.Middleware(c.IP()); err != nil {
		return err
	}

	var payload LoginUserRequest

	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse{
			Status: "fail",
			Message: "Couldn't parse the information, please send again",
		})
	}

	if err :=  validation.Validate(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse{
			Status: "fail",
			Message: "Validation failed, please send the correct login information",
		})
	}

	user, accessToken, err := h.authenticate(c.Context(), payload)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse{
			Status: "fail",
			Message: err.Error(),
		})
	}

	cookie := fiber.Cookie{
		Name: "access_token",
		Value: accessToken,
		Expires: time.Now().Add(time.Hour * 8),
		HTTPOnly: false,
		Path: "/",
		SameSite: "Lax",
		Secure: false,
	}

	c.Cookie(&cookie)

	return c.Status(fiber.StatusOK).JSON(model.DataResponse{
		Status: "success",
		Message: "user logged successfully",
		Data: UserResponse{
			Email: user.Email,
			Username: user.Username,
			AccessToken: accessToken,
		},
	})
}

func (h *userHandler) authenticate(ctx context.Context, payload LoginUserRequest) (User, string, error) {
	user , err := h.userRepo.GetUserByEmail(ctx, payload.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			return user, "", fiber.NewError(fiber.StatusNotFound, "user with the spesific credential not found")
		}

		return user, "", errors.New(
			fmt.Sprintf("%s : %s", err, "GetUserByEmail error"),
		)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(payload.Password)); err != nil {
		return user, "", fiber.NewError(fiber.StatusBadRequest, "wrong password")
	}

	accessToken, err := h.generateAccessTokenFromUser(user)
	if err != nil {
		return user, "", err
	}

	return user, accessToken, nil
}

func (h *userHandler) ChangeName(c *fiber.Ctx) error {
	if err := h.ratelimiter.Middleware(c.IP()); err != nil {
		return err
	}

	var payload ChangeNameRequest

	claim, err := jwt.GetLoggedInUser(c)
	if err != nil {
		return c.Status(fiber.StatusForbidden).JSON(model.ErrorResponse{
			Status: "fail",
			Message: "request forbidden",
		})
	}

	payload.ID = claim.UserID

	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse{
			Status: "fail",
			Message: "Couldn't parse the information, please send again",
		})
	}

	if err := validation.Validate(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse{
			Status: "fail",
			Message: err.Error(),
		})
	}

	err = h.updateName(c.Context(), payload)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(model.ErrorResponse{
			Status: "fail",
			Message: err.Error(),
		})
	}
	return c.Status(fiber.StatusAccepted).JSON(model.DataResponse{
		Status: "success",
		Message: "Name updated successfully",
	})
}

func (h *userHandler) updateName(ctx context.Context, payload ChangeNameRequest) error {
	err := h.userRepo.UpdateName(ctx, payload)
	if err != nil {
		return err
	}
	return nil
}

func (h *userHandler) VerifyAuth(c *fiber.Ctx) error {
	if err := h.ratelimiter.Middleware(c.IP()); err != nil {
		return err
	}

	_, err := jwt.GetLoggedInUser(c)
	if err != nil {
		return c.Status(fiber.StatusForbidden).JSON(model.ErrorResponse{
			Status: "fail",
			Message: "request forbidden",
		})
	}

	return c.Status(fiber.StatusOK).JSON(model.DataResponse{
		Status: "success",
		Message: "User validated",
	})
}

func (h *userHandler) generateAccessTokenFromUser(user User) (string, error) {
	claims := jwt.BuildJWTClaims(jwt.JWTUser{
		UserID: user.ID,
		Username: user.Username,
		Email: user.Email,
	}, time.Hour * 8)

	accessToken, err := h.jwtProvider.GenerateToken(claims)
	if err != nil {
		return "", err
	}

	return accessToken, nil
}