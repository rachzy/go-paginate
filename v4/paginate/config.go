package paginate

import (
	"context"
	"log/slog"
	"os"
	"strconv"
)

// GlobalConfig holds the global configuration for go-paginate
type GlobalConfig struct {
	ctx    context.Context
	logger *slog.Logger

	DefaultLimit int
	MaxLimit     int
	DebugMode    bool
}

// globalConfig is the singleton instance
var globalConfig = &GlobalConfig{
	ctx: context.Background(),

	DefaultLimit: 10,  // default value
	MaxLimit:     100, // default value
	DebugMode:    false,
}

func init() {
	loadFromEnv()
}

// loadFromEnv loads configuration from environment variables
func loadFromEnv() {
	// Load GO_PAGINATE_DEBUG
	if debugStr := os.Getenv("GO_PAGINATE_DEBUG"); debugStr != "" {
		if debug, err := strconv.ParseBool(debugStr); err == nil {
			globalConfig.DebugMode = debug
			SafeLog(globalConfig.ctx, slog.LevelInfo, "Debug mode loaded from environment",
				"GO_PAGINATE_DEBUG", debug)
		} else {
			SafeLog(globalConfig.ctx, slog.LevelWarn, "Invalid GO_PAGINATE_DEBUG value, using default",
				"value", debugStr,
				"error", err,
				"default", globalConfig.DebugMode)
		}
	}

	// Load GO_PAGINATE_DEFAULT_LIMIT
	if limitStr := os.Getenv("GO_PAGINATE_DEFAULT_LIMIT"); limitStr != "" {
		if limit, err := strconv.Atoi(limitStr); err == nil && limit > 0 {
			globalConfig.DefaultLimit = limit
			SafeLog(globalConfig.ctx, slog.LevelInfo, "Default limit loaded from environment",
				"GO_PAGINATE_DEFAULT_LIMIT", limit)
		} else {
			SafeLog(globalConfig.ctx, slog.LevelWarn, "Invalid GO_PAGINATE_DEFAULT_LIMIT value, using default",
				"value", limitStr,
				"error", err,
				"default", globalConfig.DefaultLimit)
		}
	}

	// Load GO_PAGINATE_MAX_LIMIT
	if maxLimitStr := os.Getenv("GO_PAGINATE_MAX_LIMIT"); maxLimitStr != "" {
		if maxLimit, err := strconv.Atoi(maxLimitStr); err == nil && maxLimit > 0 {
			globalConfig.MaxLimit = maxLimit
			SafeLog(globalConfig.ctx, slog.LevelInfo, "Max limit loaded from environment",
				"GO_PAGINATE_MAX_LIMIT", maxLimit)
		} else {
			SafeLog(globalConfig.ctx, slog.LevelWarn, "Invalid GO_PAGINATE_MAX_LIMIT value, using default",
				"value", maxLimitStr,
				"error", err,
				"default", globalConfig.MaxLimit)
		}
	}

	SafeLog(globalConfig.ctx, slog.LevelInfo, "Go-paginate configuration initialized",
		"defaultLimit", globalConfig.DefaultLimit,
		"maxLimit", globalConfig.MaxLimit,
		"debugMode", globalConfig.DebugMode)
}

func Init(ctx context.Context, config *GlobalConfig) {
	globalConfig.ctx = ctx
	globalConfig.MaxLimit = config.MaxLimit
	globalConfig.DefaultLimit = config.DefaultLimit
	globalConfig.DebugMode = config.DebugMode

	initLoggerOnce(globalConfig.ctx)
}

func InitWithLogger(ctx context.Context, config *GlobalConfig, logger *slog.Logger) {
	globalConfig.logger = logger

	Init(ctx, config)
}

// SetDefaultLimit sets the global default limit
func SetDefaultLimit(limit int) {
	logger := getLogger("go-paginate-config")

	if limit <= 0 {
		logger.Error("Invalid default limit value, must be greater than 0",
			"attempted_value", limit,
			"current_value", globalConfig.DefaultLimit)
		return
	}

	oldValue := globalConfig.DefaultLimit
	globalConfig.DefaultLimit = limit

	logger.Info("Default limit updated",
		"old_value", oldValue,
		"new_value", limit)
}

// SetMaxLimit sets the global maximum limit
func SetMaxLimit(maxLimit int) {
	logger := getLogger("go-paginate-config")

	if maxLimit <= 0 {
		logger.Error("Invalid max limit value, must be greater than 0",
			"attempted_value", maxLimit,
			"current_value", globalConfig.MaxLimit)
		return
	}

	oldValue := globalConfig.MaxLimit
	globalConfig.MaxLimit = maxLimit

	logger.Info("Max limit updated",
		"old_value", oldValue,
		"new_value", maxLimit)
}

// SetDebugMode sets the global debug mode
func SetDebugMode(debug bool) {
	logger := getLogger("go-paginate-config")

	oldValue := globalConfig.DebugMode
	globalConfig.DebugMode = debug

	logger.Info("Debug mode updated",
		"old_value", oldValue,
		"new_value", debug)
}

// GetDefaultLimit returns the global default limit
func GetDefaultLimit() int {
	return globalConfig.DefaultLimit
}

// GetMaxLimit returns the global maximum limit
func GetMaxLimit() int {
	return globalConfig.MaxLimit
}

// IsDebugMode returns the global debug mode status
func IsDebugMode() bool {
	return globalConfig.DebugMode
}

// SetLogger sets a custom logger for the configuration
func SetLogger(logger *slog.Logger) {
	globalConfig.logger = logger
	initLoggerOnce(globalConfig.ctx)
}

// logSQL logs SQL queries when debug mode is enabled
func logSQL(operation, query string, args []any) {
	if globalConfig.DebugMode {
		logger := getLogger("go-paginate-sql")
		logger.Info("Generated SQL query",
			"operation", operation,
			"query", query,
			"args", args,
			"args_count", len(args))
	}
}

func getLogger(component string) *slog.Logger {
	initLoggerOnce(globalConfig.ctx)
	return globalConfig.logger.With("component", component)
}
