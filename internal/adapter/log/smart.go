package log

import (
	"fmt"
	"time"
	smart "todolist/pkg/log"
)

type SmartLogAdapter struct {
	*smart.Logger
}

// Ensure it implements both interfaces
var (
	_ Log         = (*SmartLogAdapter)(nil)
	_ ExtendedLog = (*SmartLogAdapter)(nil)
)

// Print implements Logging.
func (s *SmartLogAdapter) Print(args ...any) {
	s.Logger.Print(args...)
}

// Debug implements Logging.
func (s *SmartLogAdapter) Debug(args ...any) {
	s.Logger.Debug().Msg(fmt.Sprint(args...))
}

// Info implements Logging.
func (s *SmartLogAdapter) Info(args ...any) {
	s.Logger.Info().Msg(fmt.Sprint(args...))
}

// Warn implements Logging.
func (s *SmartLogAdapter) Warn(args ...any) {
	s.Logger.Warn().Msg(fmt.Sprint(args...))
}

// Error implements Logging.
func (s *SmartLogAdapter) Error(args ...any) {
	s.Logger.Error().Msg(fmt.Sprint(args...))
}

// Fatal implements Logging.
func (s *SmartLogAdapter) Fatal(args ...any) {
	s.Logger.Fatal().Msg(fmt.Sprint(args...))
}

// Panic implements Logging.
func (s *SmartLogAdapter) Panic(args ...any) {
	s.Logger.Panic().Msg(fmt.Sprint(args...))
}

// Printf implements Logging.
func (s *SmartLogAdapter) Printf(format string, args ...any) {
	s.Logger.Printf(format, args...)
}

// Debugf implements Logging.
func (s *SmartLogAdapter) Debugf(format string, args ...any) {
	s.Logger.Debug().Msgf(format, args...)
}

// Infof implements Logging.
func (s *SmartLogAdapter) Infof(format string, args ...any) {
	s.Logger.Info().Msgf(format, args...)
}

// Warnf implements Logging.
func (s *SmartLogAdapter) Warnf(format string, args ...any) {
	s.Logger.Warn().Msgf(format, args...)
}

// Errorf implements Logging.
func (s *SmartLogAdapter) Errorf(format string, args ...any) {
	s.Logger.Error().Msgf(format, args...)
}

// Fatalf implements Logging.
func (s *SmartLogAdapter) Fatalf(format string, args ...any) {
	s.Logger.Fatal().Msgf(format, args...)
}

// Panicf implements Logging.
func (s *SmartLogAdapter) Panicf(format string, args ...any) {
	s.Logger.Panic().Msgf(format, args...)
}

// WithError implements Logging.
func (s *SmartLogAdapter) WithError(err error) Log {
	newLogger := s.With().Err(err).Logger()
	return &SmartLogAdapter{&smart.Logger{Logger: &newLogger}}
}

// WithField implements Logging.
func (s *SmartLogAdapter) WithField(key string, value any) Log {
	newLogger := s.With().Interface(key, value).Logger()
	return &SmartLogAdapter{&smart.Logger{Logger: &newLogger}}
}

// WithFields implements Logging.
func (s *SmartLogAdapter) WithFields(fields map[string]any) Log {
	newLogger := s.With().Fields(fields).Logger()
	return &SmartLogAdapter{&smart.Logger{Logger: &newLogger}}
}

// Success implements SmartLogging.
func (s *SmartLogAdapter) Success(msg string) {
	s.Logger.Success(msg)
}

// Failure implements SmartLogging.
func (s *SmartLogAdapter) Failure(msg string) {
	s.Logger.Failure(msg)
}

// Progress implements SmartLogging.
func (s *SmartLogAdapter) Progress(msg string, current, total int) {
	s.Logger.Progress(msg, current, total)
}

// Benchmark implements SmartLogging.
func (s *SmartLogAdapter) Benchmark(name string, duration time.Duration) {
	s.Logger.Benchmark(name, duration)
}

// API implements SmartLogging.
func (s *SmartLogAdapter) API(method, path string, statusCode int, duration time.Duration) {
	s.Logger.API(method, path, statusCode, duration)
}

// WithContext implements SmartLogging.
func (s *SmartLogAdapter) WithContext(ctx map[string]interface{}) ExtendedLog {
	smartLogger := s.Logger.WithContext(ctx)
	return &SmartLogAdapter{smartLogger}
}
