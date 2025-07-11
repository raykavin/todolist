package module

import (
	rptTodo "todolist/internal/domain/todo/repository"
	svcTodo "todolist/internal/domain/todo/service"

	"go.uber.org/fx"
)

type ServiceParams struct {
	fx.In

	TodoRepository      rptTodo.TodoRepository
	TodoQueryRepository rptTodo.TodoQueryRepository
}

type ServiceContainer struct {
	fx.Out
	svcTodo.TodoService
}

func NewServices(params ServiceParams) ServiceContainer {
	return ServiceContainer{
		TodoService: svcTodo.NewTodoService(params.TodoRepository, params.TodoQueryRepository),
	}
}
