package di

import (
	"todolist/internal/adapter/repository"
	rptPerson "todolist/internal/domain/person/repository"
	rptTodo "todolist/internal/domain/todo/repository"
	rptUser "todolist/internal/domain/user/repository"

	"go.uber.org/fx"
	"gorm.io/gorm"
)

// RepositoryParams defines the dependencies required to create the repositories
type RepositoryParams struct {
	fx.In
	DatabaseProvider *gorm.DB
}

// Repositories groups all repository implementations provided from Fx
type RepositoryContainer struct {
	fx.Out
	UserRepository        rptUser.UserRepository
	UserQueryRepository   rptUser.UserQueryRepository
	PersonRepository      rptPerson.PersonRepository
	PersonQueryRepository rptPerson.PersonQueryRepository
	TodoRepository        rptTodo.TodoRepository
	TodoQueryRepository   rptTodo.TodoQueryRepository
}

// NewRepositories creates all repository implementations
func NewRepositories(p RepositoryParams) RepositoryContainer {
	return RepositoryContainer{
		UserRepository:        repository.NewUserRepository(p.DatabaseProvider),
		UserQueryRepository:   repository.NewUserQueryRepository(p.DatabaseProvider),
		PersonRepository:      repository.NewPersonRepository(p.DatabaseProvider),
		PersonQueryRepository: repository.NewPersonQueryRepository(p.DatabaseProvider),
		TodoRepository:        repository.NewTodoRepository(p.DatabaseProvider),
		TodoQueryRepository:   repository.NewTodoQueryRepository(p.DatabaseProvider),
	}
}

// RepositoriesModule exports the Fx module that provides all repository dependencies
func RepositoriesModule() fx.Option {
	return fx.Module("repositories", fx.Provide(NewRepositories))
}
