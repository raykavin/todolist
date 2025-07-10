package module

import (
	"todolist/internal/config"
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
		AppConfig             config.ApplicationProvider
		DatabaseServiceConfig config.DatabaseServiceProvider
		Log                   log.Interface
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
