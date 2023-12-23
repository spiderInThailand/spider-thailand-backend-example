package database

import (
	"context"
	"fmt"
	"spider-go/config"
	"spider-go/logger"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	DB     *mongo.Database
	Client *mongo.Client
)

type MongoDB struct {
	config *config.MongoConfig
	log    *logger.Logger
}

func NewMongoDB(config *config.MongoConfig) *MongoDB {
	return &MongoDB{
		config: config,
		log:    logger.L().Named("MongoDB"),
	}
}

func (m *MongoDB) Connect() {
	var err error

	clientOption := options.Client().ApplyURI(m.toMongoURL())

	// set connect mongodb
	Client, err = mongo.Connect(context.TODO(), clientOption)
	if err != nil {
		m.log.Errorf("connect mongoDB failed, error %+v", err)
		panic(err)
	}

	// test connection
	if err := Client.Ping(context.TODO(), nil); err != nil {
		m.log.Errorf("ping mongoDB failed, error %+v", err)
		panic(err)
	}

	m.log.Info("mongoDB connection successful")
}

func (m *MongoDB) SetDB() {
	DB = Client.Database(m.config.DbName)
}

func (m *MongoDB) Close() {
	if err := Client.Disconnect(context.TODO()); err != nil {
		panic(err)
	}
}

func (m *MongoDB) toMongoURL() string {
	return fmt.Sprintf("mongodb://%s:%s@%s",
		m.config.Username,
		m.config.Password,
		m.config.HostPort,
	)
}
