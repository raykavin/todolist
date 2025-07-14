package module

import "go.uber.org/fx"

func HTTPHandlers() fx.Option {
	return fx.Module("handlers")
}
