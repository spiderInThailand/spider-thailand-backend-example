package repository

import (
	"context"
	"fmt"
	"spider-go/domain"
	"spider-go/logger"
	"spider-go/model"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type SpiderRepository struct {
	database       *mongo.Database
	log            *logger.Logger
	collectionName string
}

var (
	ErrorSpiderRepositoryDeleteSpiderIsZero = fmt.Errorf("delete spider result is zero")
)

func NewSpiderRepository(db *mongo.Database) domain.SpiderRepository {
	return &SpiderRepository{
		database:       db,
		log:            logger.L().Named("SpiderRepository"),
		collectionName: "spider",
	}
}

func (r *SpiderRepository) InsertNewSpider(ctx context.Context, data model.SpiderInfo) error {
	log := r.log.WithContext(ctx)

	log.Infof("[InsertNewSpider] start with param: %+v", data)

	coll := r.database.Collection(r.collectionName)

	result, err := coll.InsertOne(ctx, data)
	log.Debugf("[InsertNewSpider] insert data result: %+v", result)
	if err != nil {
		return err
	}

	return nil
}

func (r *SpiderRepository) FindSpiderByUUID(ctx context.Context, spiderUUID string) (*model.SpiderInfo, error) {

	log := r.log.WithContext(ctx)

	coll := r.database.Collection(r.collectionName)

	selector := bson.M{
		"spider_uuid": spiderUUID,
	}

	log.Infof("[find spider by uuid] find spider with selector: %v", selector)

	var resultSpiderInfo model.SpiderInfo

	if err := coll.FindOne(ctx, selector).Decode(&resultSpiderInfo); err != nil {
		return nil, err
	}

	return &resultSpiderInfo, nil
}

func (r *SpiderRepository) UpdateImageFileToSpiderInfo(ctx context.Context, filesName []string, spiderUUID string) error {
	log := r.log.WithContext(ctx)

	coll := r.database.Collection(r.collectionName)

	tn := time.Now()

	selector := bson.M{
		"spider_uuid": spiderUUID,
	}

	updater := bson.M{
		"$set": bson.M{
			"image_file": filesName,
			"updated_at": tn,
		},
	}

	result, err := coll.UpdateOne(ctx, selector, updater)
	if err != nil {
		log.Errorf("[UpdateImageFileToSpiderInfo] update mongo failed, error: %v", err)
		return err
	}

	log.Infof("[update image file to spider info] result: %+v", result)
	return nil
}

func (r *SpiderRepository) FindSpiderByUUIDAndStatus(ctx context.Context, spiderUUID string, isStatusActive bool) (*model.SpiderInfo, error) {

	log := r.log.WithContext(ctx)

	coll := r.database.Collection(r.collectionName)

	selector := bson.M{"spider_uuid": spiderUUID}

	if isStatusActive {
		selector["status"] = model.SPIDER_INFO_STATUS_ACTIVE
	}

	log.Infof("[find spider by uuid and status] find spider status active `%v` with selector: %v", isStatusActive, selector)

	var resultSpiderInfo model.SpiderInfo

	if err := coll.FindOne(ctx, bson.M{"spider_uuid": spiderUUID}).Decode(&resultSpiderInfo); err != nil {
		log.Errorf("[find spider by uuid and status] error find one in mongo, error: %v", err)
		if err == mongo.ErrNoDocuments {
			return nil, ErrorMongoNotFound
		}
		return nil, err
	}

	return &resultSpiderInfo, nil
}

func (r *SpiderRepository) findManySpiderWithCondition(ctx context.Context, filter bson.M, opts *options.FindOptions) (*mongo.Cursor, error) {

	coll := r.database.Collection(r.collectionName)

	cursor, err := coll.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}

	return cursor, nil
}

func (r *SpiderRepository) FindAllSpiderListWithActive(ctx context.Context) ([]model.SpiderInfo, error) {
	log := r.log.WithContext(ctx)
	var spiderInfoList []model.SpiderInfo

	selector := bson.M{
		"status": model.SPIDER_INFO_STATUS_ACTIVE,
	}

	cursor, err := r.findManySpiderWithCondition(ctx, selector, nil)
	if err != nil {
		log.Errorf("[FindAllSpiderListWithActive] find spider info list with selector %v, error: %+v", selector, err)
		return []model.SpiderInfo{}, err
	}

	if err = cursor.All(ctx, &spiderInfoList); err != nil {
		log.Errorf("[FindAllSpiderListWithActive] get spider info data from cursor error: %+v", err)
		return []model.SpiderInfo{}, err
	}

	return spiderInfoList, nil
}

