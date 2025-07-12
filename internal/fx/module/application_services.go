package module

import (
	"go.uber.org/fx"

	rptTodo "todolist/internal/domain/user/repository"
	"todolist/internal/service"
)

// ApplicationServiceParams defines the dependencies required to create services
type ApplicationServiceParams struct {
	fx.In
	UserRepository      rptTodo.UserRepository
	UserQueryRepository rptTodo.UserQueryRepository
}

// ApplicationServiceContainer provides all service implementations
type ApplicationServiceContainer struct {
	fx.Out
	UserSecurityService service.UserSecurityService
}

// NewApplicationServices creates all service implementations
func NewApplicationServices(p ApplicationServiceParams) ApplicationServiceContainer {
	return ApplicationServiceContainer{
		UserSecurityService: service.NewUserSecurityService(p.UserRepository, p.UserQueryRepository),
	}
}

// ApplicationServices returns the fx module with all service dependencies
func ApplicationServices() fx.Option {
	return fx.Module("domain_services",
		fx.Provide(NewApplicationServices),
	)
}
