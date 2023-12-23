package handler

import (
	"fmt"
	"net/http"
	api_model "spider-go/api/model"
	"spider-go/asset"
	"spider-go/domain"
	"spider-go/logger"
	"spider-go/model"
	"spider-go/usecase"
	"spider-go/utils/validator"

	"github.com/gin-gonic/gin"
)

type SpiderInfoHandler struct {
	spiderInfoUsecase      domain.SpiderInfoUsecase
	thaiGeographiesUsecase domain.ThaiGeographiesUsecase
	log                    *logger.Logger
}

func NewSpiderInfoHandler(spiderInfoUsecase domain.SpiderInfoUsecase, thaiGeographiesUsecase domain.ThaiGeographiesUsecase) *SpiderInfoHandler {
	return &SpiderInfoHandler{
		spiderInfoUsecase:      spiderInfoUsecase,
		thaiGeographiesUsecase: thaiGeographiesUsecase,
		log:                    logger.L().Named("SpiderInfoHandler"),
	}
}

// =========================================================
// get one spider info
// =========================================================
func (h *SpiderInfoHandler) GetOneSpiderInfoHandler(ctx *gin.Context) {
	log := h.log.WithContext(ctx)

	var req api_model.GetOneSpiderInfoRequester
	var resp api_model.GetOneSpiderInfoResponsor

	// read info
	if err := ctx.ShouldBind(&req); err != nil {
		log.Errorf("should bind request failed: %+v", err)
		resp.Header.ErrorCode = asset.E().GeneralSystemError.ErrorCode
		resp.Header.Message = asset.E().GeneralSystemError.ErrorMessageEN
		ctx.AbortWithStatusJSON(asset.E().GeneralSystemError.StatusCode, resp)
		return
	}

	log.Infof("get one spider info handler start with req: %v", req)

	if err := validator.Struct(req); err != nil {
		log.Errorf("validate request data fail, error: %+v", err)
		resp.Header.ErrorCode = asset.E().RequestDataFail.ErrorCode
		resp.Header.Message = asset.E().RequestDataFail.ErrorMessageEN
		ctx.JSON(asset.E().RequestDataFail.StatusCode, resp)
		return
	}

	if err := validator.Struct(req); err != nil {
		log.Errorf("validate request data fail, error: %+v", err)
		resp.Header.ErrorCode = asset.E().RequestDataFail.ErrorCode
		resp.Header.Message = asset.E().RequestDataFail.ErrorMessageEN
		ctx.JSON(asset.E().RequestDataFail.StatusCode, resp)
		return
	}

	var spiderInfo *model.SpiderInfo
	var err error
	err = nil
	// if req.Header.Username == "" || req.Header.UserKey == "" {
	// 	spiderInfo, err = h.spiderInfoUsecase.GetSpiderInfoUsecase(ctx, req.Data.SpiderUUID, "", "", false)
	// } else {
	// 	spiderInfo, err = h.spiderInfoUsecase.GetSpiderInfoUsecase(ctx, req.Data.SpiderUUID, req.Header.Username, req.Header.UserKey, true)
	// }

	spiderInfo, err = h.spiderInfoUsecase.GetSpiderInfoUsecase(ctx, req.Data.SpiderUUID, req.Header.Username)

	if err != nil {
		log.Errorf("[GetOneSpiderInfoHandler] get spider info usecase failed, error: %v", err)
		assetErr := h.mapGetOneSpiderInfoHandler(err)
		resp.Header.ErrorCode = assetErr.ErrorCode
		resp.Header.Message = assetErr.ErrorMessageEN
		ctx.JSON(assetErr.StatusCode, resp)
		return
	}

	resp.Header.ErrorCode = SUCCESS_CODE
	resp.Header.Message = ""
	resp.Data = *mapSpiderInfoModel(spiderInfo)
	ctx.JSON(http.StatusOK, resp)

}

