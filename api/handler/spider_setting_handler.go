package handler

import (
	"net/http"
	api_model "spider-go/api/model"
	"spider-go/asset"
	"spider-go/domain"
	"spider-go/logger"
	"spider-go/usecase"
	"spider-go/utils/validator"

	"github.com/gin-gonic/gin"
)

type SpiderSettingHandler struct {
	uploadImageUsecase       domain.UploadImageUsecase
	spiderSettingUsecase     domain.DeleteSpiderInfoUsecase
	updateSpiderInfoUsecase  domain.UpdateSpiderInfoUsecase
	removeSpiderImageUsecase domain.RemoveSpiderImageUsecase
	log                      *logger.Logger
}

func NewSpiderSettingHandler(
	uploadImageUsecase domain.UploadImageUsecase,
	spiderSettingUsecase domain.DeleteSpiderInfoUsecase,
	updateSpiderInfoUsecase domain.UpdateSpiderInfoUsecase,
	removeSpiderImage domain.RemoveSpiderImageUsecase,
) *SpiderSettingHandler {
	return &SpiderSettingHandler{
		uploadImageUsecase:       uploadImageUsecase,
		spiderSettingUsecase:     spiderSettingUsecase,
		updateSpiderInfoUsecase:  updateSpiderInfoUsecase,
		removeSpiderImageUsecase: removeSpiderImage,
		log:                      logger.L().Named("SpiderSettingHandler"),
	}
}

// =========================================================
// upload image
// =========================================================
func (h *SpiderSettingHandler) UploadImageSpiderHandler(ctx *gin.Context) {
	log := h.log.WithContext(ctx)

	var req api_model.SpiderImageSettingRequester
	var resp api_model.SpiderImageSettingResponser

	// read info
	if err := ctx.ShouldBind(&req); err != nil {
		log.Errorf("should bind request failed: %+v", err)
		resp.Header.ErrorCode = asset.E().GeneralSystemError.ErrorCode
		resp.Header.Message = asset.E().GeneralSystemError.ErrorMessageEN
		ctx.AbortWithStatusJSON(asset.E().GeneralSystemError.StatusCode, resp)
		return
	}

	log.Infof("upload image spider handler start")

	if err := validator.Struct(req); err != nil {
		log.Errorf("validate request data fail, error: %+v", err)
		resp.Header.ErrorCode = asset.E().RequestDataFail.ErrorCode
		resp.Header.Message = asset.E().RequestDataFail.ErrorMessageEN
		ctx.JSON(asset.E().RequestDataFail.StatusCode, resp)
		return
	}

	err := h.uploadImageUsecase.UploadImageSpiderUsecase(ctx, req.Data.SpiderUUID, req.Data.ListImageEncode)
	if err != nil {
		log.Errorf("[UploadImageSpiderHandler] upload image usecase failed, error: %v", err)
		assetErr := h.mapUploadImageHandlerErrorCode(err)
		resp.Header.ErrorCode = assetErr.ErrorCode
		resp.Header.Message = assetErr.ErrorMessageEN
		ctx.JSON(assetErr.StatusCode, resp)
		return
	}

	resp.Header.ErrorCode = SUCCESS_CODE
	resp.Header.Message = ""
	ctx.JSON(http.StatusOK, resp)
}

func (h *SpiderSettingHandler) mapUploadImageHandlerErrorCode(err error) *asset.ErrorCode {
	switch err {
	case usecase.ErrorUploadImageUsecaseVlidateSpiderUUID:
		return &asset.E().SpiderNotFound
	case usecase.ErrorRedisConnection:
		return &asset.E().ErrorTempDB
	case usecase.ErrorValidateUserFaild:
		return &asset.E().UserNotLogin
	case usecase.ErrorUploadImageUsecaseFileTypeNotMatch:
		return &asset.E().InvalidImageType
	case usecase.ErrorMongoTechnicalFail:
		return &asset.E().ErrorSpiderDB
	default:
		return &asset.E().GeneralSystemError
	}
}

// *************************************************

// =========================================================
// delete spider info
// =========================================================
func (h *SpiderSettingHandler) DeleteSpiderHandler(ctx *gin.Context) {
	log := h.log.WithContext(ctx)

	var req api_model.DeleteSpiderRequester
	var resp api_model.DeleteSpiderResponser

	if err := ctx.ShouldBind(&req); err != nil {
		log.Errorf("should bind request failed: %+v", err)
		resp.Header.ErrorCode = asset.E().GeneralSystemError.ErrorCode
		resp.Header.Message = asset.E().GeneralSystemError.ErrorMessageEN
		ctx.AbortWithStatusJSON(asset.E().GeneralSystemError.StatusCode, resp)
		return
	}
	log.Infof("[DeleteSpiderHandler] start delete spider with req: %+v", req)

	if err := validator.Struct(req); err != nil {
		log.Errorf("validate request data fail, error: %+v", err)
		resp.Header.ErrorCode = asset.E().RequestDataFail.ErrorCode
		resp.Header.Message = asset.E().RequestDataFail.ErrorMessageEN
		ctx.JSON(asset.E().RequestDataFail.StatusCode, resp)
		return
	}

	err := h.spiderSettingUsecase.DeleteSpiderInfoUsecase(ctx, req.Data.SpiderUUID)
	if err != nil {
		log.Errorf("[DeleteSpiderHandler] delete spider info usecase error: %+v", err)
		assetErr := h.mapDeleteSpiderHandler(err)
		resp.Header.ErrorCode = assetErr.ErrorCode
		resp.Header.Message = assetErr.ErrorMessageEN
		ctx.AbortWithStatusJSON(assetErr.StatusCode, resp)
		return
	}

	resp.Header.ErrorCode = SUCCESS_CODE
	resp.Header.Message = ""

	ctx.JSON(http.StatusOK, resp)
}

