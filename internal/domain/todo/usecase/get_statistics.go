package usecase

import (
	"context"
	"todolist/internal/domain/todo/dto"
	"todolist/internal/domain/todo/repository"
)

// GetStatisticsUseCase handles retrieving todo statistics
type GetStatisticsUseCase interface {
	Execute(ctx context.Context, userID int64) (*dto.TodoStatistics, error)
}

type getStatisticsUseCase struct {
	todoQueryRepo repository.TodoQueryRepository
}

// NewGetStatisticsUseCase creates a new instance of GetStatisticsUseCase
func NewGetStatisticsUseCase(todoQueryRepo repository.TodoQueryRepository) GetStatisticsUseCase {
	return &getStatisticsUseCase{
		todoQueryRepo: todoQueryRepo,
	}
}

// Execute retrieves todo statistics for a user
func (uc *getStatisticsUseCase) Execute(ctx context.Context, userID int64) (*dto.TodoStatistics, error) {
	return uc.todoQueryRepo.GetStatistics(ctx, userID)
}
