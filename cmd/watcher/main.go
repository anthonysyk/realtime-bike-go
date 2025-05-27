package main

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.uber.org/zap"
	"realtime-bike-go/config"
	"realtime-bike-go/internal/watcher"
	"realtime-bike-go/pkg/graceful"
	"realtime-bike-go/pkg/rblogger"
)

var (
	configPrefix = "rb-watcher"
)

func main() {
	ctx, _ := graceful.WaitForSignalContext(context.Background())

	logger := rblogger.New()
	cfg := config.NewBase(ctx, configPrefix)
	watcherClient, err := watcher.New(ctx, logger, cfg.MongoDB)
	if err != nil {
		logger.Error("could not create watcher", zap.Error(err))
		return
	}
	defer watcherClient.Close(ctx)

	pipeline := bson.A{}

	eventsChan, err := watcherClient.Watch(ctx, pipeline, cfg.MongoDB.Options)
	if err != nil {
		logger.Error("could not fetch watch events", zap.Error(err))
		return
	}

	err = watcherClient.SaveEvents(ctx, eventsChan)
	if err != nil {
		logger.Error("could not save events", zap.Error(err))
		return
	}

	<-ctx.Done()
}
