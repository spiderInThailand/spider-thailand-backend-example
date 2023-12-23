package usecase

import (
	"context"
	"spider-go/domain"
	"spider-go/logger"
	"spider-go/model"
	"spider-go/repository"
)

type GetFamilyListUsecase struct {
	SpiderStatisticsRepo domain.StatisticsRepository
	SpiderRepo           domain.SpiderRepository
	log                  *logger.Logger
}

func NewGetFamilyListUsecase(SpiderStatisticsRepo domain.StatisticsRepository, SpiderRepo domain.SpiderRepository) *GetFamilyListUsecase {
	return &GetFamilyListUsecase{
		SpiderStatisticsRepo: SpiderStatisticsRepo,
		SpiderRepo:           SpiderRepo,
		log:                  logger.L().Named("GetFamilyList"),
	}
}

func (u *GetFamilyListUsecase) Execute(ctx context.Context, page, size int32) ([]model.FamilyList, error) {
	log := u.log.WithContext(ctx)

	spiderStatisticeList, err := u.SpiderStatisticsRepo.FindFamilyListWithLimitSizePage(ctx, page, size)
	if err != nil {
		log.Errorf("[get family list usecase] find family list from statistice failed, error: %+v", err)
		return nil, ErrorMongoTechnicalFail
	}

	var familyList []model.FamilyList

	for _, value := range spiderStatisticeList {
		var tempFamilyList model.FamilyList

		tempFamilyList.Family = value.FamilyName

		spiderInfoList, err := u.SpiderRepo.FindSpiderInfoByFirstFamilyOrGenus(ctx, "family", value.FamilyName)
		if err != nil {
			log.Errorf("[get family list usecase] find one spider info failed, error: %+v", err)
			if err == repository.ErrorMongoNotFound {
				tempFamilyList.Author = "N/A"
				tempFamilyList.Quantity = int32(0)
				familyList = append(familyList, tempFamilyList)
				continue
			}

			return nil, ErrorMongoTechnicalFail
		}

		tempFamilyList.Author = spiderInfoList[0].Author
		tempFamilyList.Quantity = int32(len(spiderInfoList))

		familyList = append(familyList, tempFamilyList)
	}

	return familyList, nil
}
