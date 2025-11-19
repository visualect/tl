package repo

import (
	"gorm.io/gorm"
)

type TasksRepository interface {
	// AddTask(ctx context.Context, task string) error
	// GetTasks(ctx context.Context) ([]models.Task, error)
	// CompleteTask(ctx context.Context, id int) error
	// DeleteTask(ctx context.Context, id int) error
}

type tasksRepo struct {
	db *gorm.DB
}

func NewTasks(db *gorm.DB) TasksRepository {
	return &tasksRepo{db}
}

// func (t *tasksRepo) AddTask(ctx context.Context, task string) error
// func (t *tasksRepo) GetTasks(ctx context.Context) ([]models.Task, error)
// func (t *tasksRepo) CompleteTask(ctx context.Context, id int) error
// func (t *tasksRepo) DeleteTask(ctx context.Context, id int) error
