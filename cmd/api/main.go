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
	"todolist/internal/fx/module"
	"todolist/pkg/log"
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

const shutdownTimeout = 25 * time.Second

func main() {
	flag.Parse()

	app := fx.New(
		// Configure Fx logger
		fx.WithLogger(configureFxLogger),

		// Dependency injection modules
		module.Core(*configFile, *watchConfig), // Core: context, config, logger, wait group
		module.Repositories(),                  // Repositories: database repositories
		module.ApplicationServices(),           // Services: application services
		module.DomainServices(),                // Services: domain services
		module.UseCases(),                      // UseCases: business logic
		module.HTTPHandlers(),                  // HTTPHandler: HTTP server and routes

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
	terminal.PrintBanner(config.GetName())
	terminal.PrintText(config.GetDescription())
	terminal.PrintText("S.O.G.E - Sistemas Operacionais, Gerenciais e Estratégicos")
	terminal.PrintText(fmt.Sprintf("Copyright (c) %d I R Tecnologia, Todos os direitos reservados!",
		time.Now().Year()))
	terminal.PrintHeader(fmt.Sprintf("Version: %s", config.GetVersion()))
}

// runApplication registers startup and shutdown hooks
func runApplication(
	lc fx.Lifecycle,
	ctx context.Context,
	cancel context.CancelFunc,
	log log.ExtendedLog,
	wg *sync.WaitGroup,
) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			log.Success("PON Watcher application started")
			return nil
		},
		OnStop: func(ctx context.Context) error {
			return gracefulShutdown(cancel, log, wg)
		},
	})
}

// gracefulShutdown cancels the main context and waits for goroutines to finish
func gracefulShutdown(
	cancel context.CancelFunc,
	log log.ExtendedLog,
	wg *sync.WaitGroup,
) error {
	log.Info("Shutting down PON Watcher application...")

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
		log.Success("All goroutines finished successfully")
	case <-time.After(shutdownTimeout):
		log.Failure("Timeout waiting for goroutines to finish")
	}

	return nil
}

// handleAppLifecycle sets up OS signal handling for graceful shutdown
func handleAppLifecycle(
	lc fx.Lifecycle,
	shutdowner fx.Shutdowner,
	log log.ExtendedLog,
) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go signalHandler(shutdowner, log)
			return nil
		},
		OnStop: func(ctx context.Context) error {
			log.Success("Application stopped successfully")
			return nil
		},
	})
}

// signalHandler listens for OS signals to trigger application shutdown
func signalHandler(shutdowner fx.Shutdowner, log log.ExtendedLog) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	log.Info("PON Watcher is running. Press Ctrl+C to stop...")
	<-quit

	log.Info("Shutdown signal received")
	log.Warn("Closing connections and cleaning up, please wait...")

	if err := shutdowner.Shutdown(); err != nil {
		log.Failure(fmt.Sprintf("Error during shutdown: %v", err))
	}
}
