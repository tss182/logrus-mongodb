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
	var client *mongo.Client
	if opt.MongoClient != nil {
		client = opt.MongoClient
	} else {
		//connect to mongodb
		protocol := "mongodb"
		uri := fmt.Sprintf("%s://%s:%s@%s:%s",
			protocol,
			url.QueryEscape(opt.MongoUser),
			url.QueryEscape(opt.MongoPass),
			opt.MongoHost,
			opt.MongoPort,
		)
		if opt.Srv {
			protocol = "mongodb+srv"
			uri = fmt.Sprintf("%s://%s:%s@%s",
				protocol,
				url.QueryEscape(opt.MongoUser),
				url.QueryEscape(opt.MongoPass),
				opt.MongoHost,
			)
		}
		ctx := context.Background()
		c, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
		if err != nil {
			return nil, err
		}
		client = c
		defer func() {
			_ = client.Disconnect(ctx)
		}()
	}

	return &Hooker{db: client.Database(opt.MongoDBName), opt: &opt}, nil
}

func (h *Hooker) Fire(entry *logrus.Entry) error {
	data := make(logrus.Fields)
	data["level"] = entry.Level.String()
	data["time"] = entry.Time
	data["msg"] = entry.Message
	data["file"] = entry.Caller.File
	data["function"] = entry.Caller.Function
	data["line"] = entry.Caller.Line

	for k, v := range entry.Data {
		if errData, isError := v.(error); logrus.ErrorKey == k && v != nil && isError {
			data[k] = errData.Error()
		} else {
			data[k] = v
		}
	}
	_, mgoErr := h.db.Collection(h.opt.MongoCollection).InsertOne(context.Background(), data)

	if mgoErr != nil {
		return fmt.Errorf("failed to send log entry to mongodb: %v", mgoErr)
	}

	return nil
}

func (h *Hooker) Levels() []logrus.Level {
	return logrus.AllLevels
}
