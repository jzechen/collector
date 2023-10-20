/**
 * @Time: 2023/10/20 15:45
 * @Author: jzechen
 * @File: options.go
 * @Software: GoLand collector
 */

package mdb

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	_ "go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	_ "go.mongodb.org/mongo-driver/mongo/options"
	_ "go.mongodb.org/mongo-driver/mongo/readpref"
	"time"
)

type MongoImpl struct {
	ctx        context.Context
	cancelFunc context.CancelFunc
	client     *mongo.Client
}

// NewMongoDBImpl initialize a MongoDB connection.
func NewMongoDBImpl(ctx context.Context) (*MongoImpl, error) {
	ctx, cancel := context.WithTimeoutCause(ctx, 10*time.Second, context.DeadlineExceeded)
	c, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		cancel()
		return nil, err
	}
	mgo := &MongoImpl{
		ctx:        ctx,
		cancelFunc: cancel,
		client:     c,
	}
	return mgo, nil
}

//func GetDB() (*gorm.DB, error) {
//	sqlDB, err := DB.DB()
//	if err != nil {
//		return nil, err
//	}
//	if err = sqlDB.Ping(); err != nil {
//		sqlDB.Close()
//		return nil, err
//	}
//
//	return DB, nil
//}

func (m *MongoImpl) Close() error {
	err := m.client.Disconnect(m.ctx)
	if err != nil {
		return err
	}
	return nil
}
