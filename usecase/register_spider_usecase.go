package usecase

import (
	"context"
	"fmt"
	api_model "spider-go/api/model"
	"spider-go/domain"
	"spider-go/logger"
	"spider-go/model"
	"spider-go/utils/uuid"
	"time"

	"golang.org/x/exp/slices"
)

const SPIDER_PREFIX = "SPIDER_%s"

type RegisterSpiderUsecase struct {
	spiderRepo     domain.SpiderRepository
	statisticsRepo domain.StatisticsRepository
	log            *logger.Logger
}

func NewRegisterSpiderUsecase(spiderRepo domain.SpiderRepository,
	statisticsRepo domain.StatisticsRepository,
) domain.RegisterSpiderUsecase {
	return &RegisterSpiderUsecase{
		spiderRepo:     spiderRepo,
		statisticsRepo: statisticsRepo,
		log:            logger.L().Named("RegisterSpiderUsecase"),
	}

}

func (u *RegisterSpiderUsecase) Register(ctx context.Context, req api_model.SpiderInfo, username string) (string, error) {
	log := u.log.WithContext(ctx)

	log.Infof("start regiter spider")

	spiderInfo := u.prepareSpiderInfoFromRequest(req, username)

	// =======================================================
	// check family, genus, species is exist
	// =======================================================

	err := u.validateSpiderStatisticsIsFound(ctx, spiderInfo)
	if err != nil {
		return "", err
	}

	// =======================================================

	// =======================================================
	// save data to mongo
	// =======================================================

	if err := u.spiderRepo.InsertNewSpider(ctx, spiderInfo); err != nil {
		log.Errorf("[redister spider usercase] insert spider info failed, error: %+v", err)
		return "", ErrorMongoTechnicalFail
	}

	// =======================================================

	return spiderInfo.SpiderUUID, nil

}

func (u *RegisterSpiderUsecase) validateSpiderStatisticsIsFound(ctx context.Context, newSpiderInfo model.SpiderInfo) error {
	log := u.log.WithContext(ctx)

	spiderStatistics, err := u.statisticsRepo.FindSpiderStatisticsByFamily(ctx, newSpiderInfo.Family)
	if err != nil {
		// if mongo not found than insert new spider statistics
		log.Errorf("[validate spider statistics]  find spider statistics by family `%v`, error: %+v", newSpiderInfo.Family, err)
		if err.Error() == MONGO_NOT_FOUND {
			// upsert new spider statistics
			log.Infof("[validate spider statistics] insert new spider statistics, %v", newSpiderInfo.Family)
			newSpiderStatistics := u.prepareNewSpiderStatistics(newSpiderInfo)
			if err := u.statisticsRepo.UpsertSpiderStatistics(ctx, newSpiderInfo.Family, newSpiderStatistics); err != nil {
				log.Errorf("[validate spider statistics] upsert new spider statistics is failed, error: %+v", err)
				return ErrorMongoConnection
			}
			return nil
		} else {
			log.Errorf("[validate spider statistics] mongo error: %+v", err)
			return ErrorMongoConnection
		}
	}

	// find genus
	genusIndex := slices.IndexFunc(spiderStatistics.Genus, func(s model.GenusGroup) bool {
		return s.GenusName == newSpiderInfo.Genus
	})

	if genusIndex == -1 { // if genus index not found
		newGenus := model.GenusGroup{
			GenusName: newSpiderInfo.Genus,
			Species: []model.SpeciesGroup{
				{
					SpeciesName: newSpiderInfo.Species,
				},
			},
		}
		spiderStatistics.Genus = append(spiderStatistics.Genus, newGenus)

		if err := u.statisticsRepo.UpsertSpiderStatistics(ctx, newSpiderInfo.Family, *spiderStatistics); err != nil {
			log.Errorf("[validate spider statistics] upsert new spider statistics is failed, error: %+v", err)
			return ErrorMongoConnection
		}

		log.Infof("[validate spider statistics] upsert new statistics (genus and species) success")

		return nil
	}

	// find species
	speciesIndex := slices.IndexFunc(spiderStatistics.Genus[genusIndex].Species, func(s model.SpeciesGroup) bool {
		return s.SpeciesName == newSpiderInfo.Species
	})

	if speciesIndex == -1 { // if species index not found
		newSpecies := model.SpeciesGroup{
			SpeciesName: newSpiderInfo.Species,
		}
		spiderStatistics.Genus[genusIndex].Species = append(spiderStatistics.Genus[genusIndex].Species, newSpecies)

		if err := u.statisticsRepo.UpsertSpiderStatistics(ctx, newSpiderInfo.Family, *spiderStatistics); err != nil {
			log.Errorf("[validate spider statistics] upsert new spider statistics is failed, error: %+v", err)
			return ErrorMongoConnection
		}

		log.Infof("[validate spider statistics] upsert new statistics (species) success")

		return nil
	}

	return nil

}

func (u *RegisterSpiderUsecase) prepareSpiderInfoFromRequest(req api_model.SpiderInfo, username string) model.SpiderInfo {

	// generate uuid from spider_uuid

	spiderUUID := uuid.GernerateUUID32()

	var Address []model.Address
	timeNow := time.Now()

	for _, address := range req.Address {
		tempAddress := model.Address{
			Province: address.Province,
			District: address.District,
			Locality: address.Locality,
		}

		var tempPosition []model.Position

		for _, position := range address.Position {
			thisPosition := model.Position{
				Latitude:  position.Latitude,
				Longitude: position.Longitude,
				Name:      position.Name,
			}

			tempPosition = append(tempPosition, thisPosition)
		}

		tempAddress.Position = tempPosition

		Address = append(Address, tempAddress)
	}

	spiderInfo := model.SpiderInfo{
		SpiderUUID:   fmt.Sprintf(SPIDER_PREFIX, spiderUUID),
		Family:       req.Family,
		Genus:        req.Genus,
		Species:      req.Species,
		Author:       req.Author,
		PublishYear:  req.PublishYear,
		Country:      req.Country,
		CountryOther: req.OtherCountries,
		Altitude:     req.Altitude,
		Method:       req.Method,
		Habital:      req.Habital,
		Microhabital: req.Microhabital,
		Designate:    req.Designate,
		Address:      Address,
		Paper:        req.Paper,
		CreatedAt:    timeNow,
		UpdatedAt:    timeNow,
		Status:       model.SPIDER_INFO_STATUS_ACTIVE,
		CreatedBy:    username,
	}

	return spiderInfo
}

func (u *RegisterSpiderUsecase) prepareNewSpiderStatistics(newSpiderInfo model.SpiderInfo) model.SpiderStatistics {
	data := model.SpiderStatistics{
		FamilyName: newSpiderInfo.Family,
		Genus: []model.GenusGroup{
			{
				GenusName: newSpiderInfo.Genus,
				Species: []model.SpeciesGroup{
					{
						SpeciesName: newSpiderInfo.Species,
					},
				},
			},
		},
		CreatedAt: newSpiderInfo.CreatedAt,
		UpdatedAt: newSpiderInfo.CreatedAt,
	}

	return data
}
