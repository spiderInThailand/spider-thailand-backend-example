package usecase

import (
	"context"
	"fmt"
	api_model "spider-go/api/model"
	"spider-go/domain"
	"spider-go/logger"
	"spider-go/model"
	"time"
)

type UpdateSpiderInfoUsecase struct {
	spiderRepo domain.SpiderRepository
	log        *logger.Logger
}

var (
	ErrorUpdateSpiderInfoUsecaseSpiderUUIDNotFound = fmt.Errorf("spider uuid is not found in mongodb")
)

func NewUpdateSpiderInfoUsecase(spiderRepo domain.SpiderRepository) domain.UpdateSpiderInfoUsecase {
	return &UpdateSpiderInfoUsecase{
		spiderRepo: spiderRepo,
		log:        logger.L().Named("UpdateSpiderInfoUsecase"),
	}
}

func (u *UpdateSpiderInfoUsecase) UpdateSpiderInfoUsecase(ctx context.Context, spiderInfoReq api_model.SpiderInfo) error {
	log := u.log.WithContext(ctx)

	spiderUUID := spiderInfoReq.SpiderUUID
	spiderInfo := u.prepareSpiderInfoFromRequest(spiderInfoReq)

	log.Infof("[UpdateSpiderInfoUsecase] update spider with spider uuid: %v", spiderUUID)

	isUpdate, err := u.spiderRepo.UpdateSpiderInfo(ctx, spiderUUID, spiderInfo)
	if err != nil {
		return ErrorMongoTechnicalFail
	}

	if !isUpdate {
		return ErrorUpdateSpiderInfoUsecaseSpiderUUIDNotFound
	}

	return nil

}

func (u *UpdateSpiderInfoUsecase) prepareSpiderInfoFromRequest(req api_model.SpiderInfo) model.SpiderInfo {

	// generate uuid from spider_uuid

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
		SpiderUUID:   req.SpiderUUID,
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
		UpdatedAt:    timeNow,
	}

	return spiderInfo
}
