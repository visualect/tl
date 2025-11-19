package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/visualect/tl/internal/auth"
	"github.com/visualect/tl/internal/dto"
	"github.com/visualect/tl/internal/repo"
	"golang.org/x/crypto/bcrypt"
)

func New(usersRepo repo.UsersRepository, tasksRepo repo.TasksRepository) *tasksHandler {
	return &tasksHandler{usersRepo, tasksRepo}
}

type tasksHandler struct {
	usersRepo repo.UsersRepository
	tasksRepo repo.TasksRepository
}

func (t *tasksHandler) SignUp(c echo.Context) error {
	data := new(dto.RegisterUserRequest)
	if err := c.Bind(data); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := c.Validate(data); err != nil {
		return err
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(data.Password), 10)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	err = t.usersRepo.Create(c.Request().Context(), dto.RegisterUserRequest{Login: data.Login, Password: string(hashed)})
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusCreated, data)
}

func (t *tasksHandler) Login(c echo.Context) error {
	data := new(dto.LoginUserRequest)
	if err := c.Bind(data); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := c.Validate(data); err != nil {
		return err
	}

	user, err := t.usersRepo.GetUserByLogin(c.Request().Context(), data.Login)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, err.Error())
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(data.Password))
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "Authenticaton failed. Password is incorrect")
	}

	token, err := auth.GenerateJWTToken(user.ID)
	if err != nil {
		echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	respBody := dto.LoginResponse{Token: token}
	return c.JSON(http.StatusOK, respBody)
}
