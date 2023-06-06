package logrus_mongodb

import (
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"net/url"
)

func New(opt Option) (*Hooker, error) {
	//connect to mongodb
	protocol := "mongodb"
	if opt.Srv {
		protocol = "mongodb+srv"
	}
	uri := fmt.Sprintf("%s://%s:%s@%s:%s",
		protocol,
		url.QueryEscape(opt.MongoUser),
		url.QueryEscape(opt.MongoPass),
		opt.MongoHost,
		opt.MongoPort,
	)

	ctx := context.Background()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}
	return &Hooker{db: client.Database(opt.MongoDBName), opt: &opt}, nil
}

func (h *Hooker) Fire(entry *logrus.Entry) error {
	data := make(logrus.Fields)
	data["Level"] = entry.Level.String()
	data["Time"] = entry.Time
	data["Message"] = entry.Message

	for k, v := range entry.Data {
		if errData, isError := v.(error); logrus.ErrorKey == k && v != nil && isError {
			data[k] = errData.Error()
		} else {
			data[k] = v
		}
	}
	ctx := context.Background()
	_, mgoErr := h.db.Collection(h.opt.MongoCollection).InsertOne(ctx, data)

	if mgoErr != nil {
		return fmt.Errorf("failed to send log entry to mongodb: %v", mgoErr)
	}

	return nil
}

func (h *Hooker) Levels() []logrus.Level {
	return logrus.AllLevels
}
