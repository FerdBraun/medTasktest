package task

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	taskdomain "example.com/taskservice/internal/domain/task"
)

type Service struct {
	repo Repository
	now  func() time.Time
}

const MockTimeKey = "mock_time"

func (s *Service) getNow(ctx context.Context) time.Time {
	if t, ok := ctx.Value(MockTimeKey).(time.Time); ok {
		return t
	}
	return s.now()
}

func NewService(repo Repository) *Service {
	return &Service{
		repo: repo,
		now:  func() time.Time { return time.Now().UTC() },
	}
}

func (s *Service) Create(ctx context.Context, input CreateInput) (*taskdomain.Task, error) {
	normalized, err := validateCreateInput(input)
	if err != nil {
		return nil, err
	}

	model := &taskdomain.Task{
		Title:            normalized.Title,
		Description:      normalized.Description,
		Status:           normalized.Status,
		PeriodicityType:  normalized.PeriodicityType,
		PeriodicityValue: normalized.PeriodicityValue,
	}
	now := s.getNow(ctx)
	model.CreatedAt = now
	model.UpdatedAt = now

	created, err := s.repo.Create(ctx, model)
	if err != nil {
		return nil, err
	}

	return created, nil
}

func (s *Service) GetByID(ctx context.Context, id int64) (*taskdomain.Task, error) {
	if id <= 0 {
		return nil, fmt.Errorf("%w: id must be positive", ErrInvalidInput)
	}

	task, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return s.refreshPeriodicity(ctx, task)
}

func (s *Service) Update(ctx context.Context, id int64, input UpdateInput) (*taskdomain.Task, error) {
	if id <= 0 {
		return nil, fmt.Errorf("%w: id must be positive", ErrInvalidInput)
	}

	normalized, err := validateUpdateInput(input)
	if err != nil {
		return nil, err
	}

	existing, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	model := &taskdomain.Task{
		ID:               id,
		Title:            normalized.Title,
		Description:      normalized.Description,
		Status:           normalized.Status,
		PeriodicityType:  normalized.PeriodicityType,
		PeriodicityValue: normalized.PeriodicityValue,
		CreatedAt:        existing.CreatedAt,
		LastCompletedAt:  existing.LastCompletedAt,
		UpdatedAt:        s.getNow(ctx),
	}

	// Update last_completed_at if status changed to "done"
	if model.Status == taskdomain.StatusDone && existing.Status != taskdomain.StatusDone {
		now := s.getNow(ctx)
		model.LastCompletedAt = &now
	}

	updated, err := s.repo.Update(ctx, model)
	if err != nil {
		return nil, err
	}

	return updated, nil
}

func (s *Service) Delete(ctx context.Context, id int64) error {
	if id <= 0 {
		return fmt.Errorf("%w: id must be positive", ErrInvalidInput)
	}

	return s.repo.Delete(ctx, id)
}

func (s *Service) List(ctx context.Context) ([]taskdomain.Task, error) {
	list, err := s.repo.List(ctx)
	if err != nil {
		return nil, err
	}

	for i := range list {
		refreshed, err := s.refreshPeriodicity(ctx, &list[i])
		if err == nil {
			list[i] = *refreshed
		}
	}

	return list, nil
}

func (s *Service) refreshPeriodicity(ctx context.Context, task *taskdomain.Task) (*taskdomain.Task, error) {
	if task.PeriodicityType == "no" || task.Status != taskdomain.StatusDone || task.LastCompletedAt == nil {
		return task, nil
	}

	now := s.getNow(ctx)
	shouldReset := false
	todayStart := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)

	switch task.PeriodicityType {
	case "day":
		interval, _ := strconv.Atoi(task.PeriodicityValue)
		if interval <= 0 {
			interval = 1
		}
		// Compare start of days to reset at midnight of the target day
		lastDayStart := time.Date(task.LastCompletedAt.Year(), task.LastCompletedAt.Month(), task.LastCompletedAt.Day(), 0, 0, 0, 0, time.UTC)
		daysPassed := int(todayStart.Sub(lastDayStart).Hours() / 24)
		if daysPassed >= interval {
			shouldReset = true
		}
	case "month":
		day, _ := strconv.Atoi(task.PeriodicityValue)
		// Reset if it's the target day of the month AND it wasn't completed today
		if now.Day() == day && task.LastCompletedAt.Before(todayStart) {
			shouldReset = true
		}
	case "date":
		dates := strings.Split(task.PeriodicityValue, ",")
		todayISO := now.Format("2006-01-02")
		todayRU := now.Format("02.01.2006")
		for _, d := range dates {
			d = strings.TrimSpace(d)
			// Match either ISO (2006-01-02) or RU (02.01.2006) format
			if (d == todayISO || d == todayRU) && task.LastCompletedAt.Before(todayStart) {
				shouldReset = true
				break
			}
		}
	case "odd":
		isEven := now.Day()%2 == 0
		val := strings.ToLower(task.PeriodicityValue)
		if ((val == "even" && isEven) || (val == "odd" && !isEven)) &&
			task.LastCompletedAt.Before(todayStart) {
			shouldReset = true
		}
	}

	if shouldReset {
		task.Status = taskdomain.StatusNew
		task.UpdatedAt = s.getNow(ctx)
		return s.repo.Update(ctx, task)
	}

	return task, nil
}

func validateCreateInput(input CreateInput) (CreateInput, error) {
	input.Title = strings.TrimSpace(input.Title)
	input.Description = strings.TrimSpace(input.Description)

	if input.Title == "" {
		return CreateInput{}, fmt.Errorf("%w: title is required", ErrInvalidInput)
	}

	if input.Status == "" {
		input.Status = taskdomain.StatusNew
	}

	if !input.Status.Valid() {
		return CreateInput{}, fmt.Errorf("%w: invalid status", ErrInvalidInput)
	}

	if input.PeriodicityType == "" {
		input.PeriodicityType = "no"
	}

	return input, nil
}

func validateUpdateInput(input UpdateInput) (UpdateInput, error) {
	input.Title = strings.TrimSpace(input.Title)
	input.Description = strings.TrimSpace(input.Description)

	if input.Title == "" {
		return UpdateInput{}, fmt.Errorf("%w: title is required", ErrInvalidInput)
	}

	if !input.Status.Valid() {
		return UpdateInput{}, fmt.Errorf("%w: invalid status", ErrInvalidInput)
	}

	if input.PeriodicityType == "" {
		input.PeriodicityType = "no"
	}

	return input, nil
}
