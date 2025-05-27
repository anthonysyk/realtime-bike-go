package watcher

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"realtime-bike-go/config"
	"realtime-bike-go/internal/store"

	"go.mongodb.org/mongo-driver/bson"
	moption "go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

const (
	AppName = "rb-watcher"
)

type Client struct {
	logger *zap.Logger
	store  *store.Store
}

func New(ctx context.Context, log *zap.Logger, cfg config.MongoDB) (*Client, error) {
	st, err := store.New(ctx, log, cfg.URI, cfg.DatabaseName, AppName, cfg.ServerSelectionTimeout)
	if err != nil {
		return nil, fmt.Errorf("could not connect to MongoDB: %v", err)
	}
	return &Client{
		store:  st,
		logger: log,
	}, nil
}

func (c *Client) Watch(ctx context.Context, pipeline bson.A, config config.MongoDBOptions) (chan store.ChangeEvent, error) {
	opts := &moption.ChangeStreamOptions{
		BatchSize:    &config.BatchSize,
		MaxAwaitTime: &config.MaxAwaitTime,
		StartAtOperationTime: &primitive.Timestamp{
			T: config.StartAtOperationTimeT,
			I: config.StartAtOperationTimeI,
		},
	}
	opts.SetFullDocument(moption.UpdateLookup)

	stream, err := c.store.DB.Collection(store.CollectionOpenDataStation).Watch(ctx, pipeline, opts)
	if err != nil {
		return nil, err
	}

	streamChan := make(chan store.ChangeEvent)
	go func() {
		defer func() {
			close(streamChan)
			_ = stream.Close(ctx)
		}()
		for stream.Next(ctx) {
			var event store.ChangeEvent
			if err := stream.Decode(&event); err != nil {
				c.logger.Info("could not decode event:", zap.Error(err))
				continue
			}
			select {
			case streamChan <- event:
			case <-ctx.Done():
				c.logger.Info("context done, stopping change stream watcher")
				return
			}
		}
	}()

	return streamChan, nil
}

func (c *Client) SaveEvents(ctx context.Context, eventsChan chan store.ChangeEvent) error {
	go func() {
		for event := range eventsChan {
			err := c.store.StoreChangeEvent(ctx, store.CollectionOpenDataStationEvents, event)
			if err != nil {
				c.logger.Error("could not insert station", zap.Error(err))
			}
		}
	}()

	return nil
}

func (c *Client) Close(ctx context.Context) error {
	c.logger.Info("closing watcher...")
	return c.store.Close(ctx)
}