func (h *SpiderInfoHandler) mapGetOneSpiderInfoHandler(err error) *asset.ErrorCode {
	switch err {
	case usecase.ErrorRedisConnection:
		return &asset.E().ErrorTempDB
	case usecase.ErrorValidateUserFaild:
		return &asset.E().UserNotLogin
	case usecase.ErrorSpiderInfoUsecaseSpiderNotFound:
		return &asset.E().SpiderNotFound
	case usecase.ErrorMongoConnection, usecase.ErrorMongoTechnicalFail:
		return &asset.E().ErrorSpiderDB
	case usecase.ErrorSpiderInfoUsecaseAccountInsufficientPermissions:
		return &asset.E().InsufficientUserRights
	default:
		return &asset.E().GeneralSystemError
	}
}

// *************************************************

// =========================================================
// get spider image
// =========================================================
func (h *SpiderInfoHandler) GetSpiderImagesHandler(ctx *gin.Context) {
	log := h.log.WithContext(ctx)

	var req api_model.GetSpiderImageRequester
	var resp api_model.GetSpiderImageResponsor

	if err := ctx.ShouldBind(&req); err != nil {
		log.Errorf("should bind request failed: %+v", err)
		resp.Header.ErrorCode = asset.E().GeneralSystemError.ErrorCode
		resp.Header.Message = asset.E().GeneralSystemError.ErrorMessageEN
		ctx.AbortWithStatusJSON(asset.E().GeneralSystemError.StatusCode, resp)
		return
	}

	log.Infof("get spider image handler start with req: %v", req)

	if err := validator.Struct(req); err != nil {
		log.Errorf("validate request data fail, error: %+v", err)
		resp.Header.ErrorCode = asset.E().RequestDataFail.ErrorCode
		resp.Header.Message = asset.E().RequestDataFail.ErrorMessageEN
		ctx.JSON(asset.E().RequestDataFail.StatusCode, resp)
		return
	}

	spiderImageEndcode, err := h.spiderInfoUsecase.GetSpiderImagesUsecase(ctx, req.Data.SpiderImageList)
	if err != nil {
		// it's only one error, no need to map
		resp.Header.ErrorCode = asset.E().GeneralSystemError.ErrorCode
		resp.Header.Message = asset.E().GeneralSystemError.ErrorMessageEN
		ctx.AbortWithStatusJSON(asset.E().GeneralSystemError.StatusCode, resp)
		return
	}

	resp.Header.ErrorCode = SUCCESS_CODE
	resp.Header.Message = ""
	resp.Data.SpiderImageList = h.mapGetSpiderImagesHandler(spiderImageEndcode)
	ctx.JSON(http.StatusOK, resp)

}

func (h *SpiderInfoHandler) mapGetSpiderImagesHandler(SpiderImageEndcodeList []model.SpiderImageList) []api_model.SpiderImageList {
	var responseData []api_model.SpiderImageList
	for _, data := range SpiderImageEndcodeList {
		responseData = append(responseData, api_model.SpiderImageList{
			Title: data.Name,
			Src:   data.ImageBase64,
		})
	}

	return responseData
}

// *************************************************

// =========================================================
// get spider info list manager
// =========================================================
func (h *SpiderInfoHandler) GetSpiderInfoListManagerHandler(ctx *gin.Context) {
	log := h.log.WithContext(ctx)

	var req api_model.GetSpiderInfoListManagerRequester
	var resp api_model.GetSpiderInfoListManagerResponsor

	if err := ctx.ShouldBind(&req); err != nil {
		log.Errorf("should bind request failed: %+v", err)
		resp.Header.ErrorCode = asset.E().GeneralSystemError.ErrorCode
		resp.Header.Message = asset.E().GeneralSystemError.ErrorMessageEN
		ctx.AbortWithStatusJSON(asset.E().GeneralSystemError.StatusCode, resp)
		return
	}

	log.Infof("get spider info list manager handler start with req: %v", req)

	if err := validator.Struct(req); err != nil {
		log.Errorf("validate request data fail, error: %+v", err)
		resp.Header.ErrorCode = asset.E().RequestDataFail.ErrorCode
		resp.Header.Message = asset.E().RequestDataFail.ErrorMessageEN
		ctx.JSON(asset.E().RequestDataFail.StatusCode, resp)
		return
	}

	// username := req.Header.Username
	// key := req.Header.UserKey
	// page := req.Data.Page
	// size := req.Data.Size

	err := fmt.Errorf("test: JWT")
	// spiderInfoListManager, err := h.spiderInfoUsecase.GetSpiderInfoListManager(ctx, username, key, page, size)
	if err != nil {
		log.Errorf("GetSpiderInfoListManager return error: %+v", err)
		resp.Header.ErrorCode = asset.E().GeneralSystemError.ErrorCode
		resp.Header.Message = asset.E().GeneralSystemError.ErrorMessageEN
		ctx.AbortWithStatusJSON(asset.E().GeneralSystemError.StatusCode, resp)
		return
	}

	var spiderInfoList []api_model.SpiderInfo

	// for _, spiderList := range spiderInfoListManager {
	// 	spiderInfo := mapSpiderInfoModel(&spiderList)

	// 	spiderInfoList = append(spiderInfoList, *spiderInfo)
	// }

	resp.Header.ErrorCode = SUCCESS_CODE
	resp.Header.Message = ""
	resp.Data.SpiderInfoList = spiderInfoList
	ctx.JSON(http.StatusOK, resp)
}