func (r *SpiderRepository) FindAllSpiderListManager(ctx context.Context, page, limit int) ([]model.SpiderInfo, error) {
	log := r.log.WithContext(ctx)
	var spiderInfoList []model.SpiderInfo

	// set options
	opts := options.Find()

	opts.SetSort(bson.M{
		"_id": 1,
	})
	opts.SetSkip(int64(page * limit))
	opts.SetLimit(int64(limit))

	log.Infof("[FindAllSpiderListWithActive] find spider with skip `%v` and limit `%v`", int64(page), int64(limit))

	cursor, err := r.findManySpiderWithCondition(ctx, bson.M{}, opts)
	if err != nil {
		log.Errorf("[FindAllSpiderListWithActive] find spider info list error: %+v", err)
		return []model.SpiderInfo{}, err
	}

	if err = cursor.All(ctx, &spiderInfoList); err != nil {
		log.Errorf("[FindAllSpiderListWithActive] get spider info data from cursor error: %+v", err)
		return []model.SpiderInfo{}, err
	}

	return spiderInfoList, nil
}

func (r *SpiderRepository) DeleteSpiderInfoWithSpiderUUID(ctx context.Context, spiderUUID string) error {
	log := r.log.WithContext(ctx)

	log.Infof("[DeleteSpiderInfoWithSpiderUUID] spider_uuid: %v", spiderUUID)

	selector := bson.M{
		"spider_uuid": spiderUUID,
	}

	coll := r.database.Collection(r.collectionName)

	result, err := coll.DeleteOne(ctx, selector)
	if err != nil {
		log.Errorf("[DeleteSpiderInfoWithSpiderUUID] delete one at spider_uuid `%v` failed, result: %+v, error: %+v", spiderUUID, result, err)
		return err
	}

	if result.DeletedCount == 0 {
		log.Errorf("[DeleteSpiderInfoWithSpiderUUID] delete spider result is `%v`", result.DeletedCount)
		return ErrorSpiderRepositoryDeleteSpiderIsZero
	}

	log.Infof("[DeleteSpiderInfoWithSpiderUUID] delete one at spider_uuid `%v` successfull, result: %+v", spiderUUID, result)

	return nil
}

func (r *SpiderRepository) UpdateSpiderInfo(ctx context.Context, spiderUUID string, spiderInfo model.SpiderInfo) (bool, error) {
	log := r.log.WithContext(ctx)

	log.Infof("[UpdateSpiderInfo] update spider where spider_uuid: %v", spiderUUID)

	selector := bson.M{
		"spider_uuid": spiderUUID,
	}

	updater := bson.M{
		"$set": spiderInfo,
	}

	coll := r.database.Collection(r.collectionName)

	result, err := coll.UpdateOne(ctx, selector, updater)
	if err != nil {
		log.Errorf("[UpdateSpiderInfo] update spider info failed, error: %+v", err)
		return false, err
	}

	if result.MatchedCount == 0 {
		log.Warnf("[UpdateSpiderInfo] update spider info is zero update")
		return false, nil
	}

	return true, nil
}

func (r *SpiderRepository) FindSpiderInfoListByGeographies(ctx context.Context, province, district, position string) ([]model.SpiderInfo, error) {
	log := r.log.WithContext(ctx)
	var spiderInfoList []model.SpiderInfo

	opts := options.Find()

	geographiesCondition := bson.M{}

	if province != "" {
		geographiesCondition["province"] = bson.M{
			"$regex":   province,
			"$options": "i",
		}
	}

	if district != "" {
		geographiesCondition["district"] = bson.M{
			"$regex":   district,
			"$options": "i",
		}
	}

	if position != "" {
		geographiesCondition["position"] = bson.M{
			"$elemMatch": bson.M{
				"name": bson.M{
					"$regex":   position,
					"$options": "i",
				},
			},
		}
	}

	selector := bson.M{
		"address": bson.M{
			"$elemMatch": geographiesCondition,
		},
	}

	log.Infof("[FindSpiderInfoListByGeographies] find spider info with condition: %+v", selector)

	cursor, err := r.findManySpiderWithCondition(ctx, selector, opts)
	if err != nil {
		log.Errorf("[FindSpiderInfoListByGeographies] find spider info list error: %+v", err)
		return []model.SpiderInfo{}, err
	}

	if err = cursor.All(ctx, &spiderInfoList); err != nil {
		log.Errorf("[FindSpiderInfoListByGeographies] get spider info data from cursor error: %+v", err)
		return []model.SpiderInfo{}, err
	}

	if len(spiderInfoList) == 0 {
		return []model.SpiderInfo{}, ErrorMongoNotFound
	}

	return spiderInfoList, nil
}

