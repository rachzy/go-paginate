package paginate

import (
	"context"
	"log/slog"
	"os"
	"sync"
)

var loggerDefaultOpts = &slog.HandlerOptions{
	Level: slog.LevelInfo,
}

var logLevels = map[string]slog.Level{
	"INFO":  slog.LevelInfo,
	"DEBUG": slog.LevelDebug,
	"WARN":  slog.LevelWarn,
	"ERROR": slog.LevelError,
}

type Log struct {
	Level   slog.Level
	Message string
	Args    []any
}

var loggerOnce sync.Once
var waitGroup sync.WaitGroup

// logsChan is a channel for safely logging messages according to user's configuration
// useful for cases where the logger might not yet be initialized
var logsChan = make(chan Log)

// initLoggerOnce initializes the logger based on the global configuration
func initLoggerOnce(ctx context.Context) {
	loggerOnce.Do(func() {
		if globalConfig.logger == nil {
			opts := getDefaultOpts()
			globalConfig.logger = slog.New(slog.NewTextHandler(os.Stdout, opts))
		}

		waitGroup.Add(1)
		go func() {
			defer waitGroup.Done()
			select {
			case <-ctx.Done():
				return
			case log := <-logsChan:
				globalConfig.logger.Log(ctx, log.Level, log.Message, log.Args...)
			}
		}()

		waitGroup.Wait()
	})
}

func getDefaultOpts() *slog.HandlerOptions {
	// Load log level from environment variable
	if logLevelStr := os.Getenv("GO_PAGINATE_LOG_LEVEL"); logLevelStr != "" {
		if logLevel, ok := logLevels[logLevelStr]; ok {
			loggerDefaultOpts.Level = logLevel
		} else {
			SafeLog(globalConfig.ctx, slog.LevelWarn, "Invalid GO_PAGINATE_LOG_LEVEL value, using default",
				"value", logLevelStr,
				"default", loggerDefaultOpts.Level,
			)
		}
	}

	return loggerDefaultOpts
}

func SafeLog(ctx context.Context, level slog.Level, message string, args ...any) {
	if globalConfig.logger == nil {
		go func() {
			select {
			case <-ctx.Done():
				return
			case logsChan <- Log{
				Level:   level,
				Message: message,
				Args:    args,
			}:
			}
		}()
		return
	}

	globalConfig.logger.Log(ctx, level, message, args...)
}
