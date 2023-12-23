package handler

import (
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

type GetGeographinesHandler struct {
	thaiGeographiesUsecase domain.ThaiGeographiesUsecase
	log                    *logger.Logger
}

func NewGetGeographinesHandler(thaiGeographiesUsecase domain.ThaiGeographiesUsecase) *GetGeographinesHandler {
	return &GetGeographinesHandler{
		thaiGeographiesUsecase: thaiGeographiesUsecase,
		log:                    logger.L().Named("GetGeographinesHandler"),
	}
}

// =========================================================
// get all province
// =========================================================
func (h *GetGeographinesHandler) GetProvinceHandler(ctx *gin.Context) {
	log := h.log.WithContext(ctx)

	var req api_model.GetProvinceRequester
	var resp api_model.GetProvinceResponser

	if err := ctx.ShouldBind(&req); err != nil {
		log.Errorf("usecase request failed: %+v", err)
		assetError := asset.E().ErrorSpiderDB
		resp.Header.ErrorCode = assetError.ErrorCode
		resp.Header.Message = assetError.ErrorMessageEN
		ctx.JSON(assetError.StatusCode, resp)
		return
	}

	log.Infof("[GetProvinceHandler] get province handler start with req: %+v", req)

	if err := validator.Struct(req); err != nil {
		log.Errorf("validate request data fail, error: %+v", err)
		resp.Header.ErrorCode = asset.E().RequestDataFail.ErrorCode
		resp.Header.Message = asset.E().RequestDataFail.ErrorMessageEN
		ctx.JSON(asset.E().RequestDataFail.StatusCode, resp)
		return
	}

	provinceList, err := h.thaiGeographiesUsecase.GetAllProvince(ctx)
	if err != nil {
		assetErr := h.mapGetProvinceHandlerError(err)
		resp.Header.ErrorCode = assetErr.ErrorCode
		resp.Header.Message = assetErr.ErrorMessageEN
		ctx.JSON(assetErr.StatusCode, resp)
		return
	}

	resp.Data.Province = h.mapProvinceToResponseDataFormat(provinceList)

	resp.Header.ErrorCode = SUCCESS_CODE
	resp.Header.Message = SUCCESS_MESSAGE

	ctx.JSON(http.StatusOK, resp)
}

func (h *GetGeographinesHandler) mapGetProvinceHandlerError(err error) *asset.ErrorCode {
	switch err {
	case usecase.ErrorMongoTechnicalFail:
		return &asset.E().ErrorSpiderDB
	case usecase.ErrorThaiGeographiesUsecaseZeroProvince:
		return &asset.E().GeographiesNotFound
	default:
		return &asset.E().GeneralSystemError
	}
}

func (h *GetGeographinesHandler) mapProvinceToResponseDataFormat(provinceList []model.Province) []api_model.Province {

	var ProvinceListResp []api_model.Province

	for index, province := range provinceList {
		provinceResp := api_model.Province{
			Number: int16(index) + 1,
			NameTH: province.NameTH,
			NameEN: province.NameEN,
		}

		ProvinceListResp = append(ProvinceListResp, provinceResp)
	}

	return ProvinceListResp
}

// =========================================================
// get district by province name
// =========================================================
func (h *GetGeographinesHandler) GetDistrictByProvinceNameHandler(ctx *gin.Context) {
	log := h.log.WithContext(ctx)

	var req api_model.GetDistrictRequester
	var resp api_model.GetDistrictResponser

	if err := ctx.ShouldBind(&req); err != nil {
		log.Errorf("usecase request failed: %+v", err)
		assetError := asset.E().ErrorSpiderDB
		resp.Header.ErrorCode = assetError.ErrorCode
		resp.Header.Message = assetError.ErrorMessageEN
		ctx.JSON(assetError.StatusCode, resp)
		return
	}

	log.Infof("[GetDistrictByProvinceNameHandler] get province handler start with req: %+v", req)

	if err := validator.Struct(req); err != nil {
		log.Errorf("validate request data fail, error: %+v", err)
		resp.Header.ErrorCode = asset.E().RequestDataFail.ErrorCode
		resp.Header.Message = asset.E().RequestDataFail.ErrorMessageEN
		ctx.JSON(asset.E().RequestDataFail.StatusCode, resp)
		return
	}

	districtList, err := h.thaiGeographiesUsecase.GetDistictWithProvinceNameEN(ctx, req.Data.ProvinceNameEN)
	if err != nil {
		assetErr := h.mapGetDistrictByProvinceNameHandlerError(err)
		resp.Header.ErrorCode = assetErr.ErrorCode
		resp.Header.Message = assetErr.ErrorMessageEN
		ctx.JSON(assetErr.StatusCode, resp)
		return
	}

	resp.Data.District = h.mapDistrictToResponseDataFormat(districtList)

	resp.Header.ErrorCode = SUCCESS_CODE
	resp.Header.Message = SUCCESS_MESSAGE

	ctx.JSON(http.StatusOK, resp)
}

func (h *GetGeographinesHandler) mapGetDistrictByProvinceNameHandlerError(err error) *asset.ErrorCode {
	switch err {
	case usecase.ErrorMongoTechnicalFail:
		return &asset.E().ErrorSpiderDB
	case usecase.ErrorThaiGeographiesUsecaseDistrictNotFound:
		return &asset.E().GeographiesNotFound
	default:
		return &asset.E().GeneralSystemError
	}
}

func (h *GetGeographinesHandler) mapDistrictToResponseDataFormat(districtList []model.District) []api_model.District {

	var districtListResp []api_model.District

	for index, district := range districtList {
		districtResp := api_model.District{
			Number: int16(index) + 1,
			NameTH: district.NameTH,
			NameEN: district.NameEN,
		}

		districtListResp = append(districtListResp, districtResp)
	}

	return districtListResp
}
