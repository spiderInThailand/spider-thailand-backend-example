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

type RegisterHandler struct {
	registerSpiderUsecase domain.RegisterSpiderUsecase
	log                   *logger.Logger
}

func NewRegisterHandler(registerSpiderUsecase domain.RegisterSpiderUsecase) *RegisterHandler {
	return &RegisterHandler{
		registerSpiderUsecase: registerSpiderUsecase,
		log:                   logger.L().Named("RegisterHandler"),
	}
}

func (h *RegisterHandler) RegisterHandler(ctx *gin.Context) {
	log := h.log.WithContext(ctx)

	var req api_model.RegisterSpiderInfoRequester
	var resp api_model.RegisterSpiderInfoResponser

	if err := ctx.ShouldBind(&req); err != nil {
		log.Errorf("should bind request failed: %+v", err)
		resp.Header.ErrorCode = asset.E().GeneralSystemError.ErrorCode
		resp.Header.Message = asset.E().GeneralSystemError.ErrorMessageEN
		ctx.AbortWithStatusJSON(http.StatusBadRequest, resp)
		return
	}

	log.Infof("register statistics info handler start with req: %+v", req)

	if err := validator.Struct(req); err != nil {
		log.Errorf("validate request data fail, error: %+v", err)
		resp.Header.ErrorCode = asset.E().RequestDataFail.ErrorCode
		resp.Header.Message = asset.E().RequestDataFail.ErrorMessageEN
		ctx.JSON(asset.E().RequestDataFail.StatusCode, resp)
		return
	}

	spiderUUID, err := h.registerSpiderUsecase.Register(ctx, req.Data, req.Header.Username)
	if err != nil {
		assetError := h.mapRegisterStatisticsInfoHandlerError(err)
		resp.Header.ErrorCode = assetError.ErrorCode
		resp.Header.Message = assetError.ErrorMessageEN
		ctx.JSON(assetError.StatusCode, resp)
		return

	}

	resp.Data = req.Data
	resp.Data.SpiderUUID = spiderUUID
	resp.Header.ErrorCode = SUCCESS_CODE
	resp.Header.Message = SUCCESS_MESSAGE

	ctx.JSON(http.StatusOK, resp)

}

func (h *RegisterHandler) mapRegisterStatisticsInfoHandlerError(err error) *asset.ErrorCode {
	switch err {
	case usecase.ErrorRedisConnection:
		return &asset.E().ErrorTempDB
	case usecase.ErrorMongoTechnicalFail, usecase.ErrorMongoConnection:
		return &asset.E().ErrorSpiderDB
	case usecase.ErrorValidateUserFaild:
		return &asset.E().UserNotLogin
	default:
		return &asset.E().GeneralSystemError
	}
}
