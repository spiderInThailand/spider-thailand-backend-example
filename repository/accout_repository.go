package repository

import (
	"context"
	"spider-go/domain"
	"spider-go/logger"
	"spider-go/model"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type Account struct {
	database       *mongo.Database
	log            *logger.Logger
	collectionName string
}

func NewAccountRepository(db *mongo.Database) domain.AccountRepository {
	return &Account{
		database:       db,
		log:            logger.L().Named("AccountRepository"),
		collectionName: "account",
	}
}

func (r *Account) CreateAccout(ctx context.Context, acc model.Account) (err error) {
	log := r.log.WithContext(ctx)

	coll := r.database.Collection(r.collectionName)

	_, err = coll.InsertOne(ctx, acc)
	if err != nil {
		log.Errorf("insert account error: %+v", err)
		return err
	}

	// log.Debugf("insert account result: %+v", result)

	return nil
}

func (r *Account) FindAccountByUsername(ctx context.Context, username string) (AccountInfo *model.Account, err error) {
	selector := bson.M{
		"username": username,
	}

	coll := r.database.Collection(r.collectionName)

	var accountInfo model.Account

	if err := coll.FindOne(ctx, selector).Decode(&accountInfo); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, ErrorMongoNotFound
		}
		return nil, err
	}

	return &accountInfo, nil
}
