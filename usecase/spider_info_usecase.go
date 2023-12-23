package usecase

import (
	"context"
	"encoding/base64"
	"fmt"
	"net/http"
	"os"
	"path"
	"spider-go/domain"
	"spider-go/logger"
	"spider-go/model"
	"spider-go/repository"
)

type SpiderInfoUsecase struct {
	spiderRepo domain.SpiderRepository
	accRepo    domain.AccountRepository
	// separated from config to prevent unit tests from generating data races
	fileImagePath string
	log           *logger.Logger
}

var (
	ErrorSpiderInfoUsecaseSpiderNotFound                 = fmt.Errorf("[spider info usecase] spider info this uuid not found")
	ErrorSpiderInfoUsecaseAccountInsufficientPermissions = fmt.Errorf("[spider info usecase] this user account is insufficient permissions")
	ErrorSpiderInfoUsecaseReadFileFail                   = fmt.Errorf("[spider info usecase] read file fail")
	ErrorSpiderInfoUsecaseValidateDataFail               = fmt.Errorf("invalid data request")
)

func NewSpiderInfoUsecase(
	spiderRepo domain.SpiderRepository,
	accRepo domain.AccountRepository,
	fileImagePath string,
) domain.SpiderInfoUsecase {
	return &SpiderInfoUsecase{
		spiderRepo:    spiderRepo,
		accRepo:       accRepo,
		fileImagePath: fileImagePath,
		log:           logger.L().Named("SpiderInfoUsecase"),
	}
}

// ========================================================
// get spider info from mongodb
// ========================================================

func (u *SpiderInfoUsecase) GetSpiderInfoUsecase(ctx context.Context, spiderUUID, username string) (*model.SpiderInfo, error) {
	log := u.log.WithContext(ctx)

	account, err := u.accRepo.FindAccountByUsername(ctx, username)
	if err != nil && err != repository.ErrorMongoNotFound {
		log.Errorf("[get spider info] find account by username error: %+v", err)
		return &model.SpiderInfo{}, ErrorMongoTechnicalFail
	}

	SpiderInfo, err := u.spiderRepo.FindSpiderByUUID(ctx, spiderUUID)
	if err != nil {
		log.Errorf("[get spider info] error mongo, error: %v", err)
		if err == repository.ErrorMongoNotFound {
			return &model.SpiderInfo{}, ErrorSpiderInfoUsecaseSpiderNotFound
		}
		return &model.SpiderInfo{}, ErrorMongoConnection
	}

	if SpiderInfo.Status != model.SPIDER_INFO_STATUS_ACTIVE && account.Role != model.ACCOUNT_ROLE_ADMIN {
		log.Errorf("user permissions denied")
		return &model.SpiderInfo{}, ErrorSpiderInfoUsecaseAccountInsufficientPermissions
	}

	return SpiderInfo, nil

}

// ********************************************************

// ========================================================
// get multi image
// ========================================================

func (u *SpiderInfoUsecase) GetSpiderImagesUsecase(ctx context.Context, fileImages []string) ([]model.SpiderImageList, error) {
	log := u.log.WithContext(ctx)

	var imageEncodeFiles []model.SpiderImageList

	for _, imageName := range fileImages {

		pathFile := path.Join(u.fileImagePath, imageName)

		log.Debugf("path file image is `%v`", pathFile)

		bytes, err := os.ReadFile(pathFile)
		if err != nil {
			log.Errorf("[GetSpiderImagesUsecase] read file at path `%v` failed, error: %v", pathFile, err)
			return nil, ErrorSpiderInfoUsecaseReadFileFail
		}

		var imageEncode string

		mimeType := http.DetectContentType(bytes)

		switch mimeType {
		case JPEG_IMAGE_TYPE:
			imageEncode += "data:image/jpeg;base64,"
		case PNG_IMAGE_TYPE:
			imageEncode += "data:image/png;base64,"
		}

		imageEncode += base64.StdEncoding.EncodeToString(bytes)

		thisSpiderImage := model.SpiderImageList{
			Name:        imageName,
			ImageBase64: imageEncode,
		}

		imageEncodeFiles = append(imageEncodeFiles, thisSpiderImage)

	}

	return imageEncodeFiles, nil

}

// ********************************************************

// ========================================================
// get spider info list manager
// ========================================================

