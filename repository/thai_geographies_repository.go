package repository

import (
	"context"
	"spider-go/domain"
	"spider-go/logger"
	"spider-go/model"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type ThaiGeographiesRepository struct {
	database       *mongo.Database
	log            *logger.Logger
	collectionName string
}

func NewThaiGeographiesRepository(db *mongo.Database) domain.ThaiGeographiesRepository {
	return &ThaiGeographiesRepository{
		database:       db,
		log:            logger.L().Named("ThaiGeographiesRepository"),
		collectionName: "thailand_geographies",
	}
}

func (r *ThaiGeographiesRepository) GetAllProvince(ctx context.Context) ([]model.Province, error) {
	log := r.log.WithContext(ctx)

	coll := r.database.Collection(r.collectionName)

	var provinceList []model.Province

	cursor, err := coll.Find(ctx, bson.M{})
	if err != nil {
		log.Errorf("[GetAllProvince] find province cursor failed, error: %+v", err)
		return provinceList, err
	}

	if err := cursor.All(ctx, &provinceList); err != nil {
		log.Errorf("[GetAllProvince] get province from cursor failed, error: %+v", err)
		return provinceList, err
	}

	return provinceList, nil
}

func (r *ThaiGeographiesRepository) FindProvinceWithProvinceNameEN(ctx context.Context, provinceName string) (model.Province, error) {
	log := r.log.WithContext(ctx)

	var province model.Province

	coll := r.database.Collection(r.collectionName)

	selector := bson.M{
		"name_en": provinceName,
	}

	if err := coll.FindOne(ctx, selector).Decode(&province); err != nil {
		log.Errorf("[FindProvinceWithProvinceNameEN] find province with filter %+v failed, error: %+v", selector, err)
		if err == mongo.ErrNoDocuments {
			return province, ErrorMongoNotFound
		}
		return province, err
	}

	return province, nil

}
