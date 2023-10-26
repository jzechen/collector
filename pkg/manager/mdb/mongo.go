/**
 * @Time: 2023/10/20 15:45
 * @Author: jzechen
 * @File: options.go
 * @Software: GoLand collector
 */

package mdb

import (
	"context"
	"github.com/jzechen/toresa/pkg/manager/config"
	"go.mongodb.org/mongo-driver/mongo"
	_ "go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	_ "go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	_ "go.mongodb.org/mongo-driver/mongo/readpref"
	"time"
)

type Interface interface {
	GetCollection(databases, collection string) (*mongo.Collection, error)
	Close() error
}

type MongoImpl struct {
	ctx        context.Context
	cancelFunc context.CancelFunc
	client     *mongo.Client
}

// NewMongoDBImpl initialize a MongoDB connection.
func NewMongoDBImpl(ctx context.Context, cfg *config.MongoConfig) (Interface, error) {
	ctx, cancel := context.WithTimeoutCause(ctx, cfg.DialTimeout, context.DeadlineExceeded)

	// Specify BSON options that cause the driver to fallback to "json"
	// struct tags if "bson" struct tags are missing, marshal nil Go maps as
	// empty BSON documents, and marshals nil Go slices as empty BSON
	// arrays.
	bsonOpts := &options.BSONOptions{
		UseJSONStructTags: true,
		NilMapAsEmpty:     true,
		NilSliceAsEmpty:   true,
	}
	clientOpts := options.Client().ApplyURI(cfg.Addr).SetBSONOptions(bsonOpts)

	c, err := mongo.Connect(ctx, clientOpts)
	if err != nil {
		cancel()
		return nil, err
	}
	err = c.Ping(ctx, readpref.Primary())
	if err != nil {
		return nil, err
	}
	mgo := &MongoImpl{
		ctx:        ctx,
		cancelFunc: cancel,
		client:     c,
	}
	return mgo, nil
}

func (m *MongoImpl) GetCollection(databases, collection string) (*mongo.Collection, error) {
	ctx, cancel := context.WithTimeout(m.ctx, 2*time.Second)
	defer cancel()

	err := m.client.Ping(ctx, readpref.Primary())
	if err != nil {
		return nil, err
	}
	return m.client.Database(databases).Collection(collection), nil
}

func (m *MongoImpl) Close() error {
	err := m.client.Disconnect(m.ctx)
	if err != nil {
		return err
	}
	return nil
}