func (r *SpiderRepository) FindSpiderInfoBySpiderType(ctx context.Context, family, genus, species string, isLimitPage bool, page, limit int32) ([]model.SpiderInfo, error) {
	log := r.log.WithContext(ctx)

	var spiderInfoList []model.SpiderInfo

	opts := options.Find()

	if isLimitPage {
		opts.SetSkip(int64(page * limit))
		opts.SetLimit(int64(limit))
	}

	selector := bson.M{}

	if family != "" {
		selector["family"] = family
	}

	if genus != "" {
		selector["genus"] = genus
	}

	if species != "" {
		selector["species"] = species
	}

	log.Infof("[FindSPiderInfoBySpiderType] find spider info with condition: %+v", selector)

	cursor, err := r.findManySpiderWithCondition(ctx, selector, opts)
	if err != nil {
		log.Errorf("[FindSPiderInfoBySpiderType] find spider info list error: %+v", err)
		return []model.SpiderInfo{}, err
	}

	if err = cursor.All(ctx, &spiderInfoList); err != nil {
		log.Errorf("[FindSPiderInfoBySpiderType] get spider info data from cursor error: %+v", err)
		return []model.SpiderInfo{}, err
	}

	if len(spiderInfoList) == 0 {
		return []model.SpiderInfo{}, ErrorMongoNotFound
	}

	return spiderInfoList, nil
}

func (r *SpiderRepository) FindSpiderInfoByLocality(ctx context.Context, locality string, page, limit int32) ([]model.SpiderInfo, error) {
	log := r.log.WithContext(ctx)

	// set options
	opts := options.Find()

	opts.SetSort(bson.M{
		"_id": 1,
	})
	opts.SetSkip(int64(page * limit))
	opts.SetLimit(int64(limit))

	// set selecter
	selecter := bson.M{
		"address": bson.M{
			"$elemMatch": bson.M{
				"position": bson.M{
					"$elemMatch": bson.M{
						"name": locality,
					},
				},
			},
		},
	}

	log.Infof("[FindSpiderInfoByLocality] find spider with skip `%v` and limit `%v`", int64(page), int64(limit))

	var spiderInfoList []model.SpiderInfo
	cursor, err := r.findManySpiderWithCondition(ctx, selecter, opts)
	if err != nil {
		log.Errorf("[FindSpiderInfoByLocality] find spider info list error: %+v", err)
		return []model.SpiderInfo{}, err
	}

	if err = cursor.All(ctx, &spiderInfoList); err != nil {
		log.Errorf("[FindSpiderInfoByLocality] get spider info data from cursor error: %+v", err)
		return []model.SpiderInfo{}, err
	}

	if len(spiderInfoList) == 0 {
		return []model.SpiderInfo{}, ErrorMongoNotFound
	}

	return spiderInfoList, nil

}

func (r *SpiderRepository) FindSpiderInfoByFirstFamilyOrGenus(ctx context.Context, field, value string) ([]model.SpiderInfo, error) {
	log := r.log.WithContext(ctx)

	// set options
	opts := options.Find()

	selecter := bson.M{
		field: value,
	}

	var spiderInfoList []model.SpiderInfo
	cursor, err := r.findManySpiderWithCondition(ctx, selecter, opts)
	if err != nil {
		log.Errorf("[FindSpiderInfoByFirstFamilyOrGenus] find spider info list error: %+v", err)
		return []model.SpiderInfo{}, err
	}

	if err = cursor.All(ctx, &spiderInfoList); err != nil {
		log.Errorf("[FindSpiderInfoByFirstFamilyOrGenus] get spider info data from cursor error: %+v", err)
		return []model.SpiderInfo{}, err
	}

	if len(spiderInfoList) == 0 {
		return []model.SpiderInfo{}, ErrorMongoNotFound
	}

	return spiderInfoList, nil

}
