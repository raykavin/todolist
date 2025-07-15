package di

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"
	"sync"

	"github.com/kataras/iris/v12/x/errors"
	"go.uber.org/fx"

	"todolist/internal/config"
	pkgConfig "todolist/pkg/config"
)

// ConfigParams defines the dependencies to load configuration
type ConfigParams struct {
	fx.In
	ConfigFile  string `name:"config_file"`
	WatchConfig bool   `name:"watch_config"`
}

// ConfigContainer groups all core configurations implementations provide from Fx
type ConfigContainer struct {
	fx.Out
	AppConfig       config.ApplicationProvider
	DefaultDatabase config.DatabaseServiceProvider
}

// Main application context
type ContextResult struct {
	fx.Out
	Context context.Context
	Cancel  context.CancelFunc
}

// Main application wait group
type WaitGroupResult struct {
	fx.Out
	WaitGroup *sync.WaitGroup
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
		AppConfig:       cfg.GetApplication(),
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
	ext := filepath.Ext(baseName)
	filename := strings.TrimSuffix(baseName, ext)

	opts.ConfigName = filename
	opts.ConfigType = strings.TrimPrefix(ext, ".")
	opts.ConfigPaths = append(opts.ConfigPaths, filepath.Dir(configFile))
	opts.WatchConfig = watchConfig

	return pkgConfig.New(opts), nil
}

// getDefaultDatabase retrieves the default database configuration
func getDefaultDatabase(cfg *config.Config) (config.DatabaseServiceProvider, error) {
	db, err := cfg.GetDatabase("default")
	if err != nil {
		return nil, errors.New("no database found for default application name: default")
	}
	return db, nil
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

// CoreModule returns the fx module with all core dependencies
func CoreModule(configFile string, watchConfig bool) fx.Option {
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
		),
		// Register lifecycle hooks
		fx.Invoke(ContextCancelOnStop),
	)
}
