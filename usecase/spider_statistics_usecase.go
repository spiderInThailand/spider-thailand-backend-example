package usecase

import (
	"context"
	"spider-go/domain"
	"spider-go/logger"
	"spider-go/model"
)

type SpiderStatistics struct {
	SpiderStatisticsRepo domain.StatisticsRepository
	log                  *logger.Logger
}

func NewSpiderStatisticsUsecase(SpiderStatisticsRepo domain.StatisticsRepository) domain.StatisticsUsecase {
	return &SpiderStatistics{
		SpiderStatisticsRepo: SpiderStatisticsRepo,
		log:                  logger.L().Named("SpiderStatisticsUsecase"),
	}
}

func (u *SpiderStatistics) GetSpiderStatisticsList(ctx context.Context) ([]model.SpiderStatistics, error) {
	log := u.log.WithContext(ctx)

	spiderStatistics, err := u.SpiderStatisticsRepo.FindAllSpiderStatistics(ctx)
	if err != nil {
		log.Errorf("find spider statistics error: %+v", err)
		return nil, ErrorMongoTechnicalFail
	}

	return spiderStatistics, nil
}
