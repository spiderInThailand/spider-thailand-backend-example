package domain

import (
	"context"
	"spider-go/model"
)

//go:generate mockgen -source=spider_statistics_domain.go -destination=./mock/spider_statistics_domain.go
type StatisticsRepository interface {
	FindAllSpiderStatistics(ctx context.Context) ([]model.SpiderStatistics, error)
	FindSpiderStatisticsByFamily(ctx context.Context, family string) (*model.SpiderStatistics, error)
	UpsertSpiderStatistics(ctx context.Context, familyName string, data model.SpiderStatistics) error
	FindFamilyListWithLimitSizePage(ctx context.Context, page, limit int32) ([]model.SpiderStatistics, error)
}

type StatisticsUsecase interface {
	GetSpiderStatisticsList(ctx context.Context) ([]model.SpiderStatistics, error)
}

type GetFamilyListUsecase interface {
	Execute(ctx context.Context, page, size int32) ([]model.FamilyList, error)
}
