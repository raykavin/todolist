package module

import (
	"todolist/internal/config"
	rptPerson "todolist/internal/domain/person/repository"
	rptTodo "todolist/internal/domain/todo/repository"
	rptUser "todolist/internal/domain/user/repository"

	svcTodo "todolist/internal/domain/todo/service"

	ucPerson "todolist/internal/usecase/person"
	ucTodo "todolist/internal/usecase/todo"
	ucUser "todolist/internal/usecase/user"
	"todolist/pkg/log"

	"go.uber.org/fx"
)

// UseCaseParams defines the dependencies required to create use cases
type (
	UseCaseParams struct {
		fx.In
		AppConfig           config.ApplicationProvider
		PersonRepository    rptPerson.PersonRepository
		UserRepository      rptUser.UserRepository
		TodoRepository      rptTodo.TodoRepository
		TodoQueryRepository rptTodo.TodoQueryRepository
		TodoService         svcTodo.TodoService

		Log log.Interface
	}

	UseCaseContainer struct {
		fx.Out
		// Person Use Cases
		CreatePersonUseCase ucPerson.CreatePersonUseCase
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
)

func NewUseCases(params UseCaseParams) UseCaseContainer {
	return UseCaseContainer{
		CreatePersonUseCase: ucPerson.NewCreatePersonUseCase(params.PersonRepository),
		GetPersonUseCase:    ucPerson.NewGetPersonUseCase(params.PersonRepository),

		ChangePasswordUseCase: ucUser.NewChangePasswordUseCase(params.UserRepository),
		CreateUserUseCase:     ucUser.NewCreateUserUseCase(params.UserRepository, params.PersonRepository),
		LoginUseCase:          ucUser.NewLoginUseCase(params.UserRepository, params.PersonRepository, params.AppConfig.GetName()),

		CompleteTodoUseCase:  ucTodo.NewCompleteTodoUseCase(params.TodoRepository, params.TodoService),
		CreateTodoUseCase:    ucTodo.NewCreateTodoUseCase(params.TodoRepository, params.TodoService),
		DeleteTodoUseCase:    ucTodo.NewDeleteTodoUseCase(params.TodoRepository, params.TodoService),
		GetStatisticsUseCase: ucTodo.NewGetStatisticsUseCase(params.TodoQueryRepository),
		GetTodoUseCase:       ucTodo.NewGetTodoUseCase(params.TodoRepository, params.TodoService),
		ListTodoUseCase:      ucTodo.NewListTodosUseCase(params.TodoQueryRepository),
		UpdateTodoUseCase:    ucTodo.NewUpdateTodoUseCase(params.TodoRepository, params.TodoService),
	}
}

// UseCase returns the fx module with all core dependencies
func UseCase() fx.Option {
	return fx.Module("usecases", fx.Provide(NewUseCases))
}
