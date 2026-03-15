package traefik

import (
	"github.com/darkweak/storages/core"
	redisstorage "github.com/darkweak/storages/redis"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func initRedisStorage(configuration *Configuration) {
	redis := configuration.GetDefaultCache().GetRedis()
	if redis.Configuration == nil && redis.Path == "" && redis.URL == "" {
		return
	}

	if configuration.GetLogger() == nil {
		var logLevel zapcore.Level
		if configuration.GetLogLevel() == "" {
			logLevel = zapcore.FatalLevel
		} else if err := logLevel.UnmarshalText([]byte(configuration.GetLogLevel())); err != nil {
			logLevel = zapcore.FatalLevel
		}

		cfg := zap.Config{
			Encoding:         "json",
			Level:            zap.NewAtomicLevelAt(logLevel),
			OutputPaths:      []string{"stderr"},
			ErrorOutputPaths: []string{"stderr"},
			EncoderConfig: zapcore.EncoderConfig{
				MessageKey:   "message",
				LevelKey:     "level",
				EncodeLevel:  zapcore.CapitalLevelEncoder,
				TimeKey:      "time",
				EncodeTime:   zapcore.ISO8601TimeEncoder,
				CallerKey:    "caller",
				EncodeCaller: zapcore.ShortCallerEncoder,
			},
		}

		logger, _ := cfg.Build()
		configuration.SetLogger(logger.Sugar())
	}

	if storer, err := redisstorage.Factory(
		core.CacheProvider{
			Configuration: redis.Configuration,
			Path:          redis.Path,
			URL:           redis.URL,
		},
		configuration.GetLogger(),
		configuration.GetDefaultCache().GetStale(),
	); err == nil {
		core.RegisterStorage(storer)
	} else {
		configuration.GetLogger().Warnf("unable to initialize redis storer: %v", err)
	}
}
