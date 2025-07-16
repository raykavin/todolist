package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
	"todolist/internal/config"
	"todolist/internal/di"
	"todolist/pkg/logger"
	"todolist/pkg/terminal"

	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
)

// Command line flags
var (
	configFile  = flag.String("config", "config.yaml", "Path to the configuration file")
	watchConfig = flag.Bool("watch-config", false, "Watch configuration file for changes")
	fxDebug     = flag.Bool("fx-debug", false, "Enable Fx dependency injection debug logs")
)

const shutdownTimeout = 5 * time.Second

// @title           Todo List API
// @version         1.0
// @description     A simple Todo List application example

// @contact.name   SOGE Fibralink
// @contact.url    https://fibralink.net.br
// @contact.email  soge@fibralink.net.br

// @host      localhost:3000
// @BasePath  /api/v1

func main() {
	flag.Parse()

	app := fx.New(
		// Configure Fx logger
		fx.WithLogger(configureFxLogger),

		// Dependency injection modules
		di.CoreModule(*configFile, *watchConfig), // Core: context, config, wait group
		di.LoggerModule(),                        // Logger: logger infrastructure
		di.DatabasesModule(),                     // Databases: database infrastructures
		di.RepositoriesModule(),                  // Repositories: database repositories
		di.ApplicationServicesModule(),           // Services: application services
		di.DomainServicesModule(),                // Services: complex domain services business logic
		di.UseCasesModule(),                      // UseCases: specifics business logic
		di.HTTPHandlersModule(),                  // HTTPHandler: HTTP handlers
		di.HTTPServerModule(),                    // HTTPServer: HTTP server setup

		// Application lifecycle hooks
		fx.Invoke(displayAppInfo),
		fx.Invoke(runApplication),
		fx.Invoke(handleAppLifecycle),
	)

	app.Run()
}

// configureFxLogger returns the Fx logger based on the debug flag
func configureFxLogger() fxevent.Logger {
	if !*fxDebug {
		return fxevent.NopLogger
	}
	return &fxevent.ConsoleLogger{W: os.Stderr}
}

// displayAppInfo prints the application banner and basic info
func displayAppInfo(config config.ApplicationProvider) {
	displayText := "S.O.G.E - Sistemas Operacionais, Gerenciais e Estratégicos"
	displayText2 := fmt.Sprintf("Copyright (c) %d I R Tecnologia, Todos os direitos reservados!", time.Now().Year())
	displayText3 := fmt.Sprintf("Version: %s", config.GetVersion())

	terminal.PrintBanner(config.GetName())
	terminal.PrintText(config.GetDescription())
	terminal.PrintText(displayText)
	terminal.PrintText(displayText2)
	terminal.PrintHeader(displayText3)
}

// runApplication registers startup and shutdown hooks
func runApplication(
	lc fx.Lifecycle,
	ctx context.Context,
	cancel context.CancelFunc,
	logger logger.ExtendedLog,
	wg *sync.WaitGroup,
) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			logger.Success("Application started")
			return nil
		},
		OnStop: func(ctx context.Context) error {
			return gracefulShutdown(cancel, logger, wg)
		},
	})
}

// gracefulShutdown cancels the main context and waits for goroutines to finish
func gracefulShutdown(
	cancel context.CancelFunc,
	logger logger.ExtendedLog,
	wg *sync.WaitGroup,
) error {
	logger.Info("Shutting down application...")

	// Cancel the main context
	cancel()

	// Wait for all goroutines to finish with timeout
	done := make(chan struct{})
	go func() {
		defer close(done)
		wg.Wait()
	}()

	select {
	case <-done:
		logger.Success("All goroutines finished successfully")
	case <-time.After(shutdownTimeout):
		logger.Failure("Timeout waiting for goroutines to finish")
	}

	return nil
}

// handleAppLifecycle sets up OS signal handling for graceful shutdown
func handleAppLifecycle(
	lc fx.Lifecycle,
	shutdowner fx.Shutdowner,
	logger logger.ExtendedLog,
) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go signalHandler(shutdowner, logger)
			return nil
		},
		OnStop: func(ctx context.Context) error {
			logger.Success("Application stopped successfully")
			return nil
		},
	})
}

// signalHandler listens for OS signals to trigger application shutdown
func signalHandler(shutdowner fx.Shutdowner, logger logger.ExtendedLog) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	logger.Info("Application is running. Press Ctrl+C to stop...")
	<-quit

	logger.Info("Shutdown signal received")
	logger.Warn("Closing connections and cleaning up, please wait...")

	if err := shutdowner.Shutdown(); err != nil {
		logger.Failure(fmt.Sprintf("Error during shutdown: %v", err))
	}
}
