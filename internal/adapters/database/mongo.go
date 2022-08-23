package database

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

type DB struct {
	client   *mongo.Client
	database *mongo.Database
}

func NewMongoClient(ctx context.Context, uri string, dbName string) (*DB, error) {
	client, err := mongo.NewClient(options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}
	if err := client.Connect(ctx); err != nil {
		return nil, err
	}

	return &DB{client, client.Database(dbName)}, err
}

func (db *DB) SaveEvent(ctx context.Context, collection string, docs interface{}) error {
	if _, err := db.database.Collection(collection).InsertOne(ctx, docs); err != nil {
		return err
	}
	return nil
}

func (db *DB) CloseDbConnection(ctx context.Context) {
	err := db.client.Disconnect(ctx)
	if err != nil {
		log.Printf("failed to disconnect from mongo: %s", err.Error())
	}
}
