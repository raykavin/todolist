package module

import (
	"go.uber.org/fx"

	"todolist/internal/config"
	rptPerson "todolist/internal/domain/person/repository"
	rptTodo "todolist/internal/domain/todo/repository"
	svcTodo "todolist/internal/domain/todo/service"
	rptUser "todolist/internal/domain/user/repository"
	ucPerson "todolist/internal/usecase/person"
	ucTodo "todolist/internal/usecase/todo"
	ucUser "todolist/internal/usecase/user"
	"todolist/pkg/log"
)

// UseCaseParams defines the dependencies required to create use cases
type UseCaseParams struct {
	fx.In
	AppConfig           config.ApplicationProvider
	PersonRepository    rptPerson.PersonRepository
	UserRepository      rptUser.UserRepository
	TodoRepository      rptTodo.TodoRepository
	TodoQueryRepository rptTodo.TodoQueryRepository
	TodoService         svcTodo.TodoService
	Log                 log.Interface
}

// UseCaseContainer provides all use case implementations
type UseCaseContainer struct {
	fx.Out

	// Person Use Cases
	CreatePersonUseCase ucPerson.CreatePersonUseCase
	UpdatePersonUseCase ucPerson.UpdatePersonUseCase
	GetPersonUseCase    ucPerson.GetPersonUseCase

	// User Use Cases
	ChangePasswordUseCase ucUser.ChangePasswordUseCase
	CreateUserUseCase     ucUser.CreateUserUseCase
	LoginUseCase          ucUser.LoginUseCase

	// Todo Use Cases
	CompleteTodoUseCase  ucTodo.CompleteTodoUseCase
	CreateTodoUseCase    ucTodo.CreateTodoUseCase
	DeleteTodoUseCase    ucTodo.DeleteTodoUseCase
	GetStatisticsUseCase ucTodo.GetStatisticsUseCase
	GetTodoUseCase       ucTodo.GetTodoUseCase
	ListTodoUseCase      ucTodo.ListTodosUseCase
	UpdateTodoUseCase    ucTodo.UpdateTodoUseCase
}

// NewUseCases creates all use case implementations
func NewUseCases(p UseCaseParams) UseCaseContainer {
	return UseCaseContainer{
		// Person Use Cases
		CreatePersonUseCase: ucPerson.NewCreatePersonUseCase(p.PersonRepository),
		GetPersonUseCase:    ucPerson.NewGetPersonUseCase(p.PersonRepository),

		// User Use Cases
		ChangePasswordUseCase: ucUser.NewChangePasswordUseCase(p.UserRepository),
		CreateUserUseCase:     ucUser.NewCreateUserUseCase(p.UserRepository, p.PersonRepository),
		LoginUseCase:          ucUser.NewLoginUseCase(p.UserRepository, p.PersonRepository, p.AppConfig.GetName()),

		// Todo Use Cases
		CompleteTodoUseCase:  ucTodo.NewCompleteTodoUseCase(p.TodoRepository, p.TodoService),
		CreateTodoUseCase:    ucTodo.NewCreateTodoUseCase(p.TodoRepository, p.TodoService),
		DeleteTodoUseCase:    ucTodo.NewDeleteTodoUseCase(p.TodoRepository, p.TodoService),
		GetStatisticsUseCase: ucTodo.NewGetStatisticsUseCase(p.TodoQueryRepository),
		GetTodoUseCase:       ucTodo.NewGetTodoUseCase(p.TodoRepository, p.TodoService),
		ListTodoUseCase:      ucTodo.NewListTodosUseCase(p.TodoQueryRepository),
		UpdateTodoUseCase:    ucTodo.NewUpdateTodoUseCase(p.TodoRepository, p.TodoService),
	}
}

// UseCases returns the fx module with all use case dependencies
func UseCases() fx.Option {
	return fx.Module("usecases",
		fx.Provide(NewUseCases),
	)
}
