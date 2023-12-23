package usecase

import (
	"context"
	"fmt"
	"os"
	"path"
	"spider-go/config"
	"spider-go/domain"
	"spider-go/logger"
)

var (
	ErrorDeleteSpiderInfoUsecaseRemoveFileFailed       = fmt.Errorf("remove spider image file failed")
	ErrorDeleteSpiderInfoUsecaseDeleteSpiderInfoFailed = fmt.Errorf("delete spider info in mongodb failed")
)

type DeleteSpiderInfoUsecase struct {
	spiderRepo domain.SpiderRepository
	log        *logger.Logger
}

func NewDeleteSpiderInfoUsecase(spiderRepo domain.SpiderRepository) domain.DeleteSpiderInfoUsecase {
	return &DeleteSpiderInfoUsecase{
		spiderRepo: spiderRepo,
		log:        logger.L().Named("DeleteSpiderInfoUsecase"),
	}
}

func (u *DeleteSpiderInfoUsecase) DeleteSpiderInfoUsecase(ctx context.Context, spiderUUID string) error {
	log := u.log.WithContext(ctx)

	log.Infof("[DeleteSpiderInfoUsecase] delete spider info at spider_uuid is `%v`", spiderUUID)

	// find spider info for get image name to delete
	spiderInfo, err := u.spiderRepo.FindSpiderByUUID(ctx, spiderUUID)
	if err != nil {
		log.Errorf("[DeleteSpiderInfoUsecase] find spider info error: %+v", err)
		return ErrorMongoTechnicalFail
	}

	if err := u.spiderRepo.DeleteSpiderInfoWithSpiderUUID(ctx, spiderUUID); err != nil {
		log.Errorf("[DeleteSpiderInfoUsecase] delete spider info at spider_uuid `%v` failed, error: %+v", spiderUUID, err)
		return ErrorDeleteSpiderInfoUsecaseDeleteSpiderInfoFailed
	}

	go func() {
		u.removeSpiderImage(ctx, spiderInfo.ImageFile)
	}()

	return nil
}

func (u *DeleteSpiderInfoUsecase) removeSpiderImage(ctx context.Context, spiderImageList []string) {
	log := u.log.WithContext(ctx)

	for _, spiderImage := range spiderImageList {
		filepath := path.Join(config.C().File.FileImagePath, spiderImage)
		tryToRemove := true
		tryCount := 0

		var RemoveErr error

		for tryToRemove && tryCount < 3 {
			log.Debugf("[DeleteSpiderInfoUsecase] count `%v`, try to remove spider image name is `%v` failed, error: %v", tryCount, spiderImage, RemoveErr)

			if err := os.Remove(filepath); err != nil {
				RemoveErr = err
				tryCount += 1
			} else {
				tryToRemove = false
			}
		}
		if tryToRemove {
			log.Errorf("[DeleteSpiderInfoUsecase] remove spider image name is `%v` failed, error: %v", spiderImage, RemoveErr)
		} else {
			log.Infof("[DeleteSpiderInfoUsecase]remove spider image name is `%v` successfull", spiderImage)
		}

	}

	log.Infof("[DeleteSpiderInfoUsecase] removed spider image list successfull")
}