// *************************************************

// =========================================================
// get spider info filter by geographies
// =========================================================
func (h *SpiderInfoHandler) GetSpiderInfoByGeographiesHandler(ctx *gin.Context) {
	log := h.log.WithContext(ctx)

	var req api_model.GetSpiderInfoByGeographiesRequester
	var resp api_model.GetSpiderInfoByGeographieResponsor

	if err := ctx.ShouldBind(&req); err != nil {
		log.Errorf("should bind request failed: %+v", err)
		resp.Header.ErrorCode = asset.E().GeneralSystemError.ErrorCode
		resp.Header.Message = asset.E().GeneralSystemError.ErrorMessageEN
		ctx.AbortWithStatusJSON(asset.E().GeneralSystemError.StatusCode, resp)
		return
	}

	log.Infof("[GetSpiderInfoByGeographiesHandler] get spider info by geographies handler start with req: %v", req)

	if err := validator.Struct(req); err != nil {
		log.Errorf("validate request data fail, error: %+v", err)
		resp.Header.ErrorCode = asset.E().RequestDataFail.ErrorCode
		resp.Header.Message = asset.E().RequestDataFail.ErrorMessageEN
		ctx.JSON(asset.E().RequestDataFail.StatusCode, resp)
		return
	}

	spiderInfoList, err := h.spiderInfoUsecase.GetSpiderInfoListByGeographies(ctx, req.Data.Province, req.Data.District, req.Data.Position)
	if err != nil {
		log.Errorf("[GetSpiderInfoByGeographiesHandler] usecase failed, error: %+v", err)
		assetErr := h.mapGetSpiderInfoByGeographiesHandlerError(err)
		resp.Header.ErrorCode = assetErr.ErrorCode
		resp.Header.Message = assetErr.ErrorMessageEN
		ctx.JSON(assetErr.StatusCode, resp)
		return
	}

	var spiderInfoListResp []api_model.SpiderInfo

	for _, spiderInfo := range spiderInfoList {
		spiderInfoListResp = append(spiderInfoListResp, *mapSpiderInfoModel(&spiderInfo))
	}

	resp.Header.ErrorCode = SUCCESS_CODE
	resp.Header.Message = ""
	resp.Data.SpiderInfoList = spiderInfoListResp
	ctx.JSON(http.StatusOK, resp)
}

func (h *SpiderInfoHandler) mapGetSpiderInfoByGeographiesHandlerError(err error) *asset.ErrorCode {
	switch err {
	case usecase.ErrorSpiderInfoUsecaseValidateDataFail:
		return &asset.E().RequestDataFail
	case usecase.ErrorMongoTechnicalFail:
		return &asset.E().ErrorSpiderDB
	case usecase.ErrorSpiderInfoUsecaseSpiderNotFound:
		return &asset.E().SpiderNotFound
	default:
		return &asset.E().GeneralSystemError
	}
}

// *************************************************

