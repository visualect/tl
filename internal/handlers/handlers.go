package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/visualect/tl/internal/auth"
	"github.com/visualect/tl/internal/dto"
	"github.com/visualect/tl/internal/repo"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
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
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return echo.NewHTTPError(http.StatusBadRequest, "user with that login name already exists")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusCreated, dto.RegisterResponse{Login: data.Login})
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
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return echo.NewHTTPError(http.StatusUnauthorized, "user does not exists")
		}
		return echo.NewHTTPError(http.StatusUnauthorized, err.Error())
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(data.Password))
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "authenticaton failed. password is incorrect")
	}

	token, err := auth.GenerateJWTToken(user.ID, user.Login)
	if err != nil {
		echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	respBody := dto.LoginResponse{Token: token}
	return c.JSON(http.StatusOK, respBody)
}

func (t *tasksHandler) GetUser(c echo.Context) error {
	pc := auth.GetPrivateClaims(c)
	u, err := t.usersRepo.GetUserByLogin(c.Request().Context(), pc.Login)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, err.Error())
	}
	user := dto.UserResponse{Login: u.Login}
	return c.JSON(http.StatusOK, user)
}

func (t *tasksHandler) AddTask(c echo.Context) error {
	data := new(dto.AddTaskRequest)
	if err := c.Bind(data); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := c.Validate(data); err != nil {
		return err
	}

	pc := auth.GetPrivateClaims(c)

	err := t.tasksRepo.CreateTask(c.Request().Context(), pc.UserID, data.Task)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusCreated, data)
}

func (t *tasksHandler) GetTasks(c echo.Context) error {
	pc := auth.GetPrivateClaims(c)

	tasks, err := t.tasksRepo.GetTasksByUserID(c.Request().Context(), pc.UserID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return echo.NewHTTPError(http.StatusBadRequest, "user not found")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, tasks)
}

func (t *tasksHandler) ToggleCompleteTask(c echo.Context) error {
	taskID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	pc := auth.GetPrivateClaims(c)

	err = t.tasksRepo.ToggleCompleteTaskByID(c.Request().Context(), taskID, pc.UserID)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return c.NoContent(http.StatusOK)
}

func (t *tasksHandler) DeleteTask(c echo.Context) error {
	taskID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	pc := auth.GetPrivateClaims(c)

	err = t.tasksRepo.DeleteTaskByID(c.Request().Context(), taskID, pc.UserID)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return c.NoContent(http.StatusOK)
}
