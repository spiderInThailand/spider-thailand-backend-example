package handler

import (
	"net/http"
	api_model "spider-go/api/model"
	"spider-go/asset"
	"spider-go/domain"
	"spider-go/logger"
	"spider-go/usecase"

	"github.com/gin-gonic/gin"
)

type RsaHandler struct {
	GetRsaKeyUsecase domain.GetRsaKeyUsecase
	log              *logger.Logger
}

func NewRsaHandler(GetRsaKeyUsecase domain.GetRsaKeyUsecase) *RsaHandler {
	return &RsaHandler{
		GetRsaKeyUsecase: GetRsaKeyUsecase,
		log:              logger.L().Named("RsaHandler"),
	}
}

func (h *RsaHandler) RequestRsaPublicKey(ctx *gin.Context) {

	log := h.log.WithContext(ctx)

	var req api_model.GetRsaKeyRequest
	var resp api_model.GetRsaKeyResponse

	if err := ctx.ShouldBind(&req); err != nil {
		log.Errorf("should bind request failed: %+v", err)
		resp.Header.ErrorCode = asset.E().GeneralSystemError.ErrorCode
		resp.Header.Message = asset.E().GeneralSystemError.ErrorMessageEN
		ctx.AbortWithStatusJSON(http.StatusBadRequest, resp)
		return
	}

	publicKey, searchKey, err := h.GetRsaKeyUsecase.GenerateRsaKey(ctx)
	if err != nil {
		log.Errorf("usecase request failed: %+v", err)
		assetError := h.mappingRequestRsaPublicKey(err)
		resp.Header.ErrorCode = assetError.ErrorCode
		resp.Header.Message = assetError.ErrorMessageEN
		ctx.JSON(assetError.StatusCode, resp)
		return
	}

	resp.Data.PublicKey = publicKey
	resp.Data.SearchKey = searchKey

	resp.Header.ErrorCode = SUCCESS_CODE
	resp.Header.Message = SUCCESS_MESSAGE

	ctx.JSON(http.StatusOK, resp)

}

func (h *RsaHandler) mappingRequestRsaPublicKey(err error) *asset.ErrorCode {
	switch err {
	case usecase.ErrorGenerateRSAKey:
		return &asset.E().GenerateRSAError
	case usecase.ErrorSaveDataToRedis:
		return &asset.E().ErrorTempDB
	default:
		return &asset.E().GeneralSystemError
	}
}
