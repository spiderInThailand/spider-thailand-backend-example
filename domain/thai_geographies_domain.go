package domain

import (
	"context"
	"spider-go/model"
)

//go:generate mockgen -source=thai_geographies_domain.go -destination=./mock/thai_geographies_domain.go

type ThaiGeographiesRepository interface {
	GetAllProvince(ctx context.Context) ([]model.Province, error)
	FindProvinceWithProvinceNameEN(ctx context.Context, provinceName string) (model.Province, error)
}

type ThaiGeographiesUsecase interface {
	GetAllProvince(ctx context.Context) ([]model.Province, error)
	GetDistictWithProvinceNameEN(ctx context.Context, provinceNameEN string) ([]model.District, error)
	GetGeographiesBySpiderType(ctx context.Context, family, genus, species string) ([]model.LocationResult, error)
}