// =========================================================
// get geographies by spider type
// =========================================================
func (h *SpiderInfoHandler) GetGeographiesBySpiderTypeHandler(ctx *gin.Context) {
	log := h.log.WithContext(ctx)

	var req api_model.GetGeoGraphiesBySpiderTypeRequester
	var resp api_model.GetGeoGraphiesBySpiderTypeResponser

	if err := ctx.ShouldBind(&req); err != nil {
		log.Errorf("should bind request failed: %+v", err)
		resp.Header.ErrorCode = asset.E().GeneralSystemError.ErrorCode
		resp.Header.Message = asset.E().GeneralSystemError.ErrorMessageEN
		ctx.AbortWithStatusJSON(asset.E().GeneralSystemError.StatusCode, resp)
		return
	}

	log.Infof("[GetGeographiesBySpiderTypeHandler] get geographies by spider type handler start with req: %v", req)

	if err := validator.Struct(req); err != nil {
		log.Errorf("validate request data fail, error: %+v", err)
		resp.Header.ErrorCode = asset.E().RequestDataFail.ErrorCode
		resp.Header.Message = asset.E().RequestDataFail.ErrorMessageEN
		ctx.JSON(asset.E().RequestDataFail.StatusCode, resp)
		return
	}

	// NOTE: call use case
	locationResult, err := h.thaiGeographiesUsecase.GetGeographiesBySpiderType(ctx, req.Data.Family, req.Data.Genus, req.Data.Species)
	if err != nil {
		log.Errorf("[GetSpiderInfoByGeographiesHandler] usecase failed, error: %+v", err)
		assetErr := h.mapGetGeographiesBySpiderTypeHandlerError(err)
		resp.Header.ErrorCode = assetErr.ErrorCode
		resp.Header.Message = assetErr.ErrorMessageEN
		ctx.JSON(assetErr.StatusCode, resp)
		return
	}

	resp.Header.ErrorCode = SUCCESS_CODE
	resp.Header.Message = ""
	resp.Data.LocationResult = locationResult
	ctx.JSON(http.StatusOK, resp)

}

func (h *SpiderInfoHandler) mapGetGeographiesBySpiderTypeHandlerError(err error) *asset.ErrorCode {
	switch err {
	case usecase.ErrorThaiGeographiesUsecaseValidateDataRequestFailed:
		return &asset.E().RequestDataFail
	case usecase.ErrorThaiGeographiesUsecaseSpiderNotFound:
		return &asset.E().RequestDataNotFound
	case usecase.ErrorMongoTechnicalFail:
		return &asset.E().ErrorSpiderDB
	default:
		return &asset.E().GeneralSystemError
	}
}

// *************************************************

// =========================================================
// get spider info by locality
// =========================================================
func (h *SpiderInfoHandler) GetSpiderInfoByLocalityHandler(ctx *gin.Context) {
	log := h.log.WithContext(ctx)

	var req api_model.GetSpiderInfoByLocalityRequester
	var resp api_model.GetSpiderInfoByLocalityResponser

	if err := ctx.ShouldBind(&req); err != nil {
		log.Errorf("should bind request failed: %+v", err)
		resp.Header.ErrorCode = asset.E().GeneralSystemError.ErrorCode
		resp.Header.Message = asset.E().GeneralSystemError.ErrorMessageEN
		ctx.AbortWithStatusJSON(asset.E().GeneralSystemError.StatusCode, resp)
		return
	}

	log.Infof("[GetSpiderInfoByLocalityHandler] get spider info by locality handler start with req: %v", req)

	if err := validator.Struct(req); err != nil {
		log.Errorf("validate request data fail, error: %+v", err)
		resp.Header.ErrorCode = asset.E().RequestDataFail.ErrorCode
		resp.Header.Message = asset.E().RequestDataFail.ErrorMessageEN
		ctx.JSON(asset.E().RequestDataFail.StatusCode, resp)
		return
	}

	if err := validator.Struct(req); err != nil {
		log.Errorf("validate request data fail, error: %+v", err)
		resp.Header.ErrorCode = asset.E().RequestDataFail.ErrorCode
		resp.Header.Message = asset.E().RequestDataFail.ErrorMessageEN
		ctx.JSON(asset.E().RequestDataFail.StatusCode, resp)
		return
	}

	spiderInfoList, err := h.spiderInfoUsecase.GetSpiderInfoListByLocality(ctx, req.Data.LocalityName, req.Data.Page, req.Data.Size)
	if err != nil {
		assetErr := h.mapGetSpiderInfoByLocalityHandlerError(err)
		resp.Header.ErrorCode = assetErr.ErrorCode
		resp.Header.Message = assetErr.ErrorMessageEN
		ctx.JSON(assetErr.StatusCode, resp)
		return
	}

	var spiderInfoListResp []api_model.SpiderInfo

	for _, spiderInfo := range spiderInfoList {
		spiderInfoListResp = append(spiderInfoListResp, *mapSpiderInfoModel(&spiderInfo))
	}

	resp.Data.SpiderInfoList = spiderInfoListResp
	resp.Header.ErrorCode = SUCCESS_CODE
	resp.Header.Message = ""
	ctx.JSON(http.StatusOK, resp)
}

