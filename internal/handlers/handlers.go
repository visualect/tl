package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/visualect/tl/internal/dto"
	"github.com/visualect/tl/internal/repo"
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

	err := t.usersRepo.Create(c.Request().Context(), *data)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusCreated, data)
}
