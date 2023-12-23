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

type GetSpiderStatisticsHandler struct {
	spiderStatisticsUsecase domain.StatisticsUsecase
	getFamilyListUsecase    domain.GetFamilyListUsecase
	log                     *logger.Logger
}

func NewGetSpiderStatisricsHandler(spiderStatisticsUsecase domain.StatisticsUsecase, getFamilyListUsecase domain.GetFamilyListUsecase) *GetSpiderStatisticsHandler {
	return &GetSpiderStatisticsHandler{
		spiderStatisticsUsecase: spiderStatisticsUsecase,
		getFamilyListUsecase:    getFamilyListUsecase,
		log:                     logger.L().Named("GetSpiderStatisticsHandler"),
	}
}

func (h *GetSpiderStatisticsHandler) GetSpiderStatisticsListHandler(ctx *gin.Context) {

	log := h.log.WithContext(ctx)

	var req api_model.SpiderStatisticsRequester
	var resp api_model.SpiderStatisticsResponser

	if err := ctx.ShouldBind(&req); err != nil {
		log.Errorf("should bind request failed: %+v", err)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, err)
		return
	}

	result, err := h.spiderStatisticsUsecase.GetSpiderStatisticsList(ctx)
	if err != nil {
		log.Errorf("usecase request failed: %+v", err)
		assetError := asset.E().ErrorSpiderDB
		resp.Header.ErrorCode = assetError.ErrorCode
		resp.Header.Message = assetError.ErrorMessageEN
		ctx.JSON(assetError.StatusCode, resp)
		return
	}

	resp.Header.ErrorCode = SUCCESS_CODE
	resp.Header.Message = SUCCESS_MESSAGE
	resp.Data = result

	ctx.JSON(200, resp)
}

func (h *GetSpiderStatisticsHandler) GetFamilyListhandler(ctx *gin.Context) {
	log := h.log.WithContext(ctx)

	var req api_model.GetFamilyListRequester
	var resp api_model.GetFamilyListResponser

	if err := ctx.ShouldBind(&req); err != nil {
		log.Errorf("should bind request failed: %+v", err)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, err)
		return
	}

	log.Infof("get family list handler start with req: %+v", req)

	if err := validator.Struct(req); err != nil {
		log.Errorf("validate request data fail, error: %+v", err)
		resp.Header.ErrorCode = asset.E().RequestDataFail.ErrorCode
		resp.Header.Message = asset.E().RequestDataFail.ErrorMessageEN
		ctx.JSON(asset.E().RequestDataFail.StatusCode, resp)
		return
	}

	familyList, err := h.getFamilyListUsecase.Execute(ctx, req.Data.Page, req.Data.Size)
	if err != nil {
		log.Errorf("[GetFamilyListhandler] usecase request failed: %+v", err)
		assetError := h.mapGetFamilyListhandlerError(err)
		resp.Header.ErrorCode = assetError.ErrorCode
		resp.Header.Message = assetError.ErrorMessageEN
		ctx.JSON(assetError.StatusCode, resp)
		return
	}

	resp.Header.ErrorCode = SUCCESS_CODE
	resp.Header.Message = SUCCESS_MESSAGE
	resp.Data.FamilyList = familyList

	ctx.JSON(http.StatusOK, resp)

}

func (h *GetSpiderStatisticsHandler) mapGetFamilyListhandlerError(err error) *asset.ErrorCode {
	switch err {
	case usecase.ErrorMongoTechnicalFail:
		return &asset.E().ErrorSpiderDB
	default:
		return &asset.E().GeneralSystemError
	}
}
