package repository

import (
	"context"
	"spider-go/domain"
	"spider-go/logger"
	"spider-go/model"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type StatisticsRepository struct {
	database       *mongo.Database
	log            *logger.Logger
	collectionName string
}

// FindSpiderStatisticsByFamily implements domain.StatisticsRepository

func NewSpiderStatisticsRepository(db *mongo.Database) domain.StatisticsRepository {
	return &StatisticsRepository{
		database:       db,
		log:            logger.L().Named("StatisticsRepository"),
		collectionName: "spider_statistics",
	}
}

func (r *StatisticsRepository) FindAllSpiderStatistics(ctx context.Context) ([]model.SpiderStatistics, error) {

	coll := r.database.Collection(r.collectionName)

	var result []model.SpiderStatistics

	cursor, err := coll.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}

	if err = cursor.All(context.TODO(), &result); err != nil {
		return nil, err
	}

	return result, nil
}

func (r *StatisticsRepository) FindSpiderStatisticsByFamily(ctx context.Context, family string) (*model.SpiderStatistics, error) {
	log := r.log.WithContext(ctx)

	coll := r.database.Collection(r.collectionName)

	selector := bson.M{
		"family_name": family,
	}

	var result model.SpiderStatistics

	if err := coll.FindOne(ctx, selector).Decode(&result); err != nil {
		log.Errorf("[FindSpiderStatisticsByFamily], find one spider statistics error %+v", err)
		if err == mongo.ErrNoDocuments {
			return nil, ErrorMongoNotFound
		}
		return nil, err
	}

	return &result, nil

}

func (r *StatisticsRepository) UpsertSpiderStatistics(ctx context.Context, familyName string, data model.SpiderStatistics) error {
	log := r.log.WithContext(ctx)

	opts := options.Update().SetUpsert(true)

	selector := bson.M{
		"family_name": familyName,
	}

	updater := bson.M{
		"$set": data,
	}

	coll := r.database.Collection(r.collectionName)

	result, err := coll.UpdateOne(ctx, selector, updater, opts)

	if err != nil {
		log.Errorf("[upsert spider statistics] upsert data error: %+v, result: %+v", err, result)
		return err
	}

	log.Infof("[upsert spider statistics] upsert data result: %+v", result)

	return nil
}

func (r *StatisticsRepository) FindFamilyListWithLimitSizePage(ctx context.Context, page, limit int32) ([]model.SpiderStatistics, error) {

	log := r.log.WithContext(ctx)

	coll := r.database.Collection(r.collectionName)

	var result []model.SpiderStatistics
	opts := options.Find()

	opts.SetSkip(int64(page * limit))
	opts.SetLimit(int64(limit))

	cursor, err := coll.Find(ctx, bson.M{}, opts)
	if err != nil {
		log.Errorf("[FindFamilyListWithLimitSizePage] mongo error: %+v", err)
		return nil, err
	}

	if err = cursor.All(context.TODO(), &result); err != nil {
		log.Errorf("[FindFamilyListWithLimitSizePage] decode to struct failed, error: %+v", err)
		return nil, err
	}

	return result, nil
}
