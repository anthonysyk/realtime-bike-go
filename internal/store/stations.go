package store

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"realtime-bike-go/pkg/opendata"
)

func (s *Store) UpsertOpenDataStation(ctx context.Context, station opendata.Station) error {
	filter := bson.M{"stationcode": station.Stationcode}
	update := bson.M{"$set": station}
	opts := options.Update().SetUpsert(true)

	_, err := s.DB.Collection(CollectionOpenDataStation).UpdateOne(ctx, filter, update, opts)
	if err != nil {
		return err
	}
	return nil
}

func (s *Store) ListStations(ctx context.Context, name, code string) ([]opendata.Station, error) {
	filter := bson.M{}
	if code != "" {
		filter["code"] = code
	}
	if name != "" {
		filter["name"] = bson.M{
			"$regex":   name,
			"$options": "i", // not case-sensitive
		}
	}
	cursor, err := s.DB.Collection(CollectionOpenDataStation).Find(ctx, filter)
	if err != nil {
		return nil, err
	}

	var results []opendata.Station
	for cursor.Next(ctx) {
		var station opendata.Station
		if cursor.Decode(&station); err != nil {
			return nil, err
		}
		results = append(results, station)
	}

	return results, nil
}
