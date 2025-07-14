package web

import (
	"context"
	"errors"
	"fmt"
	netHttp "net/http"

	"github.com/gin-gonic/gin"
)

var (
	ErrInvalidListenAddress = errors.New("invalid listen address")
)

// Ensure it implements http interfaces
// var _ http.HttpServer = (*Engine)(nil)

// Engine wraps Gin engine and implements HttpServer interface
type Engine struct {
	router *gin.Engine
	server *netHttp.Server
	addr   string
}

// NewGin creates a new Gin engine with the specified configuration
func NewGin(host string, listen uint16, debugMode bool) (*Engine, error) {
	if listen == 0 {
		return nil, ErrInvalidListenAddress
	}

	if len(host) == 0 {
		host = ":"
	}

	// Set Gin mode based on logger level
	if debugMode {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()

	// Add default middleware (Recovery)
	router.Use(gin.Recovery())

	return &Engine{
		router: router,
		addr:   fmt.Sprint(),
		server: &netHttp.Server{
			Addr:    fmt.Sprint(host, listen),
			Handler: router,
		},
	}, nil
}

// Listen implements HttpServer interface
func (e *Engine) Listen() error {
	return e.server.ListenAndServe()
}

// Shutdown implements HttpServer interface
func (e *Engine) Shutdown(ctx context.Context) error {
	return e.server.Shutdown(ctx)
}

// Router provides access to the underlying Gin router
func (e *Engine) Router() *gin.Engine {
	return e.router
}

// Server provides access to the underlying HTTP server
func (e *Engine) Server() *netHttp.Server {
	return e.server
}

// SetupRoutes is a helper function to demonstrate route setup
// You can customize this based on your needs
func (e *Engine) SetupRoutes(routes func(router *gin.Engine)) {
	routes(e.router)
}
