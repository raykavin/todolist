package module

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"
	"sync"

	"go.uber.org/fx"

	"todolist/internal/adapter/log"
	"todolist/internal/config"
	pkgConfig "todolist/pkg/config"
	pkgLog "todolist/pkg/log"
)

// Config related types
type (
	ConfigParams struct {
		fx.In
		ConfigFile  string `name:"config_file"`
		WatchConfig bool   `name:"watch_config"`
	}

	ConfigContainer struct {
		fx.Out
		App             config.ApplicationProvider
		DefaultDatabase config.DatabaseServiceProvider
	}
)

// Context related types
type (
	ContextResult struct {
		fx.Out
		Context context.Context
		Cancel  context.CancelFunc
	}

	WaitGroupResult struct {
		fx.Out
		WaitGroup *sync.WaitGroup
	}
)

// Logger related types
type LoggerParams struct {
	fx.In
	Config config.ApplicationProvider
}

// Context providers
func NewContext() ContextResult {
	ctx, cancel := context.WithCancel(context.Background())
	return ContextResult{
		Context: ctx,
		Cancel:  cancel,
	}
}

func NewWaitGroup() WaitGroupResult {
	return WaitGroupResult{
		WaitGroup: &sync.WaitGroup{},
	}
}

// Config provider
func NewConfig(params ConfigParams) (ConfigContainer, error) {
	loader, err := createConfigLoader(params.ConfigFile, params.WatchConfig)
	if err != nil {
		return ConfigContainer{}, fmt.Errorf("failed to create config loader: %w", err)
	}

	cfg, err := loader.Load()
	if err != nil {
		return ConfigContainer{}, fmt.Errorf("failed to load config: %w", err)
	}

	defaultDB, err := getDefaultDatabase(cfg)
	if err != nil {
		return ConfigContainer{}, err
	}

	return ConfigContainer{
		App:             cfg.GetApplication(),
		DefaultDatabase: defaultDB,
	}, nil
}

// createConfigLoader creates a config loader with the specified options
func createConfigLoader(configFile string, watchConfig bool) (
	*pkgConfig.ConfigLoader[config.Config],
	error,
) {
	opts := pkgConfig.DefaultLoaderOptions[config.Config]()

	// Extract config name and type from file path
	baseName := filepath.Base(configFile)
	opts.ConfigName = strings.TrimPrefix(baseName, ".")
	opts.ConfigType = strings.TrimSuffix(baseName, ".")
	opts.ConfigPaths = append(opts.ConfigPaths, filepath.Dir(configFile))
	opts.WatchConfig = watchConfig

	return pkgConfig.New(opts), nil
}

// getDefaultDatabase retrieves the default database configuration
func getDefaultDatabase(cfg *config.Config) (config.DatabaseServiceProvider, error) {
	appName := cfg.GetApplication().GetName()
	db, err := cfg.GetDatabase(appName)
	if err != nil {
		return nil, fmt.Errorf("no database found for default application name: %s", appName)
	}
	return db, nil
}

// Logger provider
func NewLogger(params LoggerParams) (pkgLog.ExtendedLog, error) {
	const (
		timeFormat    = "2006-01-02 15:04:05"
		colorEnabled  = true
		callerEnabled = false
	)

	logger, err := pkgLog.NewSmartLog(
		params.Config.GetLoggerLevel(),
		timeFormat,
		colorEnabled,
		callerEnabled,
	)
	if err != nil {
		return nil, fmt.Errorf("error initializing smart logger: %w", err)
	}

	return &log.SmartLogAdapter{Logger: logger}, nil
}

// Lifecycle hooks
func ContextCancelOnStop(lc fx.Lifecycle, cancel context.CancelFunc) {
	lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			cancel()
			return nil
		},
	})
}

// Core returns the fx module with all core dependencies
func Core(configFile string, watchConfig bool) fx.Option {
	return fx.Module("core",
		// Supply configuration parameters
		fx.Supply(
			fx.Annotate(configFile, fx.ResultTags(`name:"config_file"`)),
			fx.Annotate(watchConfig, fx.ResultTags(`name:"watch_config"`)),
		),
		// Provide dependencies
		fx.Provide(
			NewContext,
			NewWaitGroup,
			NewConfig,
			NewLogger,
		),
		// Register lifecycle hooks
		fx.Invoke(ContextCancelOnStop),
	)
}
