package box

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"time"
)

type Env struct {
	Config *Config

	Logger *zap.Logger

	PGXPool *pgxpool.Pool
}

func NewEnv(ctx context.Context) (*Env, error) {
	config, err := provideConfig()
	if err != nil {
		return nil, fmt.Errorf("unable to load config | %w", err)
	}

	logger := provideLogger(config.Debug)

	pgxPool, err := providePGXPool(ctx, config.PostgresConfig.ConnectionString)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to postgres | %w", err)
	}

	return &Env{
		Config:  config,
		Logger:  logger,
		PGXPool: pgxPool,
	}, nil
}

func provideConfig() (*Config, error) {
	cfg, err := FromEnv()
	if err != nil {
		return nil, fmt.Errorf("cannot parse config from environment | %w", err)
	}

	return cfg, nil
}

func provideLogger(debug bool) *zap.Logger {
	encoderCfg := zapcore.EncoderConfig{
		MessageKey:  "_m",
		NameKey:     "logger",
		LevelKey:    "_l",
		EncodeLevel: zapcore.LowercaseLevelEncoder,
		TimeKey:     "_t",
		EncodeTime:  zapcore.ISO8601TimeEncoder,
	}

	var level zapcore.Level
	if debug {
		level = zap.DebugLevel
	} else {
		level = zap.InfoLevel
	}

	return zap.New(zapcore.NewCore(zapcore.NewJSONEncoder(encoderCfg), os.Stdout, level))
}

func providePGXPool(ctx context.Context, connConfig string) (*pgxpool.Pool, error) {
	ctx, cancel := context.WithTimeout(ctx, 20*time.Second)
	defer cancel()

	databaseConfig, err := pgxpool.ParseConfig(connConfig)
	if err != nil {
		return nil, fmt.Errorf("error parsing config | %w", err)
	}

	pool, err := pgxpool.NewWithConfig(ctx, databaseConfig)
	if err != nil {
		return nil, fmt.Errorf("error creating pool | %w", err)
	}

	if err = pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("error pinging pool | %w", err)
	}

	return pool, nil
}
