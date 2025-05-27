package main

import (
	"context"
	"github.com/go-co-op/gocron/v2"
	"go.uber.org/zap"
	"realtime-bike-go/config"
	"realtime-bike-go/internal/collector"
	"realtime-bike-go/pkg/graceful"
	"realtime-bike-go/pkg/monitoring"
	"realtime-bike-go/pkg/rblogger"
	"time"
)

var (
	configPrefix = "rb-collector"
)

func main() {
	ctx, _ := graceful.WaitForSignalContext(context.Background())

	logger := rblogger.New()
	cfg := config.NewBase(ctx, configPrefix)
	collectorClient, err := collector.New(ctx, logger, cfg.MongoDB)
	if err != nil {
		logger.Error("could not create collector", zap.Error(err))
		return
	}
	defer collectorClient.Close(ctx)

	s, err := gocron.NewScheduler()
	if err != nil {
		logger.Error("could not create scheduler", zap.Error(err))
		return
	}

	j, err := s.NewJob(
		gocron.DurationJob(15*time.Second),
		gocron.NewTask(func() {
			_, err = monitoring.MeasureExecutionTimeWithCtxError(ctx, collectorClient.Collect, monitoring.WithLogger(logger))
			if err != nil {
				logger.Error("error during collect", zap.Error(err))
			}
		}),
	)
	if err != nil {
		logger.Error("could not run collect", zap.Error(err))
	}

	s.Start()
	_ = j.RunNow()

	<-ctx.Done()
}
