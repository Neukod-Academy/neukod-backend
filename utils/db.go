package utils

import (
	"context"
	"fmt"
	"time"

	"github.com/Neukod-Academy/neukod-backend/pkg/env"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.mongodb.org/mongo-driver/v2/mongo/readpref"
)

type Mongo struct {
	Client  *mongo.Client
	Context context.Context
	Cancel  context.CancelFunc
}

type Stored struct {
	Data []interface{}
}

func (m *Mongo) CreateClient(uri string) error {
	var err error
	m.Context, m.Cancel = context.WithTimeout(context.Background(), 2*time.Second)
	if m.Client, err = mongo.Connect(options.Client().ApplyURI(uri)); err != nil {
		return fmt.Errorf("failed to connect: %s", err.Error())
	}

	if err = m.Client.Ping(m.Context, readpref.Primary()); err != nil {
		return fmt.Errorf("failed to connect: %s", err.Error())
	}
	return nil
}

func (m *Mongo) CloseClientDB() error {
	if m.Client == nil {
		return nil
	}
	if err := m.Client.Disconnect(m.Context); err != nil {
		return fmt.Errorf("failed to disconnect the client with the database: %s", err.Error())
	}
	fmt.Println("The client is already disconnected")
	return nil
}

func (m *Mongo) FindAllData(db_name string, coll_name string) (*Stored, error) {
	if err := m.CreateClient(env.MONGO_URI); err != nil {
		return nil, err
	}
	coll := m.Client.Database(db_name).Collection(coll_name)
	filter := bson.M{}
	options := options.Find()
	var stored Stored
	if cursor, err := coll.Find(m.Context, filter, options); err != nil {
		return nil, err
	} else {
		if err := cursor.All(m.Context, &stored); err != nil {
			return nil, err
		}
	}
	return &stored, nil
}

func (m *Mongo) InsertNewData(db_name string, coll_name string, new_data interface{}) (*mongo.InsertOneResult, error) {
	if err := m.CreateClient(env.MONGO_URI); err != nil {
		return nil, err
	}
	coll := m.Client.Database(db_name).Collection(coll_name)
	opts := options.InsertOne()
	if result, err := coll.InsertOne(m.Context, new_data, opts); err != nil {
		return nil, err
	} else {
		return result, nil
	}
}
