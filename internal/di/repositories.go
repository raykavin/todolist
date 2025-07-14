package di

import (
	rptPerson "todolist/internal/domain/person/repository"
	rptTodo "todolist/internal/domain/todo/repository"
	rptUser "todolist/internal/domain/user/repository"
	rptInfra "todolist/internal/infrastructure/database/repository"

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
		UserRepository:        rptInfra.NewUserRepository(p.DatabaseProvider),
		UserQueryRepository:   rptInfra.NewUserQueryRepository(p.DatabaseProvider),
		PersonRepository:      rptInfra.NewPersonRepository(p.DatabaseProvider),
		PersonQueryRepository: rptInfra.NewPersonQueryRepository(p.DatabaseProvider),
		TodoRepository:        rptInfra.NewTodoRepository(p.DatabaseProvider),
		TodoQueryRepository:   rptInfra.NewTodoQueryRepository(p.DatabaseProvider),
	}
}

// RepositoriesModule exports the Fx module that provides all repository dependencies
func RepositoriesModule() fx.Option {
	return fx.Module("repositories", fx.Provide(NewRepositories))
}