func (h *SpiderInfoHandler) mapGetSpiderInfoByLocalityHandlerError(err error) *asset.ErrorCode {
	switch err {
	case usecase.ErrorSpiderInfoUsecaseSpiderNotFound:
		return &asset.E().SpiderNotFound
	case usecase.ErrorMongoTechnicalFail:
		return &asset.E().ErrorSpiderDB
	case usecase.ErrorSpiderInfoUsecaseValidateDataFail:
		return &asset.E().RequestDataFail
	default:
		return &asset.E().GeneralSystemError
	}
}

// *************************************************

// =========================================================
// get spider info by locality
// =========================================================
func (h *SpiderInfoHandler) GetSpiderListBySpiderTypeHandler(ctx *gin.Context) {
	log := h.log.WithContext(ctx)

	var req api_model.GetSpiderListBySpiderTypeRequester
	var resp api_model.GetSpiderListBySpiderTypeResponser

	if err := ctx.ShouldBind(&req); err != nil {
		log.Errorf("[GetSpiderListBySpiderTypeHandler] should bind error: %+v", err)
		resp.Header.ErrorCode = asset.E().GeneralSystemError.ErrorCode
		resp.Header.Message = asset.E().GeneralSystemError.ErrorMessageEN
		ctx.AbortWithStatusJSON(asset.E().GeneralSystemError.StatusCode, resp)
		return
	}

	log.Infof("[GetSpiderListBySpiderTypeHandler] get spider list by spider type handler start with req: %v", req)

	if err := validator.Struct(req); err != nil {
		log.Errorf("validate request data fail, error: %+v", err)
		resp.Header.ErrorCode = asset.E().RequestDataFail.ErrorCode
		resp.Header.Message = asset.E().RequestDataFail.ErrorMessageEN
		ctx.JSON(asset.E().RequestDataFail.StatusCode, resp)
		return
	}

	param := model.GetSpiderListBySpiderTypeParam{
		Family:  req.Data.Family,
		Genus:   req.Data.Genus,
		Species: req.Data.Species,
		Page:    req.Data.Page,
		Size:    req.Data.Size,
	}

	spiderInfoList, err := h.spiderInfoUsecase.GetSpiderListBySpiderTypeUsecase(ctx, param)
	if err != nil {
		assetErr := h.mapGetSpiderInfoByLocalityHandlerError(err)
		resp.Header.ErrorCode = assetErr.ErrorCode
		resp.Header.Message = assetErr.ErrorMessageEN

		log.Errorf("[GetSpiderListBySpiderTypeHandler] response: %+v", resp)
		ctx.JSON(assetErr.StatusCode, resp)
		return
	}

	var spiderInfoListResp []api_model.SpiderInfo

	for _, spiderInfo := range spiderInfoList {
		spiderInfoListResp = append(spiderInfoListResp, *mapSpiderInfoModel(&spiderInfo))
	}

	resp.Data.SpiderInfoList = spiderInfoListResp
	resp.Header.ErrorCode = SUCCESS_CODE
	resp.Header.Message = ""
	log.Infof("[GetSpiderListBySpiderTypeHandler] response: %+v", resp)
	ctx.JSON(http.StatusOK, resp)
}
