package store

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.uber.org/zap"
	"time"
)

const (
	CollectionOpenDataStation       = "od_stations"
	CollectionOpenDataStationEvents = "od_stations_events"
)

type Store struct {
	DB     *mongo.Database
	logger *zap.Logger
}

func New(ctx context.Context, logger *zap.Logger, uri, databaseName, appName string, serverSelectionTimeout time.Duration) (*Store, error) {
	opts := options.Client().
		ApplyURI(uri).
		SetReadPreference(readpref.Primary()).
		SetServerSelectionTimeout(serverSelectionTimeout).
		SetAppName(appName)
	mongoClient, err := mongo.Connect(ctx, opts)
	if err != nil {
		logger.Error("Failed to create mongodb client", zap.String("uri", uri), zap.Error(err))
		return nil, err
	}
	logger.Info("Connected to mongodb database", zap.String("uri", uri))

	store := &Store{DB: mongoClient.Database(databaseName), logger: logger}

	err = store.CreateIndices(ctx)
	if err != nil {
		return nil, err
	}

	return store, nil
}

func (s *Store) Close(ctx context.Context) error {
	s.logger.Info("closing mongo client...")
	return s.DB.Client().Disconnect(ctx)
}

func (s *Store) CreateIndices(ctx context.Context) error {
	_, err := s.DB.Collection(CollectionOpenDataStation).Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys:    bson.M{"stationcode": 1},
		Options: options.Index().SetUnique(true),
	})
	return err
}
