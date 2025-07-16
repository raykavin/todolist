package log

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/rs/zerolog/pkgerrors"
)

// Logger wraps zerolog.Logger with enhanced functionality
type Logger struct {
	*zerolog.Logger
	config *Config
}

// Config holds logger configuration
type Config struct {
	Level          string
	DateTimeLayout string
	Colored        bool
	JSONFormat     bool
	UseEmoji       bool
}

// logLevel represents a log level with its properties
type logLevel struct {
	Text  string
	Emoji string
	Color *color.Color
}

// Constants for formatting
const (
	maxMessageSize = 60
	maxFileSize    = 22
	maxLineSize    = 4
	progressBarLen = 20
)

// Log levels definitions
var logLevels = map[string]logLevel{
	zerolog.LevelTraceValue: {
		Text:  "TRAC",
		Emoji: "◇",
		Color: color.New(color.FgHiBlack, color.Bold),
	},
	zerolog.LevelDebugValue: {
		Text:  "DEBG",
		Emoji: "◈",
		Color: color.New(color.FgHiBlue, color.Bold),
	},
	zerolog.LevelInfoValue: {
		Text:  "INFO",
		Emoji: "◉",
		Color: color.New(color.FgHiGreen, color.Bold),
	},
	zerolog.LevelWarnValue: {
		Text:  "WARN",
		Emoji: "◎",
		Color: color.New(color.FgHiYellow, color.Bold),
	},
	zerolog.LevelErrorValue: {
		Text:  "ERRO",
		Emoji: "✖",
		Color: color.New(color.FgHiRed, color.Bold),
	},
	zerolog.LevelFatalValue: {
		Text:  "FATL",
		Emoji: "☠",
		Color: color.New(color.FgHiRed, color.Bold),
	},
	zerolog.LevelPanicValue: {
		Text:  "PANC",
		Emoji: "☠",
		Color: color.New(color.FgWhite, color.BgRed, color.Bold, color.BlinkSlow),
	},
}

// Color scheme for components
var (
	timestampColor = color.New(color.FgHiCyan, color.Italic)
	callerColor    = color.New(color.FgHiMagenta)
	messageColor   = color.New(color.FgWhite)
	fieldKeyColor  = color.New(color.FgHiYellow)
	fieldValColor  = color.New(color.FgCyan)
)

// New creates a new Logger instance
func New(config *Config) (*Logger, error) {
	if config == nil {
		config = &Config{
			Level:          "info",
			DateTimeLayout: time.RFC3339,
			Colored:        true,
			JSONFormat:     false,
			UseEmoji:       false,
		}
	}

	// Setup error stack marshaler
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack

	// Parse log level
	logMode, err := zerolog.ParseLevel(config.Level)
	if err != nil {
		return nil, fmt.Errorf("invalid log level: %w", err)
	}
	zerolog.SetGlobalLevel(logMode)

	// Create logger based on format
	var logger zerolog.Logger
	if config.JSONFormat {
		logger = createJSONLogger(config)
	} else {
		logger = createConsoleLogger(config)
	}

	// Add caller information
	logger = logger.With().CallerWithSkipFrameCount(3).Logger()

	return &Logger{
		Logger: &logger,
		config: config,
	}, nil
}

// createJSONLogger creates a JSON formatted logger
func createJSONLogger(config *Config) zerolog.Logger {
	return log.Output(zerolog.ConsoleWriter{
		Out:           os.Stdout,
		NoColor:       !config.Colored,
		TimeFormat:    config.DateTimeLayout,
		PartsOrder:    []string{"time", "level", "caller", "message"},
		FieldsExclude: []string{"caller"},
	})
}

// createConsoleLogger creates a console formatted logger
func createConsoleLogger(config *Config) zerolog.Logger {
	output := zerolog.ConsoleWriter{
		Out:        os.Stdout,
		NoColor:    !config.Colored,
		TimeFormat: config.DateTimeLayout,
		PartsOrder: []string{"time", "level", "caller", "message"},
	}

	if config.Colored {
		// Create formatter with config
		formatter := &consoleFormatter{config: config}

		output.FormatMessage = formatter.formatMessage
		output.FormatCaller = formatter.formatCaller
		output.FormatLevel = formatter.formatLevel
		output.FormatTimestamp = formatter.formatTimestamp
		output.FormatFieldName = formatter.formatFieldName
		output.FormatFieldValue = formatter.formatFieldValue
	}

	return log.Output(output)
}

