package domain

import (
	"context"
	api_model "spider-go/api/model"
	"spider-go/model"
)

//go:generate mockgen -source=spider_domain.go -destination=./mock/spider_domain.go

type SpiderRepository interface {
	InsertNewSpider(ctx context.Context, data model.SpiderInfo) error
	FindSpiderByUUID(ctx context.Context, spiderUUID string) (*model.SpiderInfo, error)
	UpdateImageFileToSpiderInfo(ctx context.Context, filesName []string, spiderUUID string) error
	FindSpiderByUUIDAndStatus(ctx context.Context, spiderUUID string, isStatusActive bool) (*model.SpiderInfo, error)
	FindAllSpiderListWithActive(ctx context.Context) ([]model.SpiderInfo, error)
	FindAllSpiderListManager(ctx context.Context, page, limit int) ([]model.SpiderInfo, error)
	DeleteSpiderInfoWithSpiderUUID(ctx context.Context, spiderUUID string) error
	UpdateSpiderInfo(ctx context.Context, spiderUUID string, spiderInfo model.SpiderInfo) (bool, error)
	FindSpiderInfoListByGeographies(ctx context.Context, province, district, position string) ([]model.SpiderInfo, error)
	FindSpiderInfoBySpiderType(ctx context.Context, family, genus, species string, isLimitPage bool, page, limit int32) ([]model.SpiderInfo, error)
	FindSpiderInfoByLocality(ctx context.Context, locality string, page, limit int32) ([]model.SpiderInfo, error)
	FindSpiderInfoByFirstFamilyOrGenus(ctx context.Context, field, value string) ([]model.SpiderInfo, error)
}

type RegisterSpiderUsecase interface {
	Register(ctx context.Context, req api_model.SpiderInfo, username string) (string, error)
}

type SpiderInfoUsecase interface {
	GetSpiderInfoUsecase(ctx context.Context, spiderUUID, username string) (*model.SpiderInfo, error)
	GetSpiderImagesUsecase(ctx context.Context, fileImages []string) ([]model.SpiderImageList, error)
	GetSpiderInfoListManager(ctx context.Context, usecase string, page, limit int) ([]model.SpiderInfo, error)
	GetSpiderInfoListByGeographies(ctx context.Context, province, district, position string) ([]model.SpiderInfo, error)
	GetSpiderInfoListByLocality(ctx context.Context, locality string, page, size int32) ([]model.SpiderInfo, error)
	GetSpiderListBySpiderTypeUsecase(ctx context.Context, param model.GetSpiderListBySpiderTypeParam) ([]model.SpiderInfo, error)
}
