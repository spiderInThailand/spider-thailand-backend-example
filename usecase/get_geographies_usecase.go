package usecase

import (
	"context"
	"fmt"
	"spider-go/domain"
	"spider-go/logger"
	"spider-go/model"
	"spider-go/repository"
)

type ThaiGeographiesUsecase struct {
	thaiGeographiesRepo domain.ThaiGeographiesRepository
	SpiderRepo          domain.SpiderRepository
	log                 *logger.Logger
}

var (
	ErrorThaiGeographiesUsecaseZeroProvince              = fmt.Errorf("not have province in mongo")
	ErrorThaiGeographiesUsecaseDistrictNotFound          = fmt.Errorf("district not found")
	ErrorThaiGeographiesUsecaseValidateDataRequestFailed = fmt.Errorf("data request failed")
	ErrorThaiGeographiesUsecaseSpiderNotFound            = fmt.Errorf("spider not found")
)

func NewThaiGeographiesUsecase(thaiGeographiesRepo domain.ThaiGeographiesRepository, SpiderRepo domain.SpiderRepository) domain.ThaiGeographiesUsecase {
	return &ThaiGeographiesUsecase{
		thaiGeographiesRepo: thaiGeographiesRepo,
		SpiderRepo:          SpiderRepo,
		log:                 logger.L().Named("ThaiGeographiesUsecase"),
	}
}

// ==================================================================
// get normal geographies (province and district)
// ==================================================================
func (u *ThaiGeographiesUsecase) GetAllProvince(ctx context.Context) ([]model.Province, error) {
	log := u.log.WithContext(ctx)

	provinceList, err := u.thaiGeographiesRepo.GetAllProvince(ctx)
	if err != nil {
		log.Errorf("[GetAllProvince] get all province failed, error: %+v", err)
		return []model.Province{}, ErrorMongoTechnicalFail
	}

	if len(provinceList) == 0 {
		return []model.Province{}, ErrorThaiGeographiesUsecaseZeroProvince
	}

	return provinceList, nil
}

func (u *ThaiGeographiesUsecase) GetDistictWithProvinceNameEN(ctx context.Context, provinceNameEN string) ([]model.District, error) {
	log := u.log.WithContext(ctx)

	province, err := u.thaiGeographiesRepo.FindProvinceWithProvinceNameEN(ctx, provinceNameEN)

	if err != nil {
		log.Errorf("[GetDistictWithProvinceNameEN] find province with province name en failed, error: %+v", err)
		if err == repository.ErrorMongoNotFound {
			return []model.District{}, ErrorThaiGeographiesUsecaseDistrictNotFound
		}
		return []model.District{}, ErrorMongoTechnicalFail
	}

	return province.Amphure, nil
}

// ******************************************************************

func (u *ThaiGeographiesUsecase) GetGeographiesBySpiderType(ctx context.Context, family, genus, species string) ([]model.LocationResult, error) {
	log := u.log.WithContext(ctx)

	if isPass := u.validateSpiderTypeReques(family, genus, species); !isPass {
		log.Errorf("[GetGeographiesBySpiderType] validate data request failed, family: `%v`, genus: `%v`, species: `%v`", family, genus, species)
		return []model.LocationResult{}, ErrorThaiGeographiesUsecaseValidateDataRequestFailed
	}

	spiderInfoResults, err := u.SpiderRepo.FindSpiderInfoBySpiderType(ctx, family, genus, species, false, 0, 0)
	if err != nil {
		log.Errorf("[GetGeographiesBySpiderType] find spider request failed, error: %+v", err)
		if err.Error() == repository.ErrorMongoNotFound.Error() {
			return []model.LocationResult{}, ErrorThaiGeographiesUsecaseSpiderNotFound
		}
		return []model.LocationResult{}, ErrorMongoTechnicalFail
	}

	localtionResult := u.getLocationOfSpiderResult(spiderInfoResults)

	return localtionResult, nil

}

func (u *ThaiGeographiesUsecase) getLocationOfSpiderResult(spiderResult []model.SpiderInfo) []model.LocationResult {

	var allAddressList []model.Address
	var localtionResultList []model.LocationResult

	// 1. get all address from spider info
	for _, spiderInfo := range spiderResult {
		allAddressList = append(allAddressList, spiderInfo.Address...)
	}

	// 2. prepare valiable for map name of province and locality
	// NOTE: save name of province and locality to map[string]someting for easy to search and not used loop
	provinceMap := make(map[string]model.LocationResult)
	localityNameMap := make(map[string]int)

	for _, addressInfo := range allAddressList {
		// 3. validate province that have already to done
		if _, ok := provinceMap[addressInfo.Province]; ok {
			// 3.1 if province is done, validate locality that have already to done
			localityName := addressInfo.Province + addressInfo.Locality

			tempProvinceMap := provinceMap[addressInfo.Province]

			// 3.1.1 if locality is done, add unique position
			if _, ok := localityNameMap[localityName]; ok {

				tempPosition := make(map[string]bool)

				localityIndex := localityNameMap[localityName]
				for _, position := range tempProvinceMap.Locality[localityIndex].SubLocaltion {
					tempPosition[position] = true
				}

				for _, position := range addressInfo.Position {
					if _, ok := tempPosition[position.Name]; !ok {
						tempProvinceMap.Locality[localityIndex].SubLocaltion = append(tempProvinceMap.Locality[localityIndex].SubLocaltion, position.Name)
					}
				}

				provinceMap[addressInfo.Province] = tempProvinceMap

			} else {
				// 3.1.2 if locality is not done, add new locality
				var positionName []string
				for _, position := range addressInfo.Position {
					positionName = append(positionName, position.Name)
				}

				tempLocality := model.LocalityResult{
					Name:         addressInfo.Locality,
					SubLocaltion: positionName,
				}

				tempProvinceMap.Locality = append(provinceMap[addressInfo.Province].Locality, tempLocality)
				provinceMap[addressInfo.Province] = tempProvinceMap
			}

		} else {
			// 3.1 if province is not done, add new province
			var tempLocationResult model.LocationResult

			var positionName []string

			for _, position := range addressInfo.Position {
				positionName = append(positionName, position.Name)
			}
			tempLocality := model.LocalityResult{
				Name:         addressInfo.Locality,
				SubLocaltion: positionName,
			}

			tempLocationResult.Province = addressInfo.Province
			tempLocationResult.Locality = []model.LocalityResult{
				tempLocality,
			}

			localityName := addressInfo.Province + addressInfo.Locality

			localityNameMap[localityName] = 0
			provinceMap[addressInfo.Province] = tempLocationResult

		}
	}

	for _, locationResult := range provinceMap {
		localtionResultList = append(localtionResultList, locationResult)
	}

	return localtionResultList
}

func (u *ThaiGeographiesUsecase) validateSpiderTypeReques(family, genus, species string) bool {
	switch {
	case species != "" && genus == "":
		return false
	case genus != "" && family == "":
		return false
	case family == "":
		return false
	default:
		return true
	}
}
