package task

import "time"

type Status string

const (
	StatusNew        Status = "new"
	StatusInProgress Status = "in_progress"
	StatusDone       Status = "done"
)

type Task struct {
	ID          int64     `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Status      Status    `json:"status"`
	PeriodicityType  string     `json:"periodicity_type"`  // "no", "day", "month", "date", "odd"
	PeriodicityValue string     `json:"periodicity_value"` // "3", "15", "2026-04-10,2026-04-12", "even"
	LastCompletedAt  *time.Time `json:"last_completed_at,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (s Status) Valid() bool {
	switch s {
	case StatusNew, StatusInProgress, StatusDone:
		return true
	default:
		return false
	}
}
