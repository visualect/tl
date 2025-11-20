package dto

type AddTaskRequest struct {
	Task string `json:"task" validate:"required"`
}
