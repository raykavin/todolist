package usecase

import (
	"context"
	"todolist/internal/domain/todo/repository"
	vo "todolist/internal/domain/todo/valueobject"
)

// GetStatisticsUseCase handles retrieving todo statistics
type GetStatisticsUseCase interface {
	Execute(ctx context.Context, userID int64) (*vo.TodoStatistics, error)
}

type getStatisticsUseCase struct {
	todoQueryRepository repository.TodoQueryRepository
}

// NewGetStatisticsUseCase creates a new instance of GetStatisticsUseCase
func NewGetStatisticsUseCase(todoQueryRepository repository.TodoQueryRepository) GetStatisticsUseCase {
	return &getStatisticsUseCase{
		todoQueryRepository: todoQueryRepository,
	}
}

// Execute retrieves todo statistics for a user
func (uc *getStatisticsUseCase) Execute(ctx context.Context, userID int64) (*vo.TodoStatistics, error) {
	return uc.todoQueryRepository.GetStatistics(ctx, userID)
}