// consoleFormatter handles console output formatting
type consoleFormatter struct {
	config *Config
}

// formatLevel formats the log level
func (f *consoleFormatter) formatLevel(i any) string {
	levelStr, ok := i.(string)
	if !ok {
		return f.formatUnknownLevel()
	}

	level, exists := logLevels[levelStr]
	if !exists {
		return f.formatUnknownLevel()
	}

	emoji := ""
	if f.config.UseEmoji {
		emoji = level.Emoji + " "
	}

	return level.Color.Sprintf(" %s%s ", emoji, level.Text)
}

// formatUnknownLevel formats unknown levels
func (f *consoleFormatter) formatUnknownLevel() string {
	emoji := ""
	if f.config.UseEmoji {
		emoji = "❓ "
	}
	return color.New(color.FgHiWhite).Sprintf(" %sUNKN ", emoji)
}

// formatMessage formats the log message
func (f *consoleFormatter) formatMessage(i any) string {
	msg, ok := i.(string)
	if !ok || len(msg) == 0 {
		return messageColor.Sprint("│ (empty message)")
	}

	// Handle multiline messages
	if strings.Contains(msg, "\n") {
		return f.formatMultilineMessage(msg)
	}

	// Truncate and pad single line messages
	if len(msg) > maxMessageSize {
		msg = msg[:maxMessageSize]
	} else {
		msg = fmt.Sprintf("%-*s", maxMessageSize, msg)
	}

	return messageColor.Sprintf("│ %s", msg)
}

// formatMultilineMessage formats messages with multiple lines
func (f *consoleFormatter) formatMultilineMessage(msg string) string {
	lines := strings.Split(msg, "\n")
	formatted := make([]string, len(lines))

	for i, line := range lines {
		formatted[i] = messageColor.Sprintf("│ %s", line)
	}

	return strings.Join(formatted, "\n")
}

// formatCaller formats the caller information
func (f *consoleFormatter) formatCaller(i any) string {
	fname, ok := i.(string)
	if !ok || len(fname) == 0 {
		return ""
	}

	caller := filepath.Base(fname)
	parts := strings.Split(caller, ":")
	if len(parts) != 2 {
		return callerColor.Sprintf("┤ %s ├", caller)
	}

	file := f.formatFileName(parts[0])
	line := f.formatLineNumber(parts[1])

	return callerColor.Sprintf("┤ %s:%s ├", file, line)
}

// formatFileName formats the file name
func (f *consoleFormatter) formatFileName(name string) string {
	file := strings.TrimSuffix(name, ".go")
	if len(file) > maxFileSize {
		return file[:maxFileSize]
	}
	return fmt.Sprintf("%-*s", maxFileSize, file)
}

// formatLineNumber formats the line number
func (f *consoleFormatter) formatLineNumber(line string) string {
	if len(line) > maxLineSize {
		return line[len(line)-maxLineSize:]
	}
	return fmt.Sprintf("%0*s", maxLineSize, line)
}

// formatTimestamp formats the timestamp
func (f *consoleFormatter) formatTimestamp(i any) string {
	strTime, ok := i.(string)
	if !ok {
		return timestampColor.Sprintf("[ %v ]", i)
	}

	ts, err := time.ParseInLocation(time.RFC3339, strTime, time.Local)
	if err != nil {
		return timestampColor.Sprintf("[ %s ]", strTime)
	}

	formatted := ts.In(time.Local).Format(f.config.DateTimeLayout)
	return timestampColor.Sprintf("[ %s ]", formatted)
}

// formatFieldName formats field names
func (f *consoleFormatter) formatFieldName(i any) string {
	name, ok := i.(string)
	if !ok {
		return fmt.Sprintf("%v", i)
	}
	return fieldKeyColor.Sprint(name)
}

// formatFieldValue formats field values
func (f *consoleFormatter) formatFieldValue(i any) string {
	switch v := i.(type) {
	case string:
		// Only quote strings that contain special characters
		if strings.ContainsAny(v, " \t\n\r\"'") {
			return "=" + fieldValColor.Sprintf("%q", v)
		}
		return "=" + fieldValColor.Sprint(v)
	case int, int8, int16, int32, int64,
		uint, uint8, uint16, uint32, uint64:
		return fieldValColor.Sprintf("=%d", v)
	case float32, float64:
		return fieldValColor.Sprintf("=%.2f", v)
	case bool:
		if v {
			return "=" + color.HiGreenString("true")
		}
		return "=" + color.HiRedString("false")
	case nil:
		return "=" + color.HiBlackString("null")
	default:
		return fieldValColor.Sprintf("=%v", v)
	}
}

