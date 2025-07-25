package di

import (
	"todolist/internal/adapter/delivery/http/handler"
	"todolist/internal/config"
	ucPerson "todolist/internal/usecase/person"
	ucTodo "todolist/internal/usecase/todo"
	ucUser "todolist/internal/usecase/user"

	"go.uber.org/fx"
)

// HttpHandlerParams defines the dependencies required to create the http handlers
type HttpHandlerParams struct {
	fx.In
	// Application configuration
	AppConfig config.ApplicationProvider

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

// HttpHandlerContainer groups all http handlers implementations provide from Fx
type HttpHandlerContainer struct {
	fx.Out
	AuthHandler   *handler.AuthHandler
	PersonHandler *handler.PersonHandler
	TodoHandler   *handler.TodoHandler
	HealthHandler *handler.HealthHandler
}

// NewHttpHandlers creates all http handlers implementations
func NewHttpHandlers(p HttpHandlerParams) HttpHandlerContainer {
	return HttpHandlerContainer{
		AuthHandler: handler.NewAuthHandler(
			p.CreateUserUseCase,
			p.CreatePersonUseCase,
			p.LoginUseCase,
			p.ChangePasswordUseCase,
		),
		PersonHandler: handler.NewPersonHandler(
			p.CreatePersonUseCase,
			p.UpdatePersonUseCase,
			p.GetPersonUseCase,
		),
		TodoHandler: handler.NewTodoHandler(
			p.CreateTodoUseCase,
			p.UpdateTodoUseCase,
			p.CompleteTodoUseCase,
			p.DeleteTodoUseCase,
			p.GetTodoUseCase,
			p.ListTodoUseCase,
			p.GetStatisticsUseCase,
		),
		HealthHandler: handler.NewHealthHandler(p.AppConfig)}
}

// HTTPHandlersModule exports the Fx module that provides all http handlers dependencies
func HTTPHandlersModule() fx.Option {
	return fx.Module("http_handlers", fx.Provide(NewHttpHandlers))
}
