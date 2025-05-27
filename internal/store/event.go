package store

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
	"time"
)

type documentKey struct {
	ID primitive.ObjectID `bson:"_id"`
}

type ChangeEvent struct {
	ID                interface{} `bson:"_id"`
	Operation         string      `bson:"operationType"`
	Document          bson.M      `bson:"fullDocument"`
	Namespace         bson.M      `bson:"ns"`
	NewCollectionName bson.M      `bson:"to,omitempty"`
	DocumentKey       documentKey `bson:"documentKey"`
	Updates           bson.M      `bson:"updateDescription,omitempty"`
	ClusterTime       time.Time   `bson:"clusterTime"`
	Transaction       int64       `bson:"txnNumber,omitempty"`
	SessionID         bson.M      `bson:"lsid,omitempty"`
}

func (e ChangeEvent) DocumentID() (string, error) {
	id := e.DocumentKey.ID
	if id.IsZero() {
		return "", fmt.Errorf("documentKey should not be empty")
	}
	return id.Hex(), nil
}

func (s *Store) StoreChangeEvent(ctx context.Context, collection string, event ChangeEvent) error {
	if event.ID == nil {
		return fmt.Errorf("resume token (_id) is nil – cannot store event")
	}
	_, err := s.DB.Collection(collection).InsertOne(ctx, event)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			s.logger.Error("Event with resume token already exists, skipping...", zap.Any("id", event.ID))
			return nil
		}
		return fmt.Errorf("failed to insert change event: %w", err)
	}
	return nil
}

func (s *Store) GetLastResumeToken(ctx context.Context, collection string) (bson.Raw, error) {
	var lastEvent bson.M
	err := s.DB.Collection(collection).FindOne(
		ctx,
		bson.M{},
		options.FindOne().SetSort(bson.D{{"_id", -1}}),
	).Decode(&lastEvent)

	if err != nil {
		return nil, err
	}

	return lastEvent["_id"].(bson.Raw), nil
}