func (u *SpiderInfoUsecase) GetSpiderInfoListManager(ctx context.Context, usecase string, page, limit int) ([]model.SpiderInfo, error) {
	log := u.log.WithContext(ctx)

	account, err := u.accRepo.FindAccountByUsername(ctx, usecase)
	if err != nil {
		log.Errorf("[GetSpiderInfoListManager] find account by username error: %+v", err)
		return []model.SpiderInfo{}, ErrorMongoTechnicalFail
	}

	if account.Role != model.ACCOUNT_ROLE_ADMIN {
		log.Errorf("[GetSpiderInfoListManager] user permissions denied")
		return []model.SpiderInfo{}, ErrorSpiderInfoUsecaseAccountInsufficientPermissions
	}

	spiderInfoList, err := u.spiderRepo.FindAllSpiderListManager(ctx, page, limit)
	if err != nil {
		log.Errorf("[GetSpiderInfoListManager] find spider list manager repo error: %+v", err)
		return []model.SpiderInfo{}, ErrorMongoConnection
	}

	log.Infof("[GetSpiderInfoListManager] length of spider info: %v", len(spiderInfoList))

	return spiderInfoList, nil
}

// ********************************************************

// ========================================================
// get spider info list filter by province and district
// ========================================================

func (u *SpiderInfoUsecase) GetSpiderInfoListByGeographies(ctx context.Context, province, district, position string) ([]model.SpiderInfo, error) {
	log := u.log.WithContext(ctx)

	// validate province
	isPass := u.validateGeographiesRequestData(province, district, position)
	if !isPass {
		log.Errorf("[GetSpiderInfoListByGeographies] validate data not pass, province: %v, district: %v, position: %v", province, district, position)
		return []model.SpiderInfo{}, ErrorSpiderInfoUsecaseValidateDataFail
	}

	spiderInfoList, err := u.spiderRepo.FindSpiderInfoListByGeographies(ctx, province, district, position)
	if err != nil {
		log.Errorf("[GetSpiderInfoListByGeographies] find spiderinfo repo failed: error: %+v", err)
		if err.Error() == repository.ErrorMongoNotFound.Error() {
			return []model.SpiderInfo{}, ErrorSpiderInfoUsecaseSpiderNotFound

		}
		return []model.SpiderInfo{}, ErrorMongoTechnicalFail
	}

	return spiderInfoList, nil
}

func (u *SpiderInfoUsecase) validateGeographiesRequestData(province, district, position string) bool {
	if district == "" && position == "" {
		return province != ""
	}

	return true
}

// ********************************************************

// ========================================================
// get spider info list filter by locality
// ========================================================
func (u *SpiderInfoUsecase) GetSpiderInfoListByLocality(ctx context.Context, locality string, page, size int32) ([]model.SpiderInfo, error) {
	log := u.log.WithContext(ctx)

	log.Infof("[GetSpiderInfoListByLocality] start usecase with param, %v, %v, %v", locality, page, size)

	spiderInfoList, err := u.spiderRepo.FindSpiderInfoByLocality(ctx, locality, page, size)
	if err != nil {
		log.Errorf("[GetSpiderInfoListByLocality] find spiderinfo repo failed: error: %+v", err)
		if err.Error() == repository.ErrorMongoNotFound.Error() {
			return []model.SpiderInfo{}, ErrorSpiderInfoUsecaseSpiderNotFound

		}
		return []model.SpiderInfo{}, ErrorMongoTechnicalFail
	}

	return spiderInfoList, nil
}

// ********************************************************

// ========================================================
// get spider list by spider type
// ========================================================
func (u *SpiderInfoUsecase) GetSpiderListBySpiderTypeUsecase(ctx context.Context, param model.GetSpiderListBySpiderTypeParam) ([]model.SpiderInfo, error) {
	log := u.log.WithContext(ctx)

	log.Infof("[GetSpiderListBySpiderTypeUsecase] start usecase with param, %+v", param)

	spiderInfoList, err := u.spiderRepo.FindSpiderInfoBySpiderType(ctx, param.Family, param.Genus, param.Species, true, param.Page, param.Size)
	if err != nil {
		log.Errorf("[GetSpiderInfoListByLocality] find spiderinfo repo failed: error: %+v", err)
		if err.Error() == repository.ErrorMongoNotFound.Error() {
			return []model.SpiderInfo{}, ErrorSpiderInfoUsecaseSpiderNotFound
		}
		return []model.SpiderInfo{}, ErrorMongoTechnicalFail
	}

	return spiderInfoList, nil
}
