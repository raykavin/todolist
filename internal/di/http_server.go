package di

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
	"sync"
	"time"

	adptHttp "todolist/internal/adapter/http"
	"todolist/internal/adapter/http/middleware"
	"todolist/internal/config"
	"todolist/internal/http/handler"
	"todolist/pkg/log"
	"todolist/pkg/web"

	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
	"golang.org/x/net/http2"
)

// HTTPServerParams defines the dependencies required to create the HTTP server
type HTTPServerParams struct {
	fx.In
	Context       context.Context
	WaitGroup     *sync.WaitGroup
	AuthHandler   *handler.AuthHandler
	PersonHandler *handler.PersonHandler
	TodoHandler   *handler.TodoHandler
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
	webConfig := params.AppConfig.GetWeb()

	port := webConfig.GetListen()
	if port == 0 {
		port = 8080
	}

	debugMode := params.AppConfig.GetLoggerLevel() == "debug"

	engine, err := web.NewGin("", port, debugMode)
	if err != nil {
		return HTTPServerContainer{}, fmt.Errorf("failed to create HTTP engine: %w", err)
	}

	router := engine.Router()
	registerMiddleware(router, webConfig, params.Log)
	registerRoutes(router, params)

	return HTTPServerContainer{
		Engine: engine,
		Router: router,
	}, nil
}

// registerMiddleware attaches all middleware to the Gin router
func registerMiddleware(router *gin.Engine, webConfig config.WebConfigProvider, logger log.ExtendedLog) {
	router.Use(middleware.RequestLog(logger))
	router.Use(middleware.Recovery())

	// CORS
	if corsHeaders := webConfig.GetCORS(); len(corsHeaders) > 0 {
		router.Use(middleware.CORS(corsHeaders))
	} else {
		router.Use(middleware.CORS())
	}

	// Max payload size
	if maxSize := webConfig.GetMaxPayloadSize(); maxSize > 0 {
		router.Use(func(c *gin.Context) {
			c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, maxSize)
			c.Next()
		})
	}
}

// registerRoutes defines all HTTP routes for the API
func registerRoutes(router *gin.Engine, params HTTPServerParams) {
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"app":     params.AppConfig.GetName(),
			"version": params.AppConfig.GetVersion(),
		})
	})

	v1 := router.Group("/api/v1")
	{
		auth := v1.Group("/auth")
		auth.POST("/register", adptHttp.WrapHandler(params.AuthHandler.Register))
		auth.POST("/login", adptHttp.WrapHandler(params.AuthHandler.Login))

		// Protected routes
		protected := v1.Group("")
		protected.Use(middleware.AuthMiddleware(middleware.JWTConfig{Secret: "your-secret-key"}))

		users := protected.Group("/users")
		users.PUT("/password", adptHttp.WrapHandler(params.AuthHandler.ChangePassword))

		persons := protected.Group("/persons")
		persons.POST("", adptHttp.WrapHandler(params.PersonHandler.CreatePerson))
		persons.GET("/:id", adptHttp.WrapHandler(params.PersonHandler.GetPerson))
		persons.PUT("/:id", adptHttp.WrapHandler(params.PersonHandler.UpdatePerson))

		todos := protected.Group("/todos")
		todos.POST("", adptHttp.WrapHandler(params.TodoHandler.CreateTodo))
		todos.GET("", adptHttp.WrapHandler(params.TodoHandler.ListTodos))
		todos.GET("/statistics", adptHttp.WrapHandler(params.TodoHandler.GetStatistics))
		todos.GET("/:id", adptHttp.WrapHandler(params.TodoHandler.GetTodo))
		todos.PUT("/:id", adptHttp.WrapHandler(params.TodoHandler.UpdateTodo))
		todos.PUT("/:id/complete", adptHttp.WrapHandler(params.TodoHandler.CompleteTodo))
		todos.DELETE("/:id", adptHttp.WrapHandler(params.TodoHandler.DeleteTodo))
	}

	router.NoRoute(func(c *gin.Context) {
		if redirectTo := params.AppConfig.GetWeb().GetNoRouteTo(); redirectTo != "" {
			c.Redirect(302, redirectTo)
		} else {
			c.JSON(404, gin.H{"error": "Not Found", "message": "The requested resource was not found"})
		}
	})
}

// buildHTTPServer constructs the http.Server with optional TLS and HTTP/2
func buildHTTPServer(engine *gin.Engine, webConfig config.WebConfigProvider) *http.Server {
	addr := fmt.Sprintf(":%d", webConfig.GetListen())

	server := &http.Server{
		Addr:         addr,
		Handler:      engine,
		ReadTimeout:  webConfig.GetReadTimeout(),
		WriteTimeout: webConfig.GetWriteTimeout(),
		IdleTimeout:  webConfig.GetIdleTimeout(),
	}

	if webConfig.GetUseSSL() {
		server.TLSConfig = &tls.Config{
			MinVersion: tls.VersionTLS12,
		}
		_ = http2.ConfigureServer(server, &http2.Server{})
	}

	return server
}

// httpServerLifecycle manages starting and stopping the HTTP server
func httpServerLifecycle(lc fx.Lifecycle, engine *gin.Engine, params HTTPServerParams) {
	webConfig := params.AppConfig.GetWeb()
	server := buildHTTPServer(engine, webConfig)

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			params.WaitGroup.Add(1)

			go func() {
				defer params.WaitGroup.Done()
				params.Log.Info(fmt.Sprintf("Starting HTTP server on %s", server.Addr))

				var err error
				if webConfig.GetUseSSL() {
					err = server.ListenAndServeTLS(webConfig.GetSSLCert(), webConfig.GetSSLKey())
				} else {
					err = server.ListenAndServe()
				}

				if err != nil && err != http.ErrServerClosed {
					params.Log.Failure(fmt.Sprintf("HTTP server error: %v", err))
				}
			}()

			time.Sleep(100 * time.Millisecond)
			params.Log.Success("HTTP server started successfully")
			return nil
		},
		OnStop: func(ctx context.Context) error {
			params.Log.Info("Stopping HTTP server...")

			shutdownCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
			defer cancel()

			if err := server.Shutdown(shutdownCtx); err != nil {
				params.Log.Failure(fmt.Sprintf("HTTP server shutdown error: %v", err))
				return err
			}

			params.Log.Success("HTTP server stopped successfully")
			return nil
		},
	})
}

// HTTPServerModule wires the HTTP server to the Fx container
func HTTPServerModule() fx.Option {
	return fx.Module("http_server",
		fx.Provide(NewHTTPServer),
		fx.Invoke(httpServerLifecycle),
	)
}
