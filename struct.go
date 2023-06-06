package logrus_mongodb

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
)

type (
	Hooker struct {
		db  *mongo.Database
		opt *Option
		c   context.Context
	}
	Option struct {
		Srv             bool
		MongoHost       string
		MongoUser       string
		MongoPass       string
		MongoDBName     string
		MongoPort       string
		MongoCollection string
	}
)
