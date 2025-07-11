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

	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
)

// Command line flags
var (
	configFile  = flag.String("config", "config.yaml", "Configuration file path")
	watchConfig = flag.Bool("watch-config", false, "Watch configuration file for changes")
	fxDebug     = flag.Bool("fx-debug", false, "Enable/disable fx dependency injector logs")
)

const shutdownTimeout = 25 * time.Second

func main() {
	flag.Parse()

	app := fx.New(
		// Configure FX Logger
		fx.WithLogger(configureFxLogger),

		// DI Modules
		module.Core(*configFile, *watchConfig), // Core: context, config, logger, wait group
		module.Repository(),                    // Repository: database repositories
		module.Services(),                      // Service: domain services
		module.UseCase(),                       // UseCase: business logic
		module.HTTPHandler(),                   // HTTPHandler: HTTP server and handlers

		// Application lifecycle
		fx.Invoke(displayAppInfo),
		fx.Invoke(runApplication),
		fx.Invoke(handleAppLifecycle),
	)

	app.Run()
}

// configureFxLogger sets up the FX dependency injection logger based on debug flag
func configureFxLogger() fxevent.Logger {
	if !*fxDebug {
		return fxevent.NopLogger
	}
	return &fxevent.ConsoleLogger{W: os.Stderr}
}

// displayAppInfo shows application banner and information
func displayAppInfo(config config.ApplicationProvider) {
	banner.PrintBanner(config.GetName())
	banner.PrintText(config.GetDescription())
	banner.PrintText("S.O.G.E - Sistemas Operacionais, Gerências e Estratégicos")
	banner.PrintText(fmt.Sprintf("Copyright (c) %d I R Tecnologia, All Rights Reserved!", time.Now().Year()))
	banner.PrintHeader(fmt.Sprintf("Version: %s", config.GetVersion()))
}

// runApplication manages the core application lifecycle hooks
func runApplication(
	lc fx.Lifecycle,
	ctx context.Context,
	cancel context.CancelFunc,
	log log.Smart,
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

// gracefulShutdown handles the graceful shutdown process
func gracefulShutdown(cancel context.CancelFunc, log log.Smart, wg *sync.WaitGroup) error {
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
		log.Success("All goroutines finished")
	case <-time.After(shutdownTimeout):
		log.Failure("Timeout waiting for goroutines to finish")
	}

	return nil
}

// handleAppLifecycle manages signal handling and application lifecycle
func handleAppLifecycle(lc fx.Lifecycle, shutdowner fx.Shutdowner, log log.Smart) {
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

// signalHandler handles OS signals for graceful shutdown
func signalHandler(shutdowner fx.Shutdowner, log log.Smart) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	log.Info("PON Watcher is running. Press Ctrl+C to stop...")
	<-quit

	log.Info("Shutdown signal received...")
	log.Warn("Please wait, closing connections and disconnecting...")

	if err := shutdowner.Shutdown(); err != nil {
		log.Failure(fmt.Sprintf("Error during shutdown: %v", err))
	}
}
