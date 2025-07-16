package di

import (
	"context"
	"fmt"
	"sync"
	"time"

	adptHttp "todolist/internal/adapter/delivery/http"
	"todolist/internal/adapter/delivery/http/handler"
	"todolist/internal/adapter/delivery/http/middleware"
	"todolist/internal/config"
	"todolist/internal/service"
	"todolist/pkg/log"
	"todolist/pkg/web"

	"github.com/gin-gonic/gin"
	"go.uber.org/fx"

	_ "todolist/docs"
)

// HTTPServerParams defines the dependencies required to create the HTTP server
type HTTPServerParams struct {
	fx.In
	Context       context.Context
	WaitGroup     *sync.WaitGroup
	AuthHandler   *handler.AuthHandler
	PersonHandler *handler.PersonHandler
	TodoHandler   *handler.TodoHandler
	HealthHandler *handler.HealthHandler
	TokenService  service.TokenService
	Log           log.ExtendedLog
	AppConfig     config.ApplicationProvider
}

// HTTPServerContainer provides the HTTP server components
type HTTPServerContainer struct {
	fx.Out
	Engine adptHttp.HttpServer
	Router *gin.Engine
}

// NewHTTPServer creates and configures the Gin engine and routes
func NewHTTPServer(params HTTPServerParams) (HTTPServerContainer, error) {
	// Create web config from application config
	engineConfig := createEngineConfig(params.AppConfig)

	// Create the web engine
	engine, err := web.NewGinWithConfig(engineConfig)
	if err != nil {
		return HTTPServerContainer{}, fmt.Errorf("failed to create HTTP engine: %w", err)
	}

	// Setup custom middleware
	engine.SetupMiddleware(func(router *gin.Engine) {
		setupAppMiddleware(router, params)
	})

	// Setup routes
	engine.SetupRoutes(func(router *gin.Engine) {
		setupAppRoutes(router, params)
	})

	return HTTPServerContainer{
		Engine: engine,
		Router: engine.Router(),
	}, nil
}

// createEngineConfig creates web.Config from application config
func createEngineConfig(appConfig config.ApplicationProvider) *web.Config {
	webConfig := appConfig.GetWeb()

	port := webConfig.GetListen()
	if port == 0 {
		port = 8080
	}

	isDebugMode := appConfig.GetLogLevel() == "debug"
	isNoRouteJSON := webConfig.GetNoRouteTo() == ""

	return &web.Config{
		Host:           "",
		Port:           port,
		DebugMode:      isDebugMode,
		EnableHTTP2:    true,
		UseRecovery:    true,
		NoRouteJSON:    isNoRouteJSON,
		NoRouteTo:      webConfig.GetNoRouteTo(),
		ReadTimeout:    webConfig.GetReadTimeout(),
		WriteTimeout:   webConfig.GetWriteTimeout(),
		IdleTimeout:    webConfig.GetIdleTimeout(),
		UseSSL:         webConfig.GetUseSSL(),
		SSLCert:        webConfig.GetSSLCert(),
		SSLKey:         webConfig.GetSSLKey(),
		MaxPayloadSize: webConfig.GetMaxPayloadSize(),
	}
}

// setupAppMiddleware configures application-specific middleware
func setupAppMiddleware(router *gin.Engine, params HTTPServerParams) {
	webConfig := params.AppConfig.GetWeb()

	// Request logging
	router.Use(middleware.RequestLog(params.Log))

	// CORS
	if corsHeaders := webConfig.GetCORS(); len(corsHeaders) > 0 {
		router.Use(middleware.CORS(corsHeaders))
	} else {
		router.Use(middleware.CORS())
	}
}

// setupAppRoutes configures all application routes
func setupAppRoutes(router *gin.Engine, params HTTPServerParams) {
	// Basic routes
	router.GET("/", func(c *gin.Context) {
		c.Redirect(302, "/swagger/index.html")
	})

	router.GET("/health", adptHttp.WrapHandler(params.HealthHandler.HealthCheck))

	// API v1 routes
	authMiddleware := middleware.AuthMiddleware(params.TokenService)
	v1 := router.Group("/api/v1")

	// Authentication routes (public and mixed)
	auth := v1.Group("/auth")
	{
		auth.POST("/register", adptHttp.WrapHandler(params.AuthHandler.Register))
		auth.POST("/login", adptHttp.WrapHandler(params.AuthHandler.Login))
		auth.PUT("/change-password", authMiddleware, adptHttp.WrapHandler(params.AuthHandler.ChangePassword))
	}

	// Protected routes
	protected := v1.Group("", authMiddleware)
	{
		// People management
		people := protected.Group("/people")
		{
			people.POST("", adptHttp.WrapHandler(params.PersonHandler.CreatePerson))
			people.GET("/:id", adptHttp.WrapHandler(params.PersonHandler.GetPerson))
			people.PUT("/:id", adptHttp.WrapHandler(params.PersonHandler.UpdatePerson))
		}

		// Todo management
		todos := protected.Group("/todos")
		{
			todos.POST("", adptHttp.WrapHandler(params.TodoHandler.CreateTodo))
			todos.GET("", adptHttp.WrapHandler(params.TodoHandler.ListTodos))
			todos.GET("/statistics", adptHttp.WrapHandler(params.TodoHandler.GetStatistics))
			todos.GET("/:id", adptHttp.WrapHandler(params.TodoHandler.GetTodo))
			todos.PUT("/:id", adptHttp.WrapHandler(params.TodoHandler.UpdateTodo))
			todos.PUT("/:id/complete", adptHttp.WrapHandler(params.TodoHandler.CompleteTodo))
			todos.DELETE("/:id", adptHttp.WrapHandler(params.TodoHandler.DeleteTodo))
		}
	}
}

// httpServerLifecycle manages starting and stopping the HTTP server
func httpServerLifecycle(lc fx.Lifecycle, engine adptHttp.HttpServer, params HTTPServerParams) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			return startServer(engine, params)
		},
		OnStop: func(ctx context.Context) error {
			return stopServer(ctx, engine, params.Log)
		},
	})
}

// startServer starts the HTTP server
func startServer(engine adptHttp.HttpServer, params HTTPServerParams) error {
	params.WaitGroup.Add(1)

	go func() {
		defer params.WaitGroup.Done()

		// Get the actual engine to access config
		if webEngine, ok := engine.(*web.Engine); ok {
			addr := webEngine.Addr()
			isSSL := webEngine.IsSSLEnabled()

			if isSSL {
				params.Log.Info(fmt.Sprintf("Starting HTTPS server on %s", addr))
			} else {
				params.Log.Info(fmt.Sprintf("Starting HTTP server on %s", addr))
			}
		}

		if err := engine.Listen(); err != nil {
			params.Log.Fatalf("Server failed to start: %v", err)
		}
	}()

	// Give server time to start
	time.Sleep(100 * time.Millisecond)
	params.Log.Success("HTTP server started successfully")
	return nil
}

// stopServer gracefully shuts down the HTTP server
func stopServer(ctx context.Context, engine adptHttp.HttpServer, logger log.ExtendedLog) error {
	logger.Info("Stopping HTTP server...")

	shutdownCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	if err := engine.Shutdown(shutdownCtx); err != nil {
		logger.Failure(fmt.Sprintf("HTTP server shutdown error: %v", err))
		return err
	}

	logger.Success("HTTP server stopped successfully")
	return nil
}

// HTTPServerModule wires the HTTP server to the Fx container
func HTTPServerModule() fx.Option {
	return fx.Module("http_server",
		fx.Provide(NewHTTPServer),
		fx.Invoke(httpServerLifecycle),
	)
}
