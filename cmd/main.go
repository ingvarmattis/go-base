package main

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"go-base/src/box"
	"go.uber.org/zap"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

func main() {
	ctx := context.Background()

	envBox, err := box.NewEnv(ctx)
	if err != nil {
		panic(err)
	}

	gracefullShutdown(envBox.Logger, envBox.PGXPool)
}

func gracefullShutdown(logger *zap.Logger, pgxPool *pgxpool.Pool) {
	quit := make(chan os.Signal, 1)
	signal.Notify(
		quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGKILL,
	)
	<-quit

	logger.Info("shutting down service...")

	var wg sync.WaitGroup
	closeFuncs := []func(){
		func() {
			defer wg.Done()

			logger.Info("pgx pool closing...")
			pgxPool.Close()
			logger.Info("pgx pool has been closed")
		},
	}
	wg.Add(len(closeFuncs))

	for _, f := range closeFuncs {
		go f()
	}

	wg.Wait()

	logger.Info("application has been gracefully shutdown")
}
