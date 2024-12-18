package handler

import (
	"blog-backend/api"
	"blog-backend/model"
	"database/sql"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

// Login
// (POST /auth/login)
func (h *Handler) PostAuthLogin(ctx echo.Context) error {
	var req api.PostAuthLoginJSONRequestBody
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, err)
	}
	// check Email and password
	user, err := h.Repo.GetUserByEmailAndPassword(ctx, string(req.Email), string(req.Password))
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, err)
	}
	if user.ID == uuid.Nil {
		return ctx.JSON(http.StatusUnauthorized, "invalid email or password")
	}
	return ctx.JSON(http.StatusOK, user)
}

// Logout
// (POST /auth/logout)
func (h *Handler) PostAuthLogout(ctx echo.Context) error {
	return ctx.JSON(http.StatusOK, "logout successful")
}

// Register a new user
// (POST /auth/register)
func (h *Handler) PostAuthRegister(ctx echo.Context) error {
	// create a new user
	var req api.PostAuthRegisterJSONRequestBody
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, err)
	}

	// hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, err)
	}

	err = h.Repo.CreateUser(ctx, model.User{
		ID:           uuid.New(),
		Email:        sql.NullString{String: string(req.Email), Valid: req.Email != ""},
		PasswordHash: sql.NullString{String: string(hashedPassword), Valid: true},
		CreatedAt:    time.Now(),
		IsAdmin:      false,
	})
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, err)
	}
	return ctx.JSON(http.StatusCreated, "User created")
}
