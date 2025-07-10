package log

import "time"

type Log interface {
	// Context methods - returns a logger based off the root logger and decorates it with the given context and arguments.
	WithField(key string, value any) Log
	WithFields(fields map[string]any) Log
	WithError(err error) Log

	// Standard log functions
	Print(args ...any)
	Debug(args ...any)
	Info(args ...any)
	Warn(args ...any)
	Error(args ...any)
	Fatal(args ...any)
	Panic(args ...any)

	// Formatted log functions
	Printf(format string, args ...any)
	Debugf(format string, args ...any)
	Infof(format string, args ...any)
	Warnf(format string, args ...any)
	Errorf(format string, args ...any)
	Fatalf(format string, args ...any)
	Panicf(format string, args ...any)
}

// ExtendedLog interface for smart logging features
type ExtendedLog interface {
	Log // Base interface

	// Enhanced logging methods
	Success(msg string)
	Failure(msg string)
	Progress(msg string, current, total int)
	Benchmark(name string, duration time.Duration)
	API(method, path string, statusCode int, duration time.Duration)

	// Context with map support
	WithContext(ctx map[string]any) ExtendedLog
}
