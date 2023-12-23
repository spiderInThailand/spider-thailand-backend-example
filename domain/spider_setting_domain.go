package domain

import (
	"context"
	api_model "spider-go/api/model"
)

//go:generate mockgen -source=spider_setting_domain.go -destination=./mock/spider_setting_domain.go

type UploadImageUsecase interface {
	UploadImageSpiderUsecase(ctx context.Context, spiderUUID string, listImageEncode64 []string) error
}

type DeleteSpiderInfoUsecase interface {
	DeleteSpiderInfoUsecase(ctx context.Context, spider_uuid string) error
}

type UpdateSpiderInfoUsecase interface {
	UpdateSpiderInfoUsecase(ctx context.Context, spiderInfo api_model.SpiderInfo) error
}

type RemoveSpiderImageUsecase interface {
	RemoveSpiderImageBySpiderImageNameList(ctx context.Context, spiderUUID string, spiderImageList []string) error
}