func (h *SpiderSettingHandler) mapDeleteSpiderHandler(err error) *asset.ErrorCode {
	switch err {
	case usecase.ErrorMongoTechnicalFail:
		return &asset.E().SpiderNotFound
	case usecase.ErrorDeleteSpiderInfoUsecaseDeleteSpiderInfoFailed:
		return &asset.E().DeleteSpiderFailed
	default:
		return &asset.E().GeneralSystemError
	}
}

// *************************************************

// =========================================================
// edit spider info
// =========================================================
func (h *SpiderSettingHandler) EditSpiderInfoHandler(ctx *gin.Context) {
	log := h.log.WithContext(ctx)

	var req api_model.EditSpiderInfoRequester
	var resp api_model.EditSpiderInfoResponser

	if err := ctx.ShouldBind(&req); err != nil {
		log.Errorf("[EditSpiderInfoHandler] should bind request failed: %+v", err)
		resp.Header.ErrorCode = asset.E().GeneralSystemError.ErrorCode
		resp.Header.Message = asset.E().GeneralSystemError.ErrorMessageEN
		ctx.AbortWithStatusJSON(asset.E().GeneralSystemError.StatusCode, resp)
		return
	}
	log.Infof("[EditSpiderInfoHandler] start EditSpiderInfoHandler with request: %+v", req)

	if err := validator.Struct(req); err != nil {
		log.Errorf("validate request data fail, error: %+v", err)
		resp.Header.ErrorCode = asset.E().RequestDataFail.ErrorCode
		resp.Header.Message = asset.E().RequestDataFail.ErrorMessageEN
		ctx.JSON(asset.E().RequestDataFail.StatusCode, resp)
		return
	}

	err := h.updateSpiderInfoUsecase.UpdateSpiderInfoUsecase(ctx, req.Data)
	if err != nil {
		log.Errorf("[DeleteSpiderHandler] delete spider info usecase error: %+v", err)
		assetErr := h.mapEditSpiderInfoHandler(err)
		resp.Header.ErrorCode = assetErr.ErrorCode
		resp.Header.Message = assetErr.ErrorMessageEN
		ctx.AbortWithStatusJSON(assetErr.StatusCode, resp)
		return
	}

	resp.Header.ErrorCode = SUCCESS_CODE
	resp.Header.Message = ""

	ctx.JSON(http.StatusOK, resp)
}

func (h *SpiderSettingHandler) mapEditSpiderInfoHandler(err error) *asset.ErrorCode {
	switch err {
	case usecase.ErrorMongoTechnicalFail:
		return &asset.E().SpiderNotFound
	case usecase.ErrorUpdateSpiderInfoUsecaseSpiderUUIDNotFound:
		return &asset.E().SpiderNotFound
	default:
		return &asset.E().GeneralSystemError
	}
}

// *************************************************

// =========================================================
// remove spider image
// =========================================================
func (h *SpiderSettingHandler) RemoveSpiderImageHandler(ctx *gin.Context) {
	log := h.log.WithContext(ctx)

	var req api_model.RemoveSpiderImageRequester
	var resp api_model.RemoveSpiderImageResponse

	if err := ctx.ShouldBind(&req); err != nil {
		log.Errorf("[RemoveSpiderImageHandler] should bind request failed: %+v", err)
		resp.Header.ErrorCode = asset.E().GeneralSystemError.ErrorCode
		resp.Header.Message = asset.E().GeneralSystemError.ErrorMessageEN
		ctx.AbortWithStatusJSON(asset.E().GeneralSystemError.StatusCode, resp)
		return
	}

	log.Infof("[DeleteSpiderHandler] start delete spider with req: %+v", req)

	if err := validator.Struct(req); err != nil {
		log.Errorf("validate request data fail, error: %+v", err)
		resp.Header.ErrorCode = asset.E().RequestDataFail.ErrorCode
		resp.Header.Message = asset.E().RequestDataFail.ErrorMessageEN
		ctx.JSON(asset.E().RequestDataFail.StatusCode, resp)
		return
	}

	err := h.removeSpiderImageUsecase.RemoveSpiderImageBySpiderImageNameList(ctx, req.Data.SpiderUUID, req.Data.SpiderImageList)
	if err != nil {
		log.Errorf("[RemoveSpiderImageHandler] remove spider image usecase error: %+v", err)
		assetErr := h.mapRemoveSpiderImageHandler(err)
		resp.Header.ErrorCode = assetErr.ErrorCode
		resp.Header.Message = assetErr.ErrorMessageEN
		ctx.AbortWithStatusJSON(assetErr.StatusCode, resp)
		return
	}

	resp.Header.ErrorCode = SUCCESS_CODE
	resp.Header.Message = ""

	ctx.JSON(http.StatusOK, resp)
}

func (h *SpiderSettingHandler) mapRemoveSpiderImageHandler(err error) *asset.ErrorCode {
	switch err {
	case usecase.ErrorMongoTechnicalFail:
		return &asset.E().SpiderNotFound
	case usecase.ErrorUpdateSpiderInfoUsecaseSpiderUUIDNotFound:
		return &asset.E().SpiderNotFound
	default:
		return &asset.E().GeneralSystemError
	}
}
