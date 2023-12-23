package usecase

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"image/jpeg"
	"image/png"
	"os"
	"path"
	"sort"
	"spider-go/config"
	"spider-go/domain"
	"spider-go/logger"
	"spider-go/utils/uuid"
	"strings"
)

type UploadImageUsecase struct {
	spiderRepo domain.SpiderRepository
	log        *logger.Logger
}

var (
	ErrorUploadImageUsecaseVlidateSpiderUUID = fmt.Errorf("validate spider uuid failed not found")
	ErrorUploadImageUsecaseSaveImageFileFail = fmt.Errorf("save spider image file failed")
	ErrorUploadImageUsecaseFileTypeNotMatch  = fmt.Errorf("file type of image not match")
)

func NewUploadImageUsecase(spiderRepo domain.SpiderRepository) domain.UploadImageUsecase {
	return &UploadImageUsecase{
		spiderRepo: spiderRepo,
		log:        logger.L().Named("UploadImageUsecase"),
	}
}

func (u *UploadImageUsecase) UploadImageSpiderUsecase(ctx context.Context, spiderUUID string, listImageEncode64 []string) error {
	log := u.log.WithContext(ctx)

	// =======================================================
	// validate spider uuid
	// =======================================================
	spiderInfo, err := u.spiderRepo.FindSpiderByUUID(ctx, spiderUUID)
	if err != nil && err.Error() != MONGO_NOT_FOUND {
		return ErrorUploadImageUsecaseVlidateSpiderUUID
	}

	// =======================================================
	// save file
	// =======================================================
	listImageName, err := u.handleFileImage(ctx, spiderUUID, listImageEncode64)
	if err != nil {
		return err
	}

	sort.Strings(listImageName)

	// =======================================================
	// add file name to mongo
	// =======================================================

	spiderInfo.ImageFile = append(spiderInfo.ImageFile, listImageName...)

	if err := u.spiderRepo.UpdateImageFileToSpiderInfo(ctx, spiderInfo.ImageFile, spiderUUID); err != nil {
		log.Errorf("[UploadImageSpiderUsecase] update spider info mongo failed, error: %v", err)
		u.deleteFile(ctx, listImageName)
		return ErrorMongoTechnicalFail
	}

	return nil
}

func (u *UploadImageUsecase) handleFileImage(ctx context.Context, spiderUUID string, listImageEncode64 []string) ([]string, error) {
	log := u.log.WithContext(ctx)

	var listImageName []string

	for _, image := range listImageEncode64 {

		// structure image base64
		// "data:[<mediatype>][;base64],<data>
		commaIndex := strings.Index(string(image), ",")
		OriginalImage := string(image)[commaIndex+1:]

		spiderImageUUID := uuid.GernerateUUID32()

		fileName := fmt.Sprintf(config.C().File.SpiderImage, spiderUUID, spiderImageUUID)
		filePath := path.Join(config.C().File.FileImagePath, fileName)

		log.Infof("[handleFileImage] filename: %v", fileName)

		imageDecode, err := base64.StdEncoding.DecodeString(OriginalImage)
		if err != nil {
			log.Errorf("[handleFileImage] base64 decode original image error: %v", err)
			fmt.Printf("point: 1")
			return nil, ErrorTechnicalError
		}

		renderImageDecode := bytes.NewReader(imageDecode)

		imageType := strings.TrimSuffix(string(image[5:commaIndex]), ";base64")
		// imageType := JPEG_IMAGE_TYPE

		log.Infof("[handleFileImage] image type: %v", imageType)

		switch imageType {
		case PNG_IMAGE_TYPE:
			// save png
			err := u.savePNGFile(ctx, renderImageDecode, filePath)
			if err != nil {
				return nil, err
			}

			fileNameType := fmt.Sprintf("%s.png", fileName)

			listImageName = append(listImageName, fileNameType)

		case JPEG_IMAGE_TYPE:
			// save jpeg
			err := u.saveJPEGFile(ctx, renderImageDecode, filePath)
			if err != nil {
				return nil, err
			}

			fileNameType := fmt.Sprintf("%s.jpeg", fileName)

			listImageName = append(listImageName, fileNameType)

		default:
			return nil, ErrorUploadImageUsecaseFileTypeNotMatch
		}

	}

	return listImageName, nil
}

func (u *UploadImageUsecase) savePNGFile(ctx context.Context, renderImageDecode *bytes.Reader, filePath string) error {
	log := u.log.WithContext(ctx)

	pngImage, err := png.Decode(renderImageDecode)
	if err != nil {
		log.Errorf("[savePNGFile] png decode error: %v", err)
		fmt.Printf("point: 2")
		return ErrorTechnicalError
	}

	file, err := os.OpenFile(filePath+".png", os.O_WRONLY|os.O_CREATE, 0777)
	if err != nil {
		log.Errorf("[savePNGFile] open file error: %v", err)
		fmt.Printf("point: 3, error: %v", err)
		return ErrorTechnicalError
	}

	defer file.Close()

	png.Encode(file, pngImage)

	return nil
}

func (u *UploadImageUsecase) saveJPEGFile(ctx context.Context, renderImageDecode *bytes.Reader, filePath string) error {
	log := u.log.WithContext(ctx)

	jpegImage, err := jpeg.Decode(renderImageDecode)
	if err != nil {
		log.Errorf("[saveJPEGFile] jpeg decode error: %v", err)
		fmt.Printf("point: 4")
		return ErrorTechnicalError
	}

	file, err := os.OpenFile(filePath+".jpeg", os.O_WRONLY|os.O_CREATE, 0777)
	if err != nil {
		log.Errorf("[saveJPEGFile] open file error: %v", err)
		fmt.Printf("point: 5")
		return ErrorTechnicalError
	}

	defer file.Close()

	jpegOpts := jpeg.Options{Quality: 100}

	jpeg.Encode(file, jpegImage, &jpegOpts)

	return nil
}

func (u *UploadImageUsecase) deleteFile(ctx context.Context, files []string) {
	log := u.log.WithContext(ctx)

	for _, file := range files {

		filePath := fmt.Sprintf(config.C().File.FileImagePath, file)

		if err := os.Remove(filePath); err != nil {
			log.Warnf("[deleteFile] remove file: `%s` failed, error: %v", file, err)
		}
	}
}
