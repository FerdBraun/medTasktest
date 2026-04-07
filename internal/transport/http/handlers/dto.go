package handlers

import (
	"time"

	taskdomain "example.com/taskservice/internal/domain/task"
)

type taskMutationDTO struct {
	Title            string            `json:"title"`
	Description      string            `json:"description"`
	Status           taskdomain.Status `json:"status"`
	PeriodicityType  string            `json:"periodicity_type"`
	PeriodicityValue string            `json:"periodicity_value"`
}

type taskDTO struct {
	ID               int64             `json:"id"`
	Title            string            `json:"title"`
	Description      string            `json:"description"`
	Status           taskdomain.Status `json:"status"`
	PeriodicityType  string            `json:"periodicity_type"`
	PeriodicityValue string            `json:"periodicity_value"`
	LastCompletedAt  *time.Time        `json:"last_completed_at,omitempty"`
	CreatedAt        time.Time         `json:"created_at"`
	UpdatedAt        time.Time         `json:"updated_at"`
}

func newTaskDTO(task *taskdomain.Task) taskDTO {
	return taskDTO{
		ID:               task.ID,
		Title:            task.Title,
		Description:      task.Description,
		Status:           task.Status,
		PeriodicityType:  task.PeriodicityType,
		PeriodicityValue: task.PeriodicityValue,
		LastCompletedAt:  task.LastCompletedAt,
		CreatedAt:        task.CreatedAt,
		UpdatedAt:        task.UpdatedAt,
	}
}
