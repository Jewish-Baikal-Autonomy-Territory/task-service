package task

import (
	"context"
	"errors"
	"fmt"
	domaintask "task-service/internal/domain/task"

	"github.com/samber/mo"
)

var (
	ErrMismatchedOwner = errors.New("owner is not the original owner")
)

type Service interface {
	CreateTask(ctx context.Context, command *CreateTaskCommand) (*domaintask.Task, error)
	CompleteTask(ctx context.Context, command CompleteTaskCommand) error
	UpdateTask(ctx context.Context, command *UpdateTaskCommand) error
	DeleteTask(ctx context.Context, command DeleteTaskCommand) error
	RestoreTask(ctx context.Context, command RestoreTaskCommand) error
}

type service struct {
	repository domaintask.Repository
}

func (s *service) CreateTask(ctx context.Context, command *CreateTaskCommand) (*domaintask.Task, error) {
	if command == nil {
		return nil, errors.New("command is missing")
	}
	task, err := domaintask.NewTask(command.OwnerID, command.Title, command.Description, command.Priority)
	if err != nil {
		return nil, fmt.Errorf("create task: %w", err)
	}
	task.GroupID = command.GroupID
	task.IsFavorite = command.IsFavorite

	if err = s.repository.Create(ctx, task); err != nil {
		return nil, fmt.Errorf("create task: %w", err)
	}
	return task, nil
}

func (s *service) CompleteTask(ctx context.Context, command CompleteTaskCommand) error {
	task, err := s.repository.FindByID(ctx, command.ID)
	if err != nil {
		return fmt.Errorf("find task: %w", err)
	}
	if task.OwnerID != command.OwnerID {
		return ErrMismatchedOwner
	}
	if err = task.Complete(); err != nil {
		return fmt.Errorf("complete task: %w", err)
	}
	if err = s.repository.Update(ctx, task); err != nil {
		return fmt.Errorf("update task: %w", err)
	}
	return nil
}

func (s *service) UpdateTask(ctx context.Context, command *UpdateTaskCommand) error {
	if command == nil {
		return errors.New("command is missing")
	}

	task, err := s.repository.FindByID(ctx, command.ID)
	if err != nil {
		return fmt.Errorf("find task: %w", err)
	}
	if task.OwnerID != command.OwnerID {
		return errors.New("invalid owner")
	}
	if value, ok := command.GroupID.Get(); ok {
		task.GroupID = mo.Some(value)
	}
	if value, ok := command.Title.Get(); ok {
		task.Title = value
	}
	if value, ok := command.Description.Get(); ok {
		task.Description = value
	}
	if value, ok := command.Priority.Get(); ok {
		task.Priority = value
	}
	if value, ok := command.IsFavorite.Get(); ok {
		task.IsFavorite = value
	}
	err = s.repository.Update(ctx, task)
	if err != nil {
		return fmt.Errorf("update task: %w", err)
	}

	return nil
}

func (s *service) DeleteTask(ctx context.Context, command DeleteTaskCommand) error {
	task, err := s.repository.FindByID(ctx, command.ID)
	if err != nil {
		return fmt.Errorf("find task: %w", err)
	}
	if task.OwnerID != command.OwnerID {
		return errors.New("is not owner of task")
	}
	if err = task.SoftDelete(); err != nil {
		return fmt.Errorf("soft delete task: %w", err)
	}
	if err = s.repository.Update(ctx, task); err != nil {
		return fmt.Errorf("update task: %w", err)
	}
	return nil
}

func (s *service) RestoreTask(ctx context.Context, command RestoreTaskCommand) error {
	task, err := s.repository.FindByID(ctx, command.ID)
	if err != nil {
		return fmt.Errorf("find task: %w", err)
	}
	if task.OwnerID != command.OwnerID {
		return errors.New("is not owner of task")
	}
	if err = task.Restore(); err != nil {
		return fmt.Errorf("restore task: %w", err)
	}
	if err = s.repository.Update(ctx, task); err != nil {
		return fmt.Errorf("update task: %w", err)
	}
	return nil
}

func NewService(repository domaintask.Repository) (Service, error) {
	if repository == nil {
		return nil, errors.New("repository is required")
	}
	return &service{
		repository: repository,
	}, nil
}
