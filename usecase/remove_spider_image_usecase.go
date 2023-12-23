package usecase

import (
	"context"
	"fmt"
	"os"
	"path"
	"spider-go/domain"
	"spider-go/logger"
	"spider-go/repository"
)

type RemoveSpiderImageUsecase struct {
	spiderRepo    domain.SpiderRepository
	fileImagePath string
	log           *logger.Logger
}

var (
	ErrorRemoveSpiderImageSpiderUUIDNotFound = fmt.Errorf("spider info not found")
)

func NewRemoveSpiderImageUsecase(spiderRepo domain.SpiderRepository, fileImagePath string) domain.RemoveSpiderImageUsecase {
	return &RemoveSpiderImageUsecase{
		spiderRepo:    spiderRepo,
		fileImagePath: fileImagePath,
		log:           logger.L().Named("RemoveSpiderImageUsecase"),
	}
}

func (u *RemoveSpiderImageUsecase) RemoveSpiderImageBySpiderImageNameList(ctx context.Context, spiderUUID string, spiderImageListRM []string) error {
	log := u.log.WithContext(ctx)

	spiderInfo, err := u.spiderRepo.FindSpiderByUUID(ctx, spiderUUID)
	if err != nil {

		log.Errorf("[RemoveSpiderImageBySpiderImageNameList] find spider info error: %+v", err)
		if err == repository.ErrorMongoNotFound {
			return ErrorRemoveSpiderImageSpiderUUIDNotFound
		}

		return ErrorMongoTechnicalFail
	}

	removeImageList := make(map[string]bool)
	newSpiderImageList := []string{}
	spiderImageList := spiderInfo.ImageFile

	for _, name := range spiderImageListRM {
		removeImageList[name] = true
	}

	for _, spiderImage := range spiderImageList {
		if _, ok := removeImageList[spiderImage]; ok {
			continue
		}
		newSpiderImageList = append(newSpiderImageList, spiderImage)
	}

	if err := u.spiderRepo.UpdateImageFileToSpiderInfo(ctx, newSpiderImageList, spiderUUID); err != nil {
		log.Errorf("[RemoveSpiderImageBySpiderImageNameList] update spider image failed, error: %+v", err)
		return ErrorMongoTechnicalFail
	}

	go func() {
		u.removeSpiderImage(ctx, spiderImageListRM)
	}()

	return nil

}

func (u *RemoveSpiderImageUsecase) removeSpiderImage(ctx context.Context, spiderImageList []string) {
	log := u.log.WithContext(ctx)

	for _, spiderImage := range spiderImageList {
		filepath := path.Join(u.fileImagePath, spiderImage)
		tryToRemove := true
		tryCount := 0

		var RemoveErr error

		for tryToRemove && tryCount < 3 {
			log.Debugf("[RemoveSpiderImageBySpiderImageNameList] count `%v`, try to remove spider image name is `%v` failed, error: %v", tryCount, spiderImage, RemoveErr)

			if err := os.Remove(filepath); err != nil {
				RemoveErr = err
				tryCount += 1
			} else {
				tryToRemove = false
			}
		}
		if tryToRemove {
			log.Errorf("[RemoveSpiderImageBySpiderImageNameList] remove spider image name is `%v` failed, error: %v", spiderImage, RemoveErr)
		} else {
			log.Infof("[RemoveSpiderImageBySpiderImageNameList]remove spider image name is `%v` successfull", spiderImage)
		}

	}

	log.Infof("[RemoveSpiderImageBySpiderImageNameList] removed spider image list successfull")
}
