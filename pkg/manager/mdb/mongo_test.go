/**
 * @Time: 2023/10/26 16:48
 * @Author: jzechen
 * @File: mongo_test.go
 * @Software: GoLand toresa
 */

package mdb

import (
	"context"
	"flag"
	"github.com/jzechen/toresa/pkg/manager/config"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"testing"
	"time"
)

var collect *mongo.Collection

func TestMain(m *testing.M) {
	flag.Parse()
	mgo, err := NewMongoDBImpl(context.Background(), &config.MongoConfig{
		Addr:        "mongodb://root:root@172.17.0.1:27017",
		Database:    "test",
		DialTimeout: 3 * time.Second,
	})
	if err != nil {
		panic(err)
	}
	collect, _ = mgo.GetCollection("test", "golangMongoAPI")
	m.Run()
	_ = mgo.Close()
}

func TestMongoImpl_Operation(t *testing.T) {
	docs := []interface{}{
		bson.D{{"_id", "1"}, {"name", "Alice"}, {"age", 18}},
		bson.D{{"_id", "2"}, {"name", "Bob"}, {"age", 16}},
		bson.D{{"_id", "3"}, {"name", "Caren"}, {"age", 20}},
	}
	res, err := collect.InsertMany(context.TODO(), docs)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("insertOne succeed, ids: %v", res.InsertedIDs)

	//filter := bson.D{{"title", "Star Wars"}}
	opts := options.Find().SetSort(bson.D{{"age", 1}})
	cursor, err := collect.Find(context.TODO(), bson.D{}, opts)
	if err != nil {
		t.Fatal(err)
	}

	var results []bson.M
	err = cursor.All(context.TODO(), &results)
	if err != nil {
		t.Fatal(err)
	}
	for index, result := range results {
		t.Logf("[%d] %v", index+1, result)
	}

	dRes, err := collect.DeleteMany(context.TODO(), bson.D{})
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("deleteMany succeed, count: %v", dRes.DeletedCount)
}
