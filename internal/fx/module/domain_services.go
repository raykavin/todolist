package module

import (
	"go.uber.org/fx"

	rptTodo "todolist/internal/domain/todo/repository"
	svcTodo "todolist/internal/domain/todo/service"
)

// DomainServiceParams defines the dependencies required to create services
type DomainServiceParams struct {
	fx.In
	TodoRepository      rptTodo.TodoRepository
	TodoQueryRepository rptTodo.TodoQueryRepository
}

// DomainServiceContainer provides all service implementations
type DomainServiceContainer struct {
	fx.Out
	TodoService svcTodo.TodoService
}

// NewDomainServices creates all service implementations
func NewDomainServices(p DomainServiceParams) DomainServiceContainer {
	return DomainServiceContainer{
		TodoService: svcTodo.NewTodoService(p.TodoRepository, p.TodoQueryRepository),
	}
}

// DomainServices returns the fx module with all service dependencies
func DomainServices() fx.Option {
	return fx.Module("domain_services",
		fx.Provide(NewDomainServices),
	)
}
