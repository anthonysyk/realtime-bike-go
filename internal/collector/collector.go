package collector

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	"realtime-bike-go/config"
	"realtime-bike-go/internal/store"
	"realtime-bike-go/pkg/opendata"
	"sync"
)

const (
	AppName = "rb-collector"
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

func (c *Client) Collect(ctx context.Context) error {
	stationsChan := c.FetchStations()

	wg := sync.WaitGroup{}
	for i := 0; i < 3; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for station := range stationsChan {
				err := c.store.UpsertOpenDataStation(ctx, station)
				if err != nil {
					c.logger.Error("could not upsert in od_station", zap.Error(err))
				}
			}
		}()
	}
	wg.Wait()

	return nil
}

func (c *Client) FetchStations() chan opendata.Station {
	var (
		hasNext      = true
		offset       = 0
		batchSize    = 100
		stationsChan = make(chan opendata.Station, batchSize)

		total = 0
	)

	go func() {
		defer func() {
			c.logger.Info("total stations processed", zap.Int("total", total))
			close(stationsChan)
		}()

		for hasNext {
			res, err := opendata.GetStationAvailability(batchSize, offset)
			if err != nil {
				c.logger.Error("error during fetch stations", zap.Error(err))
			}
			for _, station := range res.Results {
				total++
				stationsChan <- station
			}
			hasNext = res.HasNext
			offset += len(res.Results)
		}
	}()

	return stationsChan
}

func (c *Client) Close(ctx context.Context) error {
	c.logger.Info("closing collector...")
	return c.store.Close(ctx)
}
