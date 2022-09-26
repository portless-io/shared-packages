package mongodb

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoClient struct {
	Client *mongo.Client
}

func NewMongoClient(url string) MongoClient {
	if url == "" {
		url = "mongodb://localhost:27017?retryWrites=false"
	}

	opt := options.Client().ApplyURI(url).SetRetryWrites(false)

	client, err := mongo.NewClient(opt)

	if err != nil {
		panic(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	defer cancel()
	err = client.Connect(ctx)

	if err != nil {
		panic(err)
	}

	return MongoClient{
		Client: client,
	}
}

func (m *MongoClient) Database(dbName string) *mongo.Database {
	return m.Client.Database(dbName)
}

func (m *MongoClient) Close() {
	m.Client.Disconnect(context.Background())
}

func StringToID(id string) primitive.ObjectID {
	ID, _ := primitive.ObjectIDFromHex(id)
	return ID
}

func IsNoDocumentError(err error) bool {
	return err == mongo.ErrNoDocuments
}