// Enhanced logging methods

// Success logs a success message
func (l *Logger) Success(msg string) {
	if l.config.UseEmoji {
		msg = "✅ " + msg
	}
	l.Info().Msg(msg)
}

// Failure logs a failure message
func (l *Logger) Failure(msg string) {
	if l.config.UseEmoji {
		msg = "❌ " + msg
	}
	l.Error().Msg(msg)
}

// Progress logs a progress update
func (l *Logger) Progress(msg string, current, total int) {
	percentage := float64(current) / float64(total) * 100
	progressBar := l.createProgressBar(int(percentage))

	l.Info().
		Str("progress", progressBar).
		Float64("percent", percentage).
		Int("current", current).
		Int("total", total).
		Msg(msg)
}

// Benchmark logs a benchmark result
func (l *Logger) Benchmark(name string, duration time.Duration) {
	msg := "Benchmark:"

	if l.config.UseEmoji {
		emoji := l.getDurationEmoji(duration)
		msg = fmt.Sprintf("%s %s", emoji, msg)
	}

	l.Debug().
		Str("duration", duration.String()).
		Msgf("%s %s", msg, name)
}

// API logs an API request
func (l *Logger) API(method, path string, statusCode int, duration time.Duration) {
	level := l.getStatusLevel(statusCode)
	msg := "API Request"

	if l.config.UseEmoji {
		emoji := l.getStatusEmoji(statusCode)
		msg = fmt.Sprintf("%s %s", emoji, msg)
	}

	l.WithLevel(level).
		Str("method", method).
		Str("path", path).
		Int("status_code", statusCode).
		Str("duration", duration.Round(time.Millisecond).String()).
		Msg(msg)
}

// WithContext creates a new logger with additional context
func (l *Logger) WithContext(ctx map[string]any) *Logger {
	event := l.With()
	for k, v := range ctx {
		event = event.Interface(k, v)
	}
	logger := event.Logger()
	return &Logger{
		Logger: &logger,
		config: l.config,
	}
}

// Helper methods

// createProgressBar creates a visual progress bar
func (l *Logger) createProgressBar(percentage int) string {
	filled := percentage * progressBarLen / 100

	var bar strings.Builder
	bar.WriteByte('[')

	for i := 0; i < progressBarLen; i++ {
		if i < filled {
			bar.WriteRune('█')
		} else {
			bar.WriteRune('░')
		}
	}

	bar.WriteString(fmt.Sprintf("] %d%%", percentage))
	return bar.String()
}

// getDurationEmoji returns an emoji based on duration
func (l *Logger) getDurationEmoji(duration time.Duration) string {
	switch {
	case duration < time.Millisecond:
		return "⚡" // Very fast
	case duration < 10*time.Millisecond:
		return "🚀" // Fast
	case duration < 100*time.Millisecond:
		return "🏃" // Medium
	case duration < time.Second:
		return "🚶" // Slow
	default:
		return "🐌" // Very slow
	}
}

// getStatusLevel returns the appropriate log level for HTTP status code
func (l *Logger) getStatusLevel(statusCode int) zerolog.Level {
	switch {
	case statusCode >= 200 && statusCode < 300:
		return zerolog.InfoLevel
	case statusCode >= 300 && statusCode < 400:
		return zerolog.InfoLevel
	case statusCode >= 400 && statusCode < 500:
		return zerolog.WarnLevel
	case statusCode >= 500:
		return zerolog.ErrorLevel
	default:
		return zerolog.InfoLevel
	}
}

// getStatusEmoji returns an emoji for HTTP status code
func (l *Logger) getStatusEmoji(statusCode int) string {
	switch {
	case statusCode >= 200 && statusCode < 300:
		return "✅ "
	case statusCode >= 300 && statusCode < 400:
		return "🔄 "
	case statusCode >= 400 && statusCode < 500:
		return "⚠️ "
	case statusCode >= 500:
		return "❌ "
	default:
		return "❓ "
	}
}
