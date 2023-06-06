package logrus_mongodb

import (
	"go.mongodb.org/mongo-driver/mongo"
)

type (
	Hooker struct {
		db  *mongo.Database
		opt *Option
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
