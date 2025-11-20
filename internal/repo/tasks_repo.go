package repo

import (
	"context"
	"errors"

	"github.com/visualect/tl/internal/models"
	"gorm.io/gorm"
)

type TasksRepository interface {
	CreateTask(ctx context.Context, userID int, task string) error
	GetTasksByUserID(ctx context.Context, userID int) ([]models.Task, error)
	CompleteTaskByID(ctx context.Context, taskID int, userID int) error
	DeleteTaskByID(ctx context.Context, taskID int, userID int) error
}

type tasksRepo struct {
	db *gorm.DB
}

func NewTasks(db *gorm.DB) TasksRepository {
	return &tasksRepo{db}
}

func (t *tasksRepo) CreateTask(ctx context.Context, userID int, task string) error {
	err := gorm.G[models.Task](t.db).Create(ctx, &models.Task{UserID: userID, Task: task})
	return err
}

func (t *tasksRepo) GetTasksByUserID(ctx context.Context, userID int) ([]models.Task, error) {
	tasks, err := gorm.G[models.Task](t.db).Where("user_id = ?", userID).Find(ctx)
	return tasks, err
}

func (t *tasksRepo) CompleteTaskByID(ctx context.Context, taskID int, userID int) error {
	rows, err := gorm.G[models.Task](t.db).Where("id = ? AND user_id = ?", taskID, userID).Update(ctx, "completed", true)
	if rows == 0 {
		return errors.New("unable to delete task")
	}
	return err
}

func (t *tasksRepo) DeleteTaskByID(ctx context.Context, taskID int, userID int) error {
	rows, err := gorm.G[models.Task](t.db).Where("id = ? AND user_id = ?", taskID, userID).Delete(ctx)
	if rows == 0 {
		return errors.New("unable to delete task")
	}
	return err
}
