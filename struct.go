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
		MongoClient     *mongo.Client
		Srv             bool
		Ctx             context.Context
		MongoHost       string
		MongoUser       string
		MongoPass       string
		MongoDBName     string
		MongoPort       string
		MongoCollection string
	}
)
