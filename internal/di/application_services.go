package di

import (
	"fmt"

	"go.uber.org/fx"

	"todolist/internal/config"
	rptTodo "todolist/internal/domain/user/repository"
	"todolist/internal/infrastructure/auth"
	"todolist/internal/service"
)

// ApplicationServiceParams defines the dependencies required to create services
type ApplicationServiceParams struct {
	fx.In
	UserRepository      rptTodo.UserRepository
	UserQueryRepository rptTodo.UserQueryRepository
	AppConfig           config.ApplicationProvider
}

// ApplicationServiceContainer provides all service implementations
type ApplicationServiceContainer struct {
	fx.Out
	UserSecurityService service.UserSecurityService
	TokenService        service.TokenService
}

// NewApplicationServices creates all service implementations
func NewApplicationServices(p ApplicationServiceParams) (ApplicationServiceContainer, error) {
	jwtConfig := p.AppConfig.GetJWT()
	secretKey := jwtConfig.GetSecretKey()
	accessDuration := jwtConfig.GetExpirationTime()
	refreshDuration := jwtConfig.GetRefreshExpirationTime()

	// Create JWT token service using the adapter
	tokenService, err := auth.NewJWTTokenAdapter(
		secretKey,
		accessDuration,
		refreshDuration,
	)
	if err != nil {
		return ApplicationServiceContainer{}, fmt.Errorf("failed to initialize token service: %w", err)
	}

	return ApplicationServiceContainer{
		UserSecurityService: service.NewUserSecurityService(p.UserRepository, p.UserQueryRepository),
		TokenService:        tokenService,
	}, nil
}

// ApplicationServicesModule returns the fx module with all service dependencies
func ApplicationServicesModule() fx.Option {
	return fx.Module("domain_services",
		fx.Provide(NewApplicationServices),
	)
}
